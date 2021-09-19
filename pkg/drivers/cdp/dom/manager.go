package dom

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type Manager struct {
	mu        sync.RWMutex
	logger    zerolog.Logger
	client    *cdp.Client
	mouse     *input.Mouse
	keyboard  *input.Keyboard
	mainFrame *AtomicFrameID
	frames    *AtomicFrameCollection
}

func New(
	logger zerolog.Logger,
	client *cdp.Client,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
) (manager *Manager, err error) {

	manager = new(Manager)
	manager.logger = logging.WithName(logger.With(), "dom_manager").Logger()
	manager.client = client
	manager.mouse = mouse
	manager.keyboard = keyboard
	manager.mainFrame = NewAtomicFrameID()
	manager.frames = NewAtomicFrameCollection()

	return manager, nil
}

func (m *Manager) Close() error {
	errs := make([]error, 0, m.frames.Length()+1)

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

func (m *Manager) LoadRootDocument(ctx context.Context) (*HTMLDocument, error) {
	ftRepl, err := m.client.Page.GetFrameTree(ctx)

	if err != nil {
		return nil, err
	}

	return m.LoadDocument(ctx, ftRepl.FrameTree)
}

func (m *Manager) LoadDocument(ctx context.Context, frame page.FrameTree) (*HTMLDocument, error) {
	exec, err := eval.Create(ctx, m.logger, m.client, frame.Frame.ID)

	if err != nil {
		return nil, err
	}

	inputs := input.New(m.logger, m.client, exec, m.keyboard, m.mouse)

	ref, err := exec.EvalRef(ctx, templates.GetDocument())

	if err != nil {
		return nil, errors.Wrap(err, "failed to load root element")
	}

	exec.SetLoader(NewNodeLoader(m))

	rootElement := NewHTMLElement(
		m.logger,
		m.client,
		m,
		inputs,
		exec,
		*ref.ObjectID,
	)

	return NewHTMLDocument(
		m.logger,
		m.client,
		m,
		inputs,
		exec,
		rootElement,
		frame,
	), nil
}

func (m *Manager) ResolveElement(ctx context.Context, frameID page.FrameID, id runtime.RemoteObjectID) (*HTMLElement, error) {
	doc, err := m.GetFrameNode(ctx, frameID)

	if err != nil {
		return nil, err
	}

	return NewHTMLElement(
		m.logger,
		m.client,
		m,
		doc.input,
		doc.eval,
		id,
	), nil
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

	// the frame is not loaded yet
	doc, err := m.LoadDocument(ctx, frame.tree)

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
