package dom

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
	"sync"
)

var (
	eventDocumentUpdated   = events.New("doc_updated")
	eventChildNodeInserted = events.New("child_inserted")
	eventChildNodeRemoved  = events.New("child_removed")
)

type (
	DocumentUpdatedListener func(ctx context.Context)

	AttrModifiedListener func(ctx context.Context, nodeID dom.NodeID, name, value string)

	AttrRemovedListener func(ctx context.Context, nodeID dom.NodeID, name string)

	ChildNodeCountUpdatedListener func(ctx context.Context, nodeID dom.NodeID, count int)

	ChildNodeInsertedListener func(ctx context.Context, nodeID, previousNodeID dom.NodeID, node dom.Node)

	ChildNodeRemovedListener func(ctx context.Context, nodeID, previousNodeID dom.NodeID)

	Manager struct {
		mu        sync.RWMutex
		logger    *zerolog.Logger
		client    *cdp.Client
		events    *events.Loop
		mouse     *input.Mouse
		keyboard  *input.Keyboard
		mainFrame *AtomicFrameID
		frames    *AtomicFrameCollection
		cancel    context.CancelFunc
	}
)

// a dirty workaround to let pass the vet test
func createContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func New(
	logger *zerolog.Logger,
	client *cdp.Client,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
) (manager *Manager, err error) {
	ctx, cancel := createContext()

	closers := make([]io.Closer, 0, 10)

	defer func() {
		if err != nil {
			common.CloseAll(logger, closers, "failed to close a DOM event stream")
		}
	}()

	onContentReady, err := client.Page.DOMContentEventFired(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onContentReady)

	onDocUpdated, err := client.DOM.DocumentUpdated(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onDocUpdated)

	onAttrModified, err := client.DOM.AttributeModified(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onAttrModified)

	onAttrRemoved, err := client.DOM.AttributeRemoved(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onAttrRemoved)

	onChildCountUpdated, err := client.DOM.ChildNodeCountUpdated(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onChildCountUpdated)

	onChildNodeInserted, err := client.DOM.ChildNodeInserted(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onChildNodeInserted)

	onChildNodeRemoved, err := client.DOM.ChildNodeRemoved(ctx)

	if err != nil {
		return nil, err
	}

	closers = append(closers, onChildNodeRemoved)

	eventLoop := events.NewLoop()

	eventLoop.AddSource(events.NewSource(eventDocumentUpdated, onDocUpdated, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.DocumentUpdatedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventChildNodeInserted, onChildNodeInserted, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.ChildNodeInsertedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventChildNodeRemoved, onChildNodeRemoved, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.ChildNodeRemovedClient).Recv()
	}))

	manager = new(Manager)
	manager.logger = logger
	manager.client = client
	manager.events = eventLoop
	manager.mouse = mouse
	manager.keyboard = keyboard
	manager.mainFrame = NewAtomicFrameID()
	manager.frames = NewAtomicFrameCollection()
	manager.cancel = cancel

	eventLoop.Run(ctx)

	return manager, nil
}

func (m *Manager) Close() error {
	errs := make([]error, 0, m.frames.Length()+1)

	if m.cancel != nil {
		m.cancel()
		m.cancel = nil
	}

	m.frames.ForEach(func(f Frame, key page.FrameID) bool {
		// if initialized
		if f.node != nil {
			if err := f.node.Close(); err != nil {
				errs = append(errs, err)
			}
		}

		return true
	})

	if len(errs) > 0 {
		return core.Errors(errs...)
	}

	return nil
}

func (m *Manager) GetMainFrame() *HTMLDocument {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mainFrameID := m.mainFrame.Get()

	if mainFrameID == "" {
		return nil
	}

	mainFrame, exists := m.frames.Get(mainFrameID)

	if exists {
		return mainFrame.node
	}

	return nil
}

func (m *Manager) SetMainFrame(doc *HTMLDocument) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mainFrameID := m.mainFrame.Get()

	if mainFrameID != "" {
		if err := m.removeFrameRecursivelyInternal(mainFrameID); err != nil {
			m.logger.Error().Err(err).Msg("failed to close previous main frame")
		}
	}

	m.mainFrame.Set(doc.frameTree.Frame.ID)

	m.addPreloadedFrame(doc)
}

func (m *Manager) AddFrame(frame page.FrameTree) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.addFrameInternal(frame)
}

func (m *Manager) RemoveFrame(frameID page.FrameID) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.removeFrameInternal(frameID)
}

func (m *Manager) RemoveFrameRecursively(frameID page.FrameID) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.removeFrameRecursivelyInternal(frameID)
}

func (m *Manager) RemoveFramesByParentID(parentFrameID page.FrameID) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	frame, found := m.frames.Get(parentFrameID)

	if !found {
		return errors.New("frame not found")
	}

	for _, child := range frame.tree.ChildFrames {
		if err := m.removeFrameRecursivelyInternal(child.Frame.ID); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) GetFrameNode(ctx context.Context, frameID page.FrameID) (*HTMLDocument, error) {
	return m.getFrameInternal(ctx, frameID)
}

func (m *Manager) GetFrameTree(_ context.Context, frameID page.FrameID) (page.FrameTree, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	frame, found := m.frames.Get(frameID)

	if !found {
		return page.FrameTree{}, core.ErrNotFound
	}

	return frame.tree, nil
}

func (m *Manager) GetFrameNodes(ctx context.Context) (*values.Array, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	arr := values.NewArray(m.frames.Length())

	for _, f := range m.frames.ToSlice() {
		doc, err := m.getFrameInternal(ctx, f.tree.Frame.ID)

		if err != nil {
			return nil, err
		}

		arr.Push(doc)
	}

	return arr, nil
}

func (m *Manager) AddDocumentUpdatedListener(listener DocumentUpdatedListener) events.ListenerID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.events.AddListener(eventDocumentUpdated, func(ctx context.Context, _ interface{}) bool {
		listener(ctx)

		return true
	})
}

func (m *Manager) RemoveReloadListener(listenerID events.ListenerID) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.events.RemoveListener(eventDocumentUpdated, listenerID)
}

func (m *Manager) AddChildNodeInsertedListener(listener ChildNodeInsertedListener) events.ListenerID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.events.AddListener(eventChildNodeInserted, func(ctx context.Context, message interface{}) bool {
		reply := message.(*dom.ChildNodeInsertedReply)

		listener(ctx, reply.ParentNodeID, reply.PreviousNodeID, reply.Node)

		return true
	})
}

func (m *Manager) RemoveChildNodeInsertedListener(listenerID events.ListenerID) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.events.RemoveListener(eventChildNodeInserted, listenerID)
}

func (m *Manager) AddChildNodeRemovedListener(listener ChildNodeRemovedListener) events.ListenerID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.events.AddListener(eventChildNodeRemoved, func(ctx context.Context, message interface{}) bool {
		reply := message.(*dom.ChildNodeRemovedReply)

		listener(ctx, reply.ParentNodeID, reply.NodeID)

		return true
	})
}

func (m *Manager) RemoveChildNodeRemovedListener(listenerID events.ListenerID) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.events.RemoveListener(eventChildNodeRemoved, listenerID)
}

func (m *Manager) addFrameInternal(frame page.FrameTree) {
	m.frames.Set(frame.Frame.ID, Frame{
		tree: frame,
		node: nil,
	})

	for _, child := range frame.ChildFrames {
		m.addFrameInternal(child)
	}
}

func (m *Manager) addPreloadedFrame(doc *HTMLDocument) {
	m.frames.Set(doc.frameTree.Frame.ID, Frame{
		tree: doc.frameTree,
		node: doc,
	})

	for _, child := range doc.frameTree.ChildFrames {
		m.addFrameInternal(child)
	}
}

func (m *Manager) getFrameInternal(ctx context.Context, frameID page.FrameID) (*HTMLDocument, error) {
	frame, found := m.frames.Get(frameID)

	if !found {
		return nil, core.ErrNotFound
	}

	// frame is initialized
	if frame.node != nil {
		return frame.node, nil
	}

	// the frames is not loaded yet
	node, execID, err := resolveFrame(ctx, m.client, frameID)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve frame node: %s", frameID)
	}

	doc, err := LoadHTMLDocument(
		ctx,
		m.logger,
		m.client,
		m,
		m.mouse,
		m.keyboard,
		node,
		frame.tree,
		execID,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load frame document")
	}

	frame.node = doc

	return doc, nil
}

func (m *Manager) removeFrameInternal(frameID page.FrameID) error {
	current, exists := m.frames.Get(frameID)

	if !exists {
		return core.Error(core.ErrNotFound, "frame")
	}

	m.frames.Remove(frameID)

	mainFrameID := m.mainFrame.Get()

	if frameID == mainFrameID {
		m.mainFrame.Reset()
	}

	if current.node == nil {
		return nil
	}

	return current.node.Close()
}

func (m *Manager) removeFrameRecursivelyInternal(frameID page.FrameID) error {
	parent, exists := m.frames.Get(frameID)

	if !exists {
		return core.Error(core.ErrNotFound, "frame")
	}

	for _, child := range parent.tree.ChildFrames {
		if err := m.removeFrameRecursivelyInternal(child.Frame.ID); err != nil {
			return err
		}
	}

	return m.removeFrameInternal(frameID)
}
