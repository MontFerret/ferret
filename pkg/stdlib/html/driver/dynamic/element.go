package dynamic

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/events"
	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"strconv"
	"sync"
	"time"
)

const DefaultTimeout = time.Second * 30

type HtmlElement struct {
	sync.Mutex
	client         *cdp.Client
	broker         *events.EventBroker
	connected      values.Boolean
	id             dom.NodeID
	nodeType       values.Int
	nodeName       values.String
	innerHtml      values.String
	innerText      *common.LazyValue
	value          core.Value
	attributes     *common.LazyValue
	children       []dom.NodeID
	loadedChildren *common.LazyValue
}

func LoadElement(
	client *cdp.Client,
	broker *events.EventBroker,
	id dom.NodeID,
) (*HtmlElement, error) {
	if client == nil {
		return nil, core.Error(core.ErrMissedArgument, "client")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), DefaultTimeout)

	defer cancelFn()

	node, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetNodeID(id).
			SetDepth(1),
	)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	innerHtml, err := loadInnerHtml(client, id)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	return NewHtmlElement(
		client,
		broker,
		id,
		node.Node,
		innerHtml,
	), nil
}

func NewHtmlElement(
	client *cdp.Client,
	broker *events.EventBroker,
	id dom.NodeID,
	node dom.Node,
	innerHtml values.String,
) *HtmlElement {
	el := new(HtmlElement)
	el.client = client
	el.broker = broker
	el.connected = values.True
	el.id = id
	el.nodeType = values.NewInt(node.NodeType)
	el.nodeName = values.NewString(node.NodeName)
	el.innerHtml = innerHtml
	el.innerText = common.NewLazyValue(func() (core.Value, error) {
		h := el.InnerHtml()

		if h == values.EmptyString {
			return h, nil
		}

		buff := bytes.NewBuffer([]byte(h))

		parsed, err := goquery.NewDocumentFromReader(buff)

		if err != nil {
			return values.EmptyString, err
		}

		return values.NewString(parsed.Text()), nil
	})
	el.attributes = common.NewLazyValue(func() (core.Value, error) {
		return parseAttrs(node.Attributes), nil
	})
	el.value = values.EmptyString
	el.loadedChildren = common.NewLazyValue(func() (core.Value, error) {
		if !el.IsConnected() {
			return values.NewArray(0), nil
		}

		return loadNodes(client, broker, el.children)
	})

	if node.Value != nil {
		el.value = values.NewString(*node.Value)
	}

	el.children = createChildrenArray(node.Children)

	broker.AddEventListener("reload", el.handlePageReload)
	broker.AddEventListener("attr:modified", el.handleAttrModified)
	broker.AddEventListener("attr:removed", el.handleAttrRemoved)
	broker.AddEventListener("children:count", el.handleChildrenCountChanged)
	broker.AddEventListener("children:inserted", el.handleChildInserted)
	broker.AddEventListener("children:deleted", el.handleChildDeleted)

	return el
}

func (el *HtmlElement) Close() error {
	el.Lock()
	defer el.Unlock()

	// already closed
	if !el.connected {
		return nil
	}

	el.connected = false
	el.broker.RemoveEventListener("reload", el.handlePageReload)
	el.broker.RemoveEventListener("attr:modified", el.handleAttrModified)
	el.broker.RemoveEventListener("attr:removed", el.handleAttrRemoved)
	el.broker.RemoveEventListener("children:count", el.handleChildrenCountChanged)
	el.broker.RemoveEventListener("children:inserted", el.handleChildInserted)
	el.broker.RemoveEventListener("children:deleted", el.handleChildDeleted)

	return nil
}

func (el *HtmlElement) Type() core.Type {
	return core.HtmlElementType
}

func (el *HtmlElement) MarshalJSON() ([]byte, error) {
	el.Lock()
	defer el.Unlock()

	val, err := el.innerText.Value()

	if err != nil {
		return nil, err
	}

	return json.Marshal(val.String())
}

func (el *HtmlElement) String() string {
	return el.InnerHtml().String()
}

func (el *HtmlElement) Compare(other core.Value) int {
	switch other.Type() {
	case core.HtmlDocumentType:
		other := other.(*HtmlElement)

		id := int(el.id)
		otherId := int(other.id)

		if id == otherId {
			return 0
		}

		if id > otherId {
			return 1
		}

		return -1
	default:
		if other.Type() > core.HtmlElementType {
			return -1
		}

		return 1
	}
}

func (el *HtmlElement) Unwrap() interface{} {
	return el
}

func (el *HtmlElement) Hash() int {
	h := sha512.New()

	out, err := h.Write([]byte(strconv.Itoa(int(el.id))))

	if err != nil {
		return 0
	}

	return out
}

func (el *HtmlElement) Value() core.Value {
	if !el.IsConnected() {
		return el.value
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	val, err := eval.Property(ctx, el.client, el.id, "value")

	if err != nil {
		return el.value
	}

	el.value = val

	return val
}

func (el *HtmlElement) Clone() core.Value {
	return values.None
}

func (el *HtmlElement) Length() values.Int {
	return values.NewInt(len(el.children))
}

func (el *HtmlElement) NodeType() values.Int {
	return el.nodeType
}

func (el *HtmlElement) NodeName() values.String {
	return el.nodeName
}

func (el *HtmlElement) GetAttributes() core.Value {
	val, err := el.attributes.Value()

	if err != nil {
		return values.None
	}

	// returning shallow copy
	return val.Clone()
}

func (el *HtmlElement) GetAttribute(name values.String) core.Value {
	attrs, err := el.attributes.Value()

	if err != nil {
		return values.None
	}

	val, found := attrs.(*values.Object).Get(name)

	if !found {
		return values.None
	}

	return val
}

func (el *HtmlElement) GetChildNodes() core.Value {
	val, err := el.loadedChildren.Value()

	if err != nil {
		return values.NewArray(0)
	}

	return val
}

func (el *HtmlElement) GetChildNode(idx values.Int) core.Value {
	val, err := el.loadedChildren.Value()

	if err != nil {
		return values.None
	}

	return val.(*values.Array).Get(idx)
}

func (el *HtmlElement) QuerySelector(selector values.String) core.Value {
	if !el.IsConnected() {
		return values.NewArray(0)
	}

	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorArgs(el.id, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		return values.None
	}

	res, err := LoadElement(el.client, el.broker, found.NodeID)

	if err != nil {
		return values.None
	}

	return res
}

func (el *HtmlElement) QuerySelectorAll(selector values.String) core.Value {
	if !el.IsConnected() {
		return values.NewArray(0)
	}

	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorAllArgs(el.id, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		return values.None
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		childEl, err := LoadElement(el.client, el.broker, id)

		if err != nil {
			return values.None
		}

		arr.Push(childEl)
	}

	return arr
}

func (el *HtmlElement) InnerText() values.String {
	val, err := el.innerText.Value()

	if err != nil {
		return values.EmptyString
	}

	return val.(values.String)
}

func (el *HtmlElement) InnerHtml() values.String {
	el.Lock()
	defer el.Unlock()

	return el.innerHtml
}

func (el *HtmlElement) Click() (values.Boolean, error) {
	ctx, cancel := contextWithTimeout()

	defer cancel()

	return events.DispatchEvent(ctx, el.client, el.id, "click")
}

func (el *HtmlElement) Input(value core.Value, timeout values.Int) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	return el.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(el.id, "value", value.String()))
}

func (el *HtmlElement) IsConnected() values.Boolean {
	el.Lock()
	defer el.Unlock()

	return el.connected
}

func (el *HtmlElement) handlePageReload(message interface{}) {
	el.Close()
}

func (el *HtmlElement) handleAttrModified(message interface{}) {
	reply, ok := message.(*dom.AttributeModifiedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != el.id {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	val, err := el.attributes.Value()

	// failed to load
	if err != nil {
		return
	}

	el.Lock()
	defer el.Unlock()

	attrs, ok := val.(*values.Object)

	if !ok {
		return
	}

	attrs.Set(values.NewString(reply.Name), values.NewString(reply.Value))
}

func (el *HtmlElement) handleAttrRemoved(message interface{}) {
	reply, ok := message.(*dom.AttributeRemovedReply)

	// well....
	if !ok {
		return
	}

	// it's not for this element
	if reply.NodeID != el.id {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	val, err := el.attributes.Value()

	// failed to load
	if err != nil {
		return
	}

	el.Lock()
	defer el.Unlock()

	attrs, ok := val.(*values.Object)

	if !ok {
		return
	}

	// TODO: actually, we need to sync it too...
	attrs.Remove(values.NewString(reply.Name))
}

func (el *HtmlElement) handleChildrenCountChanged(message interface{}) {
	reply, ok := message.(*dom.ChildNodeCountUpdatedReply)

	if !ok {
		return
	}

	if reply.NodeID != el.id {
		return
	}

	node, err := el.client.DOM.DescribeNode(context.Background(), dom.NewDescribeNodeArgs())

	if err != nil {
		return
	}

	el.Lock()
	defer el.Unlock()

	el.children = createChildrenArray(node.Node.Children)
}

func (el *HtmlElement) handleChildInserted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeInsertedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id {
		return
	}

	targetIdx := -1
	prevId := reply.PreviousNodeID
	nextId := reply.Node.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id == prevId {
			targetIdx = idx
			break
		}
	}

	if targetIdx == -1 {
		return
	}

	arr := el.children
	el.children = append(arr[:targetIdx], append([]dom.NodeID{nextId}, arr[targetIdx:]...)...)

	if !el.loadedChildren.Ready() {
		return
	}

	loaded, err := el.loadedChildren.Value()

	if err != nil {
		return
	}

	loadedArr := loaded.(*values.Array)

	loadedEl, err := LoadElement(el.client, el.broker, nextId)

	if err != nil {
		return
	}

	loadedArr.Insert(values.NewInt(targetIdx), loadedEl)

	newInnerHtml, err := loadInnerHtml(el.client, el.id)

	if err != nil {
		return
	}

	el.innerHtml = newInnerHtml
}

func (el *HtmlElement) handleChildDeleted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeRemovedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id {
		return
	}

	targetIdx := -1
	targetId := reply.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id == targetId {
			targetIdx = idx
			break
		}
	}

	if targetIdx == -1 {
		return
	}

	arr := el.children
	el.children = append(arr[:targetIdx], arr[targetIdx+1:]...)

	if !el.loadedChildren.Ready() {
		return
	}

	loaded, err := el.loadedChildren.Value()

	if err != nil {
		return
	}

	loadedArr := loaded.(*values.Array)
	loadedArr.RemoveAt(values.NewInt(targetIdx))

	newInnerHtml, err := loadInnerHtml(el.client, el.id)

	if err != nil {
		return
	}

	el.innerHtml = newInnerHtml
}
