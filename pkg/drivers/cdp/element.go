package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/gofrs/uuid"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/rs/zerolog"
)

var emptyNodeID = dom.NodeID(0)
var emptyBackendID = dom.BackendNodeID(0)

type (
	HTMLElementIdentity struct {
		nodeID    dom.NodeID
		backendID dom.BackendNodeID
		objectID  runtime.RemoteObjectID
	}

	HTMLElement struct {
		mu             sync.Mutex
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

	broker.AddEventListener(events.EventReload, el.handlePageReload)
	broker.AddEventListener(events.EventAttrModified, el.handleAttrModified)
	broker.AddEventListener(events.EventAttrRemoved, el.handleAttrRemoved)
	broker.AddEventListener(events.EventChildNodeCountUpdated, el.handleChildrenCountChanged)
	broker.AddEventListener(events.EventChildNodeInserted, el.handleChildInserted)
	broker.AddEventListener(events.EventChildNodeRemoved, el.handleChildRemoved)

	return el
}

func (el *HTMLElement) Close() error {
	el.mu.Lock()
	defer el.mu.Unlock()

	// already closed
	if !el.connected {
		return nil
	}

	el.connected = values.False
	el.events.RemoveEventListener(events.EventReload, el.handlePageReload)
	el.events.RemoveEventListener(events.EventAttrModified, el.handleAttrModified)
	el.events.RemoveEventListener(events.EventAttrRemoved, el.handleAttrRemoved)
	el.events.RemoveEventListener(events.EventChildNodeCountUpdated, el.handleChildrenCountChanged)
	el.events.RemoveEventListener(events.EventChildNodeInserted, el.handleChildInserted)
	el.events.RemoveEventListener(events.EventChildNodeRemoved, el.handleChildRemoved)

	return nil
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	val, err := el.innerText.Read(context.Background())

	if err != nil {
		return nil, err
	}

	return json.Marshal(val.String())
}

func (el *HTMLElement) String() string {
	return el.InnerHTML(context.Background()).String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		ctx := context.Background()

		return el.InnerHTML(ctx).Compare(other.InnerHTML(ctx))
	default:
		return drivers.Compare(el.Type(), other.Type())
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el
}

func (el *HTMLElement) Hash() uint64 {
	el.mu.Lock()
	defer el.mu.Unlock()

	h := fnv.New64a()

	h.Write([]byte(el.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(el.innerHTML))

	return h.Sum64()
}

func (el *HTMLElement) Copy() core.Value {
	return values.None
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInElement(ctx, el, path)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInElement(ctx, el, path, value)
}

func (el *HTMLElement) GetValue(ctx context.Context) core.Value {
	if !el.IsConnected() {
		return el.value
	}

	val, err := eval.Property(ctx, el.client, el.id.objectID, "value")

	if err != nil {
		el.logError(err).Msg("failed to get node value")

		return el.value
	}

	el.value = val

	return val
}

func (el *HTMLElement) SetValue(ctx context.Context, value core.Value) error {
	if !el.IsConnected() {
		// TODO: Return an error
		return nil
	}

	return el.client.DOM.SetNodeValue(ctx, dom.NewSetNodeValueArgs(el.id.nodeID, value.String()))
}

func (el *HTMLElement) NodeType() values.Int {
	return el.nodeType
}

func (el *HTMLElement) NodeName() values.String {
	return el.nodeName
}

func (el *HTMLElement) Length() values.Int {
	return values.NewInt(len(el.children))
}

func (el *HTMLElement) GetAttributes(ctx context.Context) *values.Object {
	val, err := el.attributes.Read(ctx)

	if err != nil {
		return values.NewObject()
	}

	attrs := val.(*values.Object)

	// returning shallow copy
	return attrs.Copy().(*values.Object)
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name values.String) core.Value {
	attrs, err := el.attributes.Read(ctx)

	if err != nil {
		return values.None
	}

	val, found := attrs.(*values.Object).Get(name)

	if !found {
		return values.None
	}

	return val
}

func (el *HTMLElement) SetAttribute(ctx context.Context, name, value values.String) error {
	return el.client.DOM.SetAttributeValue(
		ctx,
		dom.NewSetAttributeValueArgs(el.id.nodeID, string(name), string(value)),
	)
}

func (el *HTMLElement) GetChildNodes(ctx context.Context) core.Value {
	val, err := el.loadedChildren.Read(ctx)

	if err != nil {
		return values.NewArray(0)
	}

	return val
}

func (el *HTMLElement) GetChildNode(ctx context.Context, idx values.Int) core.Value {
	val, err := el.loadedChildren.Read(ctx)

	if err != nil {
		return values.None
	}

	return val.(*values.Array).Get(idx)
}

func (el *HTMLElement) QuerySelector(ctx context.Context, selector values.String) core.Value {
	if !el.IsConnected() {
		return values.None
	}

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(el.id.nodeID, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.None
	}

	if found.NodeID == emptyNodeID {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to find a node by selector. returned 0 NodeID")

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

func (el *HTMLElement) QuerySelectorAll(ctx context.Context, selector values.String) core.Value {
	if !el.IsConnected() {
		return values.NewArray(0)
	}

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
		if id == emptyNodeID {
			el.logError(err).
				Str("selector", selector.String()).
				Msg("failed to find a node by selector. returned 0 NodeID")

			continue
		}

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

func (el *HTMLElement) InnerText(ctx context.Context) values.String {
	val, err := el.innerText.Read(ctx)

	if err != nil {
		return values.EmptyString
	}

	if val == values.None {
		return values.EmptyString
	}

	return val.(values.String)
}

func (el *HTMLElement) InnerTextBySelector(ctx context.Context, selector values.String) values.String {
	if !el.IsConnected() {
		return values.EmptyString
	}

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	found, err := el.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(el.id.nodeID, selector.String()))

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.EmptyString
	}

	if found.NodeID == emptyNodeID {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to find a node by selector. returned 0 NodeID")

		return values.EmptyString
	}

	childNodeID := found.NodeID

	obj, err := el.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(childNodeID))

	if err != nil {
		el.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child el")

		return values.EmptyString
	}

	if obj.Object.ObjectID == nil {
		el.logError(err).
			Int("childNodeID", int(childNodeID)).
			Str("selector", selector.String()).
			Msg("failed to resolve remote object for child el")

		return values.EmptyString
	}

	objID := *obj.Object.ObjectID

	text, err := eval.Property(ctx, el.client, objID, "innerText")

	if err != nil {
		el.logError(err).
			Str("childObjectID", string(objID)).
			Str("selector", selector.String()).
			Msg("failed to load inner text for found child el")

		return values.EmptyString
	}

	return values.NewString(text.String())
}

func (el *HTMLElement) InnerTextBySelectorAll(ctx context.Context, selector values.String) *values.Array {
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
		if id == emptyNodeID {
			el.logError(err).
				Str("selector", selector.String()).
				Msg("failed to find a node by selector. returned 0 NodeID")

			continue
		}

		obj, err := el.client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id))

		if err != nil {
			el.logError(err).
				Int("index", idx).
				Int("childNodeID", int(id)).
				Str("selector", selector.String()).
				Msg("failed to resolve remote object for child el")

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
				Msg("failed to load inner text for found child el")

			continue
		}

		arr.Push(text)
	}

	return arr
}

func (el *HTMLElement) InnerHTML(_ context.Context) values.String {
	el.mu.Lock()
	defer el.mu.Unlock()

	return el.innerHTML
}

func (el *HTMLElement) InnerHTMLBySelector(ctx context.Context, selector values.String) values.String {
	if !el.IsConnected() {
		return values.EmptyString
	}

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
			Msg("failed to load inner HTML for found child el")

		return values.EmptyString
	}

	return text
}

func (el *HTMLElement) InnerHTMLBySelectorAll(ctx context.Context, selector values.String) *values.Array {
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
				Msg("failed to load inner HTML for found child el")

			// return what we have
			return arr
		}

		arr.Push(text)
	}

	return arr
}

func (el *HTMLElement) CountBySelector(ctx context.Context, selector values.String) values.Int {
	if !el.IsConnected() {
		return values.ZeroInt
	}

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorAllArgs(el.id.nodeID, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.ZeroInt
	}

	return values.NewInt(len(res.NodeIDs))
}

func (el *HTMLElement) ExistsBySelector(ctx context.Context, selector values.String) values.Boolean {
	if !el.IsConnected() {
		return values.False
	}

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(el.id.nodeID, selector.String())
	res, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		el.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.False
	}

	if res.NodeID == 0 {
		return values.False
	}

	return values.True
}

func (el *HTMLElement) WaitForClass(ctx context.Context, class values.String, when drivers.WaitEvent) error {
	task := events.NewWaitTask(
		func(ctx2 context.Context) (core.Value, error) {
			current := el.GetAttribute(ctx2, "class")

			if current.Type() != types.String {
				return values.None, nil
			}

			str := current.(values.String)
			classStr := string(class)
			classes := strings.Split(string(str), " ")

			if when != drivers.WaitEventAbsence {
				for _, c := range classes {
					if c == classStr {
						// The value does not really matter if it's not None
						// None indicates that operation needs to be repeated
						return values.True, nil
					}
				}
			} else {
				var found values.Boolean

				for _, c := range classes {
					if c == classStr {
						found = values.True
						break
					}
				}

				if found == values.False {
					// The value does not really matter if it's not None
					// None indicates that operation needs to be repeated
					return values.False, nil
				}
			}

			return values.None, nil
		},
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) Click(ctx context.Context) (values.Boolean, error) {
	return events.DispatchEvent(ctx, el.client, el.id.objectID, "click")
}

func (el *HTMLElement) Input(ctx context.Context, value core.Value, delay values.Int) error {
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

func (el *HTMLElement) Select(ctx context.Context, value *values.Array) (*values.Array, error) {
	var attrID = "data-ferret-select"

	if el.NodeName() != "SELECT" {
		return nil, core.Error(core.ErrInvalidOperation, "element is not a <select> element.")
	}

	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	err = el.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(el.id.nodeID, attrID, id.String()))

	if err != nil {
		return nil, err
	}

	res, err := eval.Eval(
		ctx,
		el.client,
		fmt.Sprintf(`
			var el = document.querySelector('[%s="%s"]');
			if (el == null) {
				return [];
			}
			var values = %s;
			if (el.nodeName.toLowerCase() !== 'select') {
				throw new Error('element is not a <select> element.');
			}
			var options = Array.from(el.options);
      		el.value = undefined;
			for (var option of options) {
        		option.selected = values.includes(option.value);
        	
				if (option.selected && !el.multiple) {
          			break;
				}
      		}
      		el.dispatchEvent(new Event('input', { 'bubbles': true }));
      		el.dispatchEvent(new Event('change', { 'bubbles': true }));
      		
			return options.filter(option => option.selected).map(option => option.value);
		`,
			attrID,
			id.String(),
			value.String(),
		),
		true,
		false,
	)

	el.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(el.id.nodeID, attrID))

	if err != nil {
		return nil, err
	}

	arr, ok := res.(*values.Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(types.Array, res.Type())
}

func (el *HTMLElement) ScrollIntoView(ctx context.Context) error {
	var attrID = "data-ferret-scroll"

	id, err := uuid.NewV4()

	if err != nil {
		return err
	}

	err = el.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(el.id.nodeID, attrID, id.String()))

	if err != nil {
		return err
	}

	_, err = eval.Eval(
		ctx,
		el.client,
		fmt.Sprintf(`
			var el = document.querySelector('[%s="%s"]');
			if (el == null) {
				throw new Error('element not found');
			}
			
			el.scrollIntoView({
    			behavior: 'instant',
				inline: 'center',
				block: 'center'
  			});
		`,
			attrID,
			id.String(),
		), false, false)

	el.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(el.id.nodeID, attrID))

	return err
}

func (el *HTMLElement) Hover(ctx context.Context) error {
	err := el.ScrollIntoView(ctx)

	if err != nil {
		return err
	}

	q, err := getClickablePoint(ctx, el.client, el.id)

	if err != nil {
		return err
	}

	return el.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", q.X, q.Y),
	)
}

func (el *HTMLElement) IsConnected() values.Boolean {
	el.mu.Lock()
	defer el.mu.Unlock()

	return el.connected
}

func (el *HTMLElement) loadInnerText(ctx context.Context) (core.Value, error) {
	if el.IsConnected() {
		text, err := loadInnerText(ctx, el.client, el.id)

		if err == nil {
			return text, nil
		}

		el.logError(err).Msg("failed to get inner text from remote object")

		// and just parse cached innerHTML
	}

	h := el.InnerHTML(ctx)

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

func (el *HTMLElement) loadAttrs(_ context.Context) (core.Value, error) {
	return parseAttrs(el.rawAttrs), nil
}

func (el *HTMLElement) loadChildren(ctx context.Context) (core.Value, error) {
	if !el.IsConnected() {
		return values.NewArray(0), nil
	}

	loaded := values.NewArray(len(el.children))

	for _, childID := range el.children {
		child, err := LoadElement(
			ctx,
			el.logger,
			el.client,
			el.events,
			childID.nodeID,
			childID.backendID,
		)

		if err != nil {
			el.logError(err).Msg("failed to load child elements")

			continue
		}

		loaded.Push(child)
	}

	return loaded, nil
}

func (el *HTMLElement) handlePageReload(_ context.Context, _ interface{}) {
	el.Close()
}

func (el *HTMLElement) handleAttrModified(ctx context.Context, message interface{}) {
	reply, ok := message.(*dom.AttributeModifiedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this el
	if reply.NodeID != el.id.nodeID {
		return
	}

	el.attributes.Write(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update element")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Set(values.NewString(reply.Name), values.NewString(reply.Value))
	})
}

func (el *HTMLElement) handleAttrRemoved(ctx context.Context, message interface{}) {
	reply, ok := message.(*dom.AttributeRemovedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this el
	if reply.NodeID != el.id.nodeID {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	el.attributes.Write(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update element")

			return
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Remove(values.NewString(reply.Name))
	})
}

func (el *HTMLElement) handleChildrenCountChanged(ctx context.Context, message interface{}) {
	reply, ok := message.(*dom.ChildNodeCountUpdatedReply)

	if !ok {
		return
	}

	if reply.NodeID != el.id.nodeID {
		return
	}

	node, err := el.client.DOM.DescribeNode(
		ctx,
		dom.NewDescribeNodeArgs().SetObjectID(el.id.objectID),
	)

	if err != nil {
		el.logError(err).Msg("failed to update element")

		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	el.children = createChildrenArray(node.Node.Children)
}

func (el *HTMLElement) handleChildInserted(ctx context.Context, message interface{}) {
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

	el.mu.Lock()
	defer el.mu.Unlock()

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

	el.loadedChildren.Write(ctx, func(v core.Value, err error) {
		loadedArr := v.(*values.Array)
		loadedEl, err := LoadElement(ctx, el.logger, el.client, el.events, nextID, emptyBackendID)

		if err != nil {
			el.logError(err).Msg("failed to load an inserted element")

			return
		}

		loadedArr.Insert(values.NewInt(targetIDx), loadedEl)

		newInnerHTML, err := loadInnerHTML(ctx, el.client, el.id)

		if err != nil {
			el.logError(err).Msg("failed to update element")

			return
		}

		el.innerHTML = newInnerHTML
		el.innerText.Reset()
	})
}

func (el *HTMLElement) handleChildRemoved(ctx context.Context, message interface{}) {
	reply, ok := message.(*dom.ChildNodeRemovedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id.nodeID {
		return
	}

	targetIDx := -1
	targetID := reply.NodeID

	el.mu.Lock()
	defer el.mu.Unlock()

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

	el.loadedChildren.Write(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(el.id.nodeID)).
				Msg("failed to update element")

			return
		}

		loadedArr := v.(*values.Array)
		loadedArr.RemoveAt(values.NewInt(targetIDx))

		newInnerHTML, err := loadInnerHTML(ctx, el.client, el.id)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(el.id.nodeID)).
				Msg("failed to update element")

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
