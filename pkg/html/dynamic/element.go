package dynamic

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MontFerret/ferret/pkg/html/common"
	"github.com/MontFerret/ferret/pkg/html/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/html/dynamic/events"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/rs/zerolog"
)

const DefaultTimeout = time.Second * 30

var emptyBackendID = dom.BackendNodeID(0)
var emptyObjectID = ""

type (
	HTMLElementIdentity struct {
		nodeID    dom.NodeID
		backendID dom.BackendNodeID
		objectID  runtime.RemoteObjectID
	}

	HTMLElement struct {
		sync.Mutex
		logger         *zerolog.Logger
		client         *cdp.Client
		events         *events.EventBroker
		connected      values.Boolean
		id             *HTMLElementIdentity
		nodeType       values.Int
		nodeName       values.String
		innerHTML      values.String
		innerText      *common.LazyValue
		value          core.Value
		rawAttrs       []string
		attributes     *common.LazyValue
		children       []*HTMLElementIdentity
		loadedChildren *common.LazyValue
	}
)

func LoadElement(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	nodeID dom.NodeID,
	backendID dom.BackendNodeID,
) (*HTMLElement, error) {
	if client == nil {
		return nil, core.Error(core.ErrMissedArgument, "client")
	}

	// getting a remote object that represents the current DOM Node
	var args *dom.ResolveNodeArgs

	if backendID > 0 {
		args = dom.NewResolveNodeArgs().SetBackendNodeID(backendID)
	} else {
		args = dom.NewResolveNodeArgs().SetNodeID(nodeID)
	}

	obj, err := client.DOM.ResolveNode(ctx, args)

	if err != nil {
		return nil, err
	}

	if obj.Object.ObjectID == nil {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("element %d", nodeID))
	}

	objectID := *obj.Object.ObjectID

	node, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetObjectID(objectID).
			SetDepth(1),
	)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(nodeID)))
	}

	id := new(HTMLElementIdentity)
	id.nodeID = nodeID
	id.objectID = objectID

	if backendID > 0 {
		id.backendID = backendID
	} else {
		id.backendID = node.Node.BackendNodeID
	}

	innerHTML, err := loadInnerHTML(ctx, client, id)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(nodeID)))
	}

	var val string

	if node.Node.Value != nil {
		val = *node.Node.Value
	}

	return NewHTMLElement(
		logger,
		client,
		broker,
		id,
		node.Node.NodeType,
		node.Node.NodeName,
		node.Node.Attributes,
		val,
		innerHTML,
		createChildrenArray(node.Node.Children),
	), nil
}

func NewHTMLElement(
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	id *HTMLElementIdentity,
	nodeType int,
	nodeName string,
	attributes []string,
	value string,
	innerHTML values.String,
	children []*HTMLElementIdentity,
) *HTMLElement {
	el := new(HTMLElement)
	el.logger = logger
	el.client = client
	el.events = broker
	el.connected = values.True
	el.id = id
	el.nodeType = values.NewInt(nodeType)
	el.nodeName = values.NewString(nodeName)
	el.innerHTML = innerHTML
	el.innerText = common.NewLazyValue(el.loadInnerText)
	el.rawAttrs = attributes
	el.attributes = common.NewLazyValue(el.loadAttrs)
	el.value = values.EmptyString
	el.loadedChildren = common.NewLazyValue(el.loadChildren)
	el.value = values.NewString(value)
	el.children = children

	broker.AddEventListener("reload", el.handlePageReload)
	broker.AddEventListener("attr:modified", el.handleAttrModified)
	broker.AddEventListener("attr:removed", el.handleAttrRemoved)
	broker.AddEventListener("children:count", el.handleChildrenCountChanged)
	broker.AddEventListener("children:inserted", el.handleChildInserted)
	broker.AddEventListener("children:deleted", el.handleChildDeleted)

	return el
}

func (el *HTMLElement) Close() error {
	el.Lock()
	defer el.Unlock()

	// already closed
	if !el.connected {
		return nil
	}

	el.connected = false
	el.events.RemoveEventListener("reload", el.handlePageReload)
	el.events.RemoveEventListener("attr:modified", el.handleAttrModified)
	el.events.RemoveEventListener("attr:removed", el.handleAttrRemoved)
	el.events.RemoveEventListener("children:count", el.handleChildrenCountChanged)
	el.events.RemoveEventListener("children:inserted", el.handleChildInserted)
	el.events.RemoveEventListener("children:deleted", el.handleChildDeleted)

	return nil
}

func (el *HTMLElement) Type() core.Type {
	return core.HTMLElementType
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	val, err := el.innerText.Read()

	if err != nil {
		return nil, err
	}

	return json.Marshal(val.String())
}

func (el *HTMLElement) String() string {
	return el.InnerHTML().String()
}

func (el *HTMLElement) Compare(other core.Value) int {
	switch other.Type() {
	case core.HTMLDocumentType:
		other := other.(*HTMLElement)

		id := int(el.id.backendID)
		otherID := int(other.id.backendID)

		if id == otherID {
			return 0
		}

		if id > otherID {
			return 1
		}

		return -1
	default:
		if other.Type() > core.HTMLElementType {
			return -1
		}

		return 1
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el
}

func (el *HTMLElement) Hash() uint64 {
	el.Lock()
	defer el.Unlock()

	h := fnv.New64a()

	h.Write([]byte(el.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(el.innerHTML))

	return h.Sum64()
}

func (el *HTMLElement) Value() core.Value {
	if !el.IsConnected() {
		return el.value
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	val, err := eval.Property(ctx, el.client, el.id.objectID, "value")

	if err != nil {
		el.logError(err).Msg("failed to get node value")

		return el.value
	}

	el.value = val

	return val
}

func (el *HTMLElement) Clone() core.Value {
	return values.None
}

func (el *HTMLElement) Length() values.Int {
	return values.NewInt(len(el.children))
}

func (el *HTMLElement) NodeType() values.Int {
	return el.nodeType
}

func (el *HTMLElement) NodeName() values.String {
	return el.nodeName
}

func (el *HTMLElement) GetAttributes() core.Value {
	val, err := el.attributes.Read()

	if err != nil {
		return values.None
	}

	// returning shallow copy
	return val.Clone()
}

func (el *HTMLElement) GetAttribute(name values.String) core.Value {
	attrs, err := el.attributes.Read()

	if err != nil {
		return values.None
	}

	val, found := attrs.(*values.Object).Get(name)

	if !found {
		return values.None
	}

	return val
}

func (el *HTMLElement) GetChildNodes() core.Value {
	val, err := el.loadedChildren.Read()

	if err != nil {
		return values.NewArray(0)
	}

	return val
}

func (el *HTMLElement) GetChildNode(idx values.Int) core.Value {
	val, err := el.loadedChildren.Read()

	if err != nil {
		return values.None
	}

	return val.(*values.Array).Get(idx)
}

func (el *HTMLElement) QuerySelector(selector values.String) core.Value {
	if !el.IsConnected() {
		return values.None
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(el.id.nodeID, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.None
	}

	res, err := LoadElement(ctx, el.logger, el.client, el.events, found.NodeID, emptyBackendID)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to load a child node by selector")

		return values.None
	}

	return res
}

func (el *HTMLElement) QuerySelectorAll(selector values.String) core.Value {
	if !el.IsConnected() {
		return values.NewArray(0)
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(el.id.nodeID, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.None
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		childEl, err := LoadElement(ctx, el.logger, el.client, el.events, id, emptyBackendID)

		if err != nil {
			el.logError(err).
				Str("selector", selector.String()).
				Msg("failed to load nodes by selector")

			// close elements that are already loaded, but won't be used because of the error
			if arr.Length() > 0 {
				arr.ForEach(func(e core.Value, _ int) bool {
					e.(*HTMLElement).Close()

					return true
				})
			}

			return values.None
		}

		arr.Push(childEl)
	}

	return arr
}

func (el *HTMLElement) WaitForClass(class values.String, timeout values.Int) error {
	task := events.NewWaitTask(
		func() (core.Value, error) {
			current := el.GetAttribute("class")

			if current.Type() != core.StringType {
				return values.None, nil
			}

			str := current.(values.String)
			classStr := string(class)
			classes := strings.Split(string(str), " ")

			for _, c := range classes {
				if c == classStr {
					return values.True, nil
				}
			}

			return values.None, nil
		},
		time.Millisecond*time.Duration(timeout),
		events.DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (el *HTMLElement) InnerText() values.String {
	val, err := el.innerText.Read()

	if err != nil {
		return values.EmptyString
	}

	return val.(values.String)
}

func (el *HTMLElement) InnerTextBySelector(selector values.String) values.String {
	if !el.IsConnected() {
		return values.EmptyString
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	found, err := el.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(el.id.nodeID, selector.String()))

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.EmptyString
	}

	childNodeID := found.NodeID

	obj, err := el.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(childNodeID))

	if err != nil {
		el.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child element")

		return values.EmptyString
	}

	if obj.Object.ObjectID == nil {
		el.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child element")

		return values.EmptyString
	}

	objID := *obj.Object.ObjectID

	text, err := eval.Property(ctx, el.client, objID, "innerText")

	if err != nil {
		el.logError(err).
			Str("childObjectID", string(objID)).
			Str("selector", selector.String()).
			Msg("failed to load inner text for found child element")

		return values.EmptyString
	}

	return values.NewString(text.String())
}

func (el *HTMLElement) InnerTextBySelectorAll(selector values.String) *values.Array {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	res, err := el.client.DOM.QuerySelectorAll(ctx, dom.NewQuerySelectorAllArgs(el.id.nodeID, selector.String()))

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.NewArray(0)
	}

	arr := values.NewArray(len(res.NodeIDs))

	for idx, id := range res.NodeIDs {
		obj, err := el.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id))

		if err != nil {
			el.logError(err).
				Int("index", idx).
				Int("childNodeID", int(id)).
				Str("selector", selector.String()).
				Msg("failed to resolve remote object for child element")

			continue
		}

		if obj.Object.ObjectID == nil {
			continue
		}

		objID := *obj.Object.ObjectID

		text, err := eval.Property(ctx, el.client, objID, "innerText")

		if err != nil {
			el.logError(err).
				Str("childObjectID", string(objID)).
				Str("selector", selector.String()).
				Msg("failed to load inner text for found child element")

			continue
		}

		arr.Push(text)
	}

	return arr
}

func (el *HTMLElement) InnerHTML() values.String {
	el.Lock()
	defer el.Unlock()

	return el.innerHTML
}

func (el *HTMLElement) InnerHTMLBySelector(selector values.String) values.String {
	if !el.IsConnected() {
		return values.EmptyString
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	found, err := el.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(el.id.nodeID, selector.String()))

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.EmptyString
	}

	text, err := loadInnerHTML(ctx, el.client, &HTMLElementIdentity{
		nodeID: found.NodeID,
	})

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to load inner HTML for found child element")

		return values.EmptyString
	}

	return text
}

func (el *HTMLElement) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(el.id.nodeID, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.NewArray(0)
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		text, err := loadInnerHTML(ctx, el.client, &HTMLElementIdentity{
			nodeID: id,
		})

		if err != nil {
			el.logError(err).
				Str("selector", selector.String()).
				Msg("failed to load inner HTML for found child element")

			// return what we have
			return arr
		}

		arr.Push(text)
	}

	return arr
}

func (el *HTMLElement) Click() (values.Boolean, error) {
	ctx, cancel := contextWithTimeout()

	defer cancel()

	return events.DispatchEvent(ctx, el.client, el.id.objectID, "click")
}

func (el *HTMLElement) Input(value core.Value, delay values.Int) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	if err := el.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(el.id.objectID)); err != nil {
		el.logError(err).Msg("failed to focus")

		return err
	}

	delayMs := time.Duration(delay)

	time.Sleep(delayMs * time.Millisecond)

	valStr := value.String()

	for _, ch := range valStr {
		for _, ev := range []string{"keyDown", "keyUp"} {
			ke := input.NewDispatchKeyEventArgs(ev).SetText(string(ch))

			if err := el.client.Input.DispatchKeyEvent(ctx, ke); err != nil {
				el.logError(err).Str("value", value.String()).Msg("failed to input a value")

				return err
			}

			time.Sleep(delayMs * time.Millisecond)
		}
	}

	return nil
}

func (el *HTMLElement) IsConnected() values.Boolean {
	el.Lock()
	defer el.Unlock()

	return el.connected
}

func (el *HTMLElement) loadInnerText() (core.Value, error) {
	if el.IsConnected() {
		ctx, cancel := contextWithTimeout()
		defer cancel()

		text, err := eval.Property(ctx, el.client, el.id.objectID, "innerText")

		if err == nil {
			return text, nil
		}

		el.logError(err).Msg("failed to read 'innerText' property of remote object")

		// and just parse innerHTML
	}

	h := el.InnerHTML()

	if h == values.EmptyString {
		return h, nil
	}

	parsed, err := parseInnerText(h.String())

	if err != nil {
		el.logError(err).Msg("failed to parse inner html")

		return values.EmptyString, err
	}

	return parsed, nil
}

func (el *HTMLElement) loadAttrs() (core.Value, error) {
	return parseAttrs(el.rawAttrs), nil
}

func (el *HTMLElement) loadChildren() (core.Value, error) {
	if !el.IsConnected() {
		return values.NewArray(0), nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	loaded := values.NewArray(len(el.children))

	for _, childId := range el.children {
		child, err := LoadElement(
			ctx,
			el.logger,
			el.client,
			el.events,
			childId.nodeID,
			childId.backendID,
		)

		if err != nil {
			el.logError(err).Msg("failed to load child nodes")

			continue
		}

		loaded.Push(child)
	}

	return loaded, nil
}

func (el *HTMLElement) handlePageReload(_ interface{}) {
	el.Close()
}

func (el *HTMLElement) handleAttrModified(message interface{}) {
	reply, ok := message.(*dom.AttributeModifiedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != el.id.nodeID {
		return
	}

	el.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update node")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Set(values.NewString(reply.Name), values.NewString(reply.Value))
	})
}

func (el *HTMLElement) handleAttrRemoved(message interface{}) {
	reply, ok := message.(*dom.AttributeRemovedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != el.id.nodeID {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	el.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update node")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Remove(values.NewString(reply.Name))
	})
}

func (el *HTMLElement) handleChildrenCountChanged(message interface{}) {
	reply, ok := message.(*dom.ChildNodeCountUpdatedReply)

	if !ok {
		return
	}

	if reply.NodeID != el.id.nodeID {
		return
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	node, err := el.client.DOM.DescribeNode(
		ctx,
		dom.NewDescribeNodeArgs().SetObjectID(el.id.objectID),
	)

	if err != nil {
		el.logError(err).Msg("failed to update node")

		return
	}

	el.Lock()
	defer el.Unlock()

	el.children = createChildrenArray(node.Node.Children)
}

func (el *HTMLElement) handleChildInserted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeInsertedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id.nodeID {
		return
	}

	targetIDx := -1
	prevID := reply.PreviousNodeID
	nextID := reply.Node.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id.nodeID == prevID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	nextIdentity := &HTMLElementIdentity{
		nodeID:    reply.Node.NodeID,
		backendID: reply.Node.BackendNodeID,
	}

	arr := el.children
	el.children = append(arr[:targetIDx], append([]*HTMLElementIdentity{nextIdentity}, arr[targetIDx:]...)...)

	if !el.loadedChildren.Ready() {
		return
	}

	el.loadedChildren.Write(func(v core.Value, err error) {
		ctx, cancel := contextWithTimeout()
		defer cancel()

		loadedArr := v.(*values.Array)
		loadedEl, err := LoadElement(ctx, el.logger, el.client, el.events, nextID, emptyBackendID)

		if err != nil {
			el.logError(err).Msg("failed to load an inserted node")

			return
		}

		loadedArr.Insert(values.NewInt(targetIDx), loadedEl)

		newInnerHTML, err := loadInnerHTML(ctx, el.client, el.id)

		if err != nil {
			el.logError(err).Msg("failed to update node")

			return
		}

		el.innerHTML = newInnerHTML
		el.innerText.Reset()
	})
}

func (el *HTMLElement) handleChildDeleted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeRemovedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id.nodeID {
		return
	}

	targetIDx := -1
	targetID := reply.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id.nodeID == targetID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	arr := el.children
	el.children = append(arr[:targetIDx], arr[targetIDx+1:]...)

	if !el.loadedChildren.Ready() {
		return
	}

	el.loadedChildren.Write(func(v core.Value, err error) {
		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(el.id.nodeID)).
				Msg("failed to update node")

			return
		}

		ctx, cancel := contextWithTimeout()
		defer cancel()

		loadedArr := v.(*values.Array)
		loadedArr.RemoveAt(values.NewInt(targetIDx))

		newInnerHTML, err := loadInnerHTML(ctx, el.client, el.id)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(el.id.nodeID)).
				Msg("failed to update node")

			return
		}

		el.innerHTML = newInnerHTML
		el.innerText.Reset()
	})
}

func (el *HTMLElement) logError(err error) *zerolog.Event {
	return el.logger.
		Error().
		Timestamp().
		Int("nodeID", int(el.id.nodeID)).
		Int("backendID", int(el.id.backendID)).
		Str("objectID", string(el.id.objectID)).
		Err(err)
}
