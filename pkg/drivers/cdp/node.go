package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/gofrs/uuid"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/rs/zerolog"
)

const DefaultTimeout = time.Second * 30

var emptyNodeID = dom.NodeID(0)
var emptyBackendID = dom.BackendNodeID(0)

type (
	HTMLNodeIdentity struct {
		nodeID    dom.NodeID
		backendID dom.BackendNodeID
		objectID  runtime.RemoteObjectID
	}

	HTMLNode struct {
		mu             sync.Mutex
		logger         *zerolog.Logger
		client         *cdp.Client
		events         *events.EventBroker
		connected      values.Boolean
		id             *HTMLNodeIdentity
		nodeType       values.Int
		nodeName       values.String
		innerHTML      values.String
		innerText      *common.LazyValue
		value          core.Value
		rawAttrs       []string
		attributes     *common.LazyValue
		children       []*HTMLNodeIdentity
		loadedChildren *common.LazyValue
	}
)

func LoadNode(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	nodeID dom.NodeID,
	backendID dom.BackendNodeID,
) (*HTMLNode, error) {
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

	id := new(HTMLNodeIdentity)
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

	return NewHTMLNode(
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

func NewHTMLNode(
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	id *HTMLNodeIdentity,
	nodeType int,
	nodeName string,
	attributes []string,
	value string,
	innerHTML values.String,
	children []*HTMLNodeIdentity,
) *HTMLNode {
	el := new(HTMLNode)
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

	broker.AddEventListener(events.EventReload, el.handlePageReload)
	broker.AddEventListener(events.EventAttrModified, el.handleAttrModified)
	broker.AddEventListener(events.EventAttrRemoved, el.handleAttrRemoved)
	broker.AddEventListener(events.EventChildNodeCountUpdated, el.handleChildrenCountChanged)
	broker.AddEventListener(events.EventChildNodeInserted, el.handleChildInserted)
	broker.AddEventListener(events.EventChildNodeRemoved, el.handleChildRemoved)

	return el
}

func (nd *HTMLNode) Close() error {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	// already closed
	if !nd.connected {
		return nil
	}

	nd.connected = values.False
	nd.events.RemoveEventListener(events.EventReload, nd.handlePageReload)
	nd.events.RemoveEventListener(events.EventAttrModified, nd.handleAttrModified)
	nd.events.RemoveEventListener(events.EventAttrRemoved, nd.handleAttrRemoved)
	nd.events.RemoveEventListener(events.EventChildNodeCountUpdated, nd.handleChildrenCountChanged)
	nd.events.RemoveEventListener(events.EventChildNodeInserted, nd.handleChildInserted)
	nd.events.RemoveEventListener(events.EventChildNodeRemoved, nd.handleChildRemoved)

	return nil
}

func (nd *HTMLNode) Type() core.Type {
	return drivers.DHTMLNodeType
}

func (nd *HTMLNode) MarshalJSON() ([]byte, error) {
	val, err := nd.innerText.Read()

	if err != nil {
		return nil, err
	}

	return json.Marshal(val.String())
}

func (nd *HTMLNode) String() string {
	return nd.InnerHTML().String()
}

func (nd *HTMLNode) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.DHTMLNodeType:
		other := other.(drivers.DHTMLNode)

		return nd.InnerHTML().Compare(other.InnerHTML())
	default:
		return drivers.Compare(nd.Type(), other.Type())
	}
}

func (nd *HTMLNode) Unwrap() interface{} {
	return nd
}

func (nd *HTMLNode) Hash() uint64 {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	h := fnv.New64a()

	h.Write([]byte(nd.Type().Name()))
	h.Write([]byte(":"))
	h.Write([]byte(nd.innerHTML))

	return h.Sum64()
}

func (nd *HTMLNode) Copy() core.Value {
	return values.None
}

func (nd *HTMLNode) Iterate(_ context.Context) (collections.CollectionIterator, error) {
	return common.NewIterator(nd)
}

func (nd *HTMLNode) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetIn(ctx, nd, path)
}

func (nd *HTMLNode) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetIn(ctx, nd, path, value)
}

func (nd *HTMLNode) GetValue() core.Value {
	if !nd.IsConnected() {
		return nd.value
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	val, err := eval.Property(ctx, nd.client, nd.id.objectID, "value")

	if err != nil {
		nd.logError(err).Msg("failed to get node value")

		return nd.value
	}

	nd.value = val

	return val
}

func (nd *HTMLNode) SetValue(value core.Value) error {
	if !nd.IsConnected() {
		// TODO: Return an error
		return nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	return nd.client.DOM.SetNodeValue(ctx, dom.NewSetNodeValueArgs(nd.id.nodeID, value.String()))
}

func (nd *HTMLNode) NodeType() values.Int {
	return nd.nodeType
}

func (nd *HTMLNode) NodeName() values.String {
	return nd.nodeName
}

func (nd *HTMLNode) Length() values.Int {
	return values.NewInt(len(nd.children))
}

func (nd *HTMLNode) GetAttributes() core.Value {
	val, err := nd.attributes.Read()

	if err != nil {
		return values.None
	}

	// returning shallow copy
	return val.Copy()
}

func (nd *HTMLNode) GetAttribute(name values.String) core.Value {
	attrs, err := nd.attributes.Read()

	if err != nil {
		return values.None
	}

	val, found := attrs.(*values.Object).Get(name)

	if !found {
		return values.None
	}

	return val
}

func (nd *HTMLNode) SetAttribute(name, value values.String) error {
	return nd.client.DOM.SetAttributeValue(
		context.Background(),
		dom.NewSetAttributeValueArgs(nd.id.nodeID, string(name), string(value)),
	)
}

func (nd *HTMLNode) GetChildNodes() core.Value {
	val, err := nd.loadedChildren.Read()

	if err != nil {
		return values.NewArray(0)
	}

	return val
}

func (nd *HTMLNode) GetChildNode(idx values.Int) core.Value {
	// TODO: Add lazy loading
	val, err := nd.loadedChildren.Read()

	if err != nil {
		return values.None
	}

	return val.(*values.Array).Get(idx)
}

func (nd *HTMLNode) QuerySelector(selector values.String) core.Value {
	if !nd.IsConnected() {
		return values.None
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(nd.id.nodeID, selector.String())
	found, err := nd.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.None
	}

	if found.NodeID == emptyNodeID {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to find a node by selector. returned 0 NodeID")

		return values.None
	}

	res, err := LoadNode(ctx, nd.logger, nd.client, nd.events, found.NodeID, emptyBackendID)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to load a child node by selector")

		return values.None
	}

	return res
}

func (nd *HTMLNode) QuerySelectorAll(selector values.String) core.Value {
	if !nd.IsConnected() {
		return values.NewArray(0)
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(nd.id.nodeID, selector.String())
	res, err := nd.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.None
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		if id == emptyNodeID {
			nd.logError(err).
				Str("selector", selector.String()).
				Msg("failed to find a node by selector. returned 0 NodeID")

			continue
		}

		childEl, err := LoadNode(ctx, nd.logger, nd.client, nd.events, id, emptyBackendID)

		if err != nil {
			nd.logError(err).
				Str("selector", selector.String()).
				Msg("failed to load nodes by selector")

			// close elements that are already loaded, but won't be used because of the error
			if arr.Length() > 0 {
				arr.ForEach(func(e core.Value, _ int) bool {
					e.(*HTMLNode).Close()

					return true
				})
			}

			return values.None
		}

		arr.Push(childEl)
	}

	return arr
}

func (nd *HTMLNode) InnerText() values.String {
	val, err := nd.innerText.Read()

	if err != nil {
		return values.EmptyString
	}

	if val == values.None {
		return values.EmptyString
	}

	return val.(values.String)
}

func (nd *HTMLNode) InnerTextBySelector(selector values.String) values.String {
	if !nd.IsConnected() {
		return values.EmptyString
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	found, err := nd.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(nd.id.nodeID, selector.String()))

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.EmptyString
	}

	if found.NodeID == emptyNodeID {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to find a node by selector. returned 0 NodeID")

		return values.EmptyString
	}

	childNodeID := found.NodeID

	obj, err := nd.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(childNodeID))

	if err != nil {
		nd.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child element")

		return values.EmptyString
	}

	if obj.Object.ObjectID == nil {
		nd.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child element")

		return values.EmptyString
	}

	objID := *obj.Object.ObjectID

	text, err := eval.Property(ctx, nd.client, objID, "innerText")

	if err != nil {
		nd.logError(err).
			Str("childObjectID", string(objID)).
			Str("selector", selector.String()).
			Msg("failed to load inner text for found child element")

		return values.EmptyString
	}

	return values.NewString(text.String())
}

func (nd *HTMLNode) InnerTextBySelectorAll(selector values.String) *values.Array {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	res, err := nd.client.DOM.QuerySelectorAll(ctx, dom.NewQuerySelectorAllArgs(nd.id.nodeID, selector.String()))

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.NewArray(0)
	}

	arr := values.NewArray(len(res.NodeIDs))

	for idx, id := range res.NodeIDs {
		if id == emptyNodeID {
			nd.logError(err).
				Str("selector", selector.String()).
				Msg("failed to find a node by selector. returned 0 NodeID")

			continue
		}

		obj, err := nd.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id))

		if err != nil {
			nd.logError(err).
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

		text, err := eval.Property(ctx, nd.client, objID, "innerText")

		if err != nil {
			nd.logError(err).
				Str("childObjectID", string(objID)).
				Str("selector", selector.String()).
				Msg("failed to load inner text for found child element")

			continue
		}

		arr.Push(text)
	}

	return arr
}

func (nd *HTMLNode) InnerHTML() values.String {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	return nd.innerHTML
}

func (nd *HTMLNode) InnerHTMLBySelector(selector values.String) values.String {
	if !nd.IsConnected() {
		return values.EmptyString
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	found, err := nd.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(nd.id.nodeID, selector.String()))

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.EmptyString
	}

	text, err := loadInnerHTML(ctx, nd.client, &HTMLNodeIdentity{
		nodeID: found.NodeID,
	})

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to load inner HTML for found child element")

		return values.EmptyString
	}

	return text
}

func (nd *HTMLNode) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(nd.id.nodeID, selector.String())
	res, err := nd.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.NewArray(0)
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		text, err := loadInnerHTML(ctx, nd.client, &HTMLNodeIdentity{
			nodeID: id,
		})

		if err != nil {
			nd.logError(err).
				Str("selector", selector.String()).
				Msg("failed to load inner HTML for found child element")

			// return what we have
			return arr
		}

		arr.Push(text)
	}

	return arr
}

func (nd *HTMLNode) CountBySelector(selector values.String) values.Int {
	if !nd.IsConnected() {
		return values.ZeroInt
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(nd.id.nodeID, selector.String())
	res, err := nd.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.ZeroInt
	}

	return values.NewInt(len(res.NodeIDs))
}

func (nd *HTMLNode) ExistsBySelector(selector values.String) values.Boolean {
	if !nd.IsConnected() {
		return values.False
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(nd.id.nodeID, selector.String())
	res, err := nd.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		nd.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.False
	}

	if res.NodeID == 0 {
		return values.False
	}

	return values.True
}

func (nd *HTMLNode) WaitForClass(class values.String, timeout values.Int) error {
	task := events.NewWaitTask(
		func() (core.Value, error) {
			current := nd.GetAttribute("class")

			if current.Type() != types.String {
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

func (nd *HTMLNode) Click() (values.Boolean, error) {
	ctx, cancel := contextWithTimeout()

	defer cancel()

	return events.DispatchEvent(ctx, nd.client, nd.id.objectID, "click")
}

func (nd *HTMLNode) Input(value core.Value, delay values.Int) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	if err := nd.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(nd.id.objectID)); err != nil {
		nd.logError(err).Msg("failed to focus")

		return err
	}

	delayMs := time.Duration(delay)

	time.Sleep(delayMs * time.Millisecond)

	valStr := value.String()

	for _, ch := range valStr {
		for _, ev := range []string{"keyDown", "keyUp"} {
			ke := input.NewDispatchKeyEventArgs(ev).SetText(string(ch))

			if err := nd.client.Input.DispatchKeyEvent(ctx, ke); err != nil {
				nd.logError(err).Str("value", value.String()).Msg("failed to input a value")

				return err
			}

			time.Sleep(delayMs * time.Millisecond)
		}
	}

	return nil
}

func (nd *HTMLNode) Select(value *values.Array) (*values.Array, error) {
	var attrID = "data-ferret-select"

	if nd.NodeName() != "SELECT" {
		return nil, core.Error(core.ErrInvalidOperation, "Element is not a <select> element.")
	}

	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	err = nd.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(nd.id.nodeID, attrID, id.String()))

	if err != nil {
		return nil, err
	}

	res, err := eval.Eval(
		nd.client,
		fmt.Sprintf(`
			var element = document.querySelector('[%s="%s"]');

			if (element == null) {
				return [];
			}

			var values = %s;

			if (element.nodeName.toLowerCase() !== 'select') {
				throw new Error('Element is not a <select> element.');
			}

			var options = Array.from(element.options);
      		element.value = undefined;

			for (var option of options) {
        		option.selected = values.includes(option.value);
        	
				if (option.selected && !element.multiple) {
          			break;
				}
      		}

      		element.dispatchEvent(new Event('input', { 'bubbles': true }));
      		element.dispatchEvent(new Event('change', { 'bubbles': true }));
      		
			return options.filter(option => option.selected).map(option => option.value);
		`,
			attrID,
			id.String(),
			value.String(),
		),
		true,
		false,
	)

	nd.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(nd.id.nodeID, attrID))

	if err != nil {
		return nil, err
	}

	arr, ok := res.(*values.Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(types.Array, res.Type())
}

func (nd *HTMLNode) ScrollIntoView() error {
	var attrID = "data-ferret-scroll"

	id, err := uuid.NewV4()

	if err != nil {
		return err
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	err = nd.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(nd.id.nodeID, attrID, id.String()))

	if err != nil {
		return err
	}

	_, err = eval.Eval(nd.client, fmt.Sprintf(`
		var nd = document.querySelector('[%s="%s"]');

		if (nd == null) {
			throw new Error('element not found');
		}

		nd.scrollIntoView({
    		behavior: 'instant',
			inline: 'center',
			block: 'center'
  		});
	`,
		attrID,
		id.String(),
	), false, false)

	nd.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(nd.id.nodeID, attrID))

	return err
}

func (nd *HTMLNode) Hover() error {
	err := nd.ScrollIntoView()

	if err != nil {
		return err
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	q, err := getClickablePoint(ctx, nd.client, nd.id)

	if err != nil {
		return err
	}

	return nd.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", q.X, q.Y),
	)
}

func (nd *HTMLNode) IsConnected() values.Boolean {
	nd.mu.Lock()
	defer nd.mu.Unlock()

	return nd.connected
}

func (nd *HTMLNode) loadInnerText() (core.Value, error) {
	if nd.IsConnected() {
		ctx, cancel := contextWithTimeout()
		defer cancel()

		text, err := loadInnerText(ctx, nd.client, nd.id)

		if err == nil {
			return text, nil
		}

		nd.logError(err).Msg("failed to get get inner text from remote object")

		// and just parse cached innerHTML
	}

	h := nd.InnerHTML()

	if h == values.EmptyString {
		return h, nil
	}

	parsed, err := parseInnerText(h.String())

	if err != nil {
		nd.logError(err).Msg("failed to parse inner html")

		return values.EmptyString, err
	}

	return parsed, nil
}

func (nd *HTMLNode) loadAttrs() (core.Value, error) {
	return parseAttrs(nd.rawAttrs), nil
}

func (nd *HTMLNode) loadChildren() (core.Value, error) {
	if !nd.IsConnected() {
		return values.NewArray(0), nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	loaded := values.NewArray(len(nd.children))

	for _, childID := range nd.children {
		child, err := LoadNode(
			ctx,
			nd.logger,
			nd.client,
			nd.events,
			childID.nodeID,
			childID.backendID,
		)

		if err != nil {
			nd.logError(err).Msg("failed to load child nodes")

			continue
		}

		loaded.Push(child)
	}

	return loaded, nil
}

func (nd *HTMLNode) handlePageReload(_ interface{}) {
	nd.Close()
}

func (nd *HTMLNode) handleAttrModified(message interface{}) {
	reply, ok := message.(*dom.AttributeModifiedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != nd.id.nodeID {
		return
	}

	nd.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			nd.logError(err).Msg("failed to update node")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Set(values.NewString(reply.Name), values.NewString(reply.Value))
	})
}

func (nd *HTMLNode) handleAttrRemoved(message interface{}) {
	reply, ok := message.(*dom.AttributeRemovedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != nd.id.nodeID {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !nd.attributes.Ready() {
		return
	}

	nd.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			nd.logError(err).Msg("failed to update node")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Remove(values.NewString(reply.Name))
	})
}

func (nd *HTMLNode) handleChildrenCountChanged(message interface{}) {
	reply, ok := message.(*dom.ChildNodeCountUpdatedReply)

	if !ok {
		return
	}

	if reply.NodeID != nd.id.nodeID {
		return
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	node, err := nd.client.DOM.DescribeNode(
		ctx,
		dom.NewDescribeNodeArgs().SetObjectID(nd.id.objectID),
	)

	if err != nil {
		nd.logError(err).Msg("failed to update node")

		return
	}

	nd.mu.Lock()
	defer nd.mu.Unlock()

	nd.children = createChildrenArray(node.Node.Children)
}

func (nd *HTMLNode) handleChildInserted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeInsertedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != nd.id.nodeID {
		return
	}

	targetIDx := -1
	prevID := reply.PreviousNodeID
	nextID := reply.Node.NodeID

	nd.mu.Lock()
	defer nd.mu.Unlock()

	for idx, id := range nd.children {
		if id.nodeID == prevID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	nextIdentity := &HTMLNodeIdentity{
		nodeID:    reply.Node.NodeID,
		backendID: reply.Node.BackendNodeID,
	}

	arr := nd.children
	nd.children = append(arr[:targetIDx], append([]*HTMLNodeIdentity{nextIdentity}, arr[targetIDx:]...)...)

	if !nd.loadedChildren.Ready() {
		return
	}

	nd.loadedChildren.Write(func(v core.Value, err error) {
		ctx, cancel := contextWithTimeout()
		defer cancel()

		loadedArr := v.(*values.Array)
		loadedEl, err := LoadNode(ctx, nd.logger, nd.client, nd.events, nextID, emptyBackendID)

		if err != nil {
			nd.logError(err).Msg("failed to load an inserted node")

			return
		}

		loadedArr.Insert(values.NewInt(targetIDx), loadedEl)

		newInnerHTML, err := loadInnerHTML(ctx, nd.client, nd.id)

		if err != nil {
			nd.logError(err).Msg("failed to update node")

			return
		}

		nd.innerHTML = newInnerHTML
		nd.innerText.Reset()
	})
}

func (nd *HTMLNode) handleChildRemoved(message interface{}) {
	reply, ok := message.(*dom.ChildNodeRemovedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != nd.id.nodeID {
		return
	}

	targetIDx := -1
	targetID := reply.NodeID

	nd.mu.Lock()
	defer nd.mu.Unlock()

	for idx, id := range nd.children {
		if id.nodeID == targetID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	arr := nd.children
	nd.children = append(arr[:targetIDx], arr[targetIDx+1:]...)

	if !nd.loadedChildren.Ready() {
		return
	}

	nd.loadedChildren.Write(func(v core.Value, err error) {
		if err != nil {
			nd.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(nd.id.nodeID)).
				Msg("failed to update node")

			return
		}

		ctx, cancel := contextWithTimeout()
		defer cancel()

		loadedArr := v.(*values.Array)
		loadedArr.RemoveAt(values.NewInt(targetIDx))

		newInnerHTML, err := loadInnerHTML(ctx, nd.client, nd.id)

		if err != nil {
			nd.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(nd.id.nodeID)).
				Msg("failed to update node")

			return
		}

		nd.innerHTML = newInnerHTML
		nd.innerText.Reset()
	})
}

func (nd *HTMLNode) logError(err error) *zerolog.Event {
	return nd.logger.
		Error().
		Timestamp().
		Int("nodeID", int(nd.id.nodeID)).
		Int("backendID", int(nd.id.backendID)).
		Str("objectID", string(nd.id.objectID)).
		Err(err)
}
