package dynamic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/events"
	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/rs/zerolog"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"
)

const DefaultTimeout = time.Second * 30

type HTMLElement struct {
	sync.Mutex
	logger         *zerolog.Logger
	client         *cdp.Client
	broker         *events.EventBroker
	connected      values.Boolean
	id             dom.NodeID
	nodeType       values.Int
	nodeName       values.String
	innerHTML      values.String
	innerText      *common.LazyValue
	value          core.Value
	rawAttrs       []string
	attributes     *common.LazyValue
	children       []dom.NodeID
	loadedChildren *common.LazyValue
}

func LoadElement(
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	id dom.NodeID,
) (*HTMLElement, error) {
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

	innerHTML, err := loadInnerHTML(client, id)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	return NewHTMLElement(
		logger,
		client,
		broker,
		id,
		node.Node,
		innerHTML,
	), nil
}

func NewHTMLElement(
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	id dom.NodeID,
	node dom.Node,
	innerHTML values.String,
) *HTMLElement {
	el := new(HTMLElement)
	el.logger = logger
	el.client = client
	el.broker = broker
	el.connected = values.True
	el.id = id
	el.nodeType = values.NewInt(node.NodeType)
	el.nodeName = values.NewString(node.NodeName)
	el.innerHTML = innerHTML
	el.innerText = common.NewLazyValue(el.loadInnerText)
	el.rawAttrs = node.Attributes[:]
	el.attributes = common.NewLazyValue(el.loadAttrs)
	el.value = values.EmptyString
	el.loadedChildren = common.NewLazyValue(el.loadChildren)

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

func (el *HTMLElement) Close() error {
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

		id := int(el.id)
		otherID := int(other.id)

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

	val, err := eval.Property(ctx, el.client, el.id, "value")

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Msg("failed to get node value")

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

	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorArgs(el.id, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return values.None
	}

	res, err := LoadElement(el.logger, el.client, el.broker, found.NodeID)

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
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

	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorAllArgs(el.id, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Str("selector", selector.String()).
			Msg("failed to retrieve nodes by selector")

		return values.None
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		childEl, err := LoadElement(el.logger, el.client, el.broker, id)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
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

func (el *HTMLElement) InnerHTML() values.String {
	el.Lock()
	defer el.Unlock()

	return el.innerHTML
}

func (el *HTMLElement) Click() (values.Boolean, error) {
	ctx, cancel := contextWithTimeout()

	defer cancel()

	return events.DispatchEvent(ctx, el.client, el.id, "click")
}

func (el *HTMLElement) Input(value core.Value) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	return el.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(el.id, "value", value.String()))
}

func (el *HTMLElement) IsConnected() values.Boolean {
	el.Lock()
	defer el.Unlock()

	return el.connected
}

func (el *HTMLElement) loadInnerText() (core.Value, error) {
	h := el.InnerHTML()

	if h == values.EmptyString {
		return h, nil
	}

	buff := bytes.NewBuffer([]byte(h))

	parsed, err := goquery.NewDocumentFromReader(buff)

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Msg("failed to parse inner html")

		return values.EmptyString, err
	}

	return values.NewString(parsed.Text()), nil
}

func (el *HTMLElement) loadAttrs() (core.Value, error) {
	return parseAttrs(el.rawAttrs), nil
}

func (el *HTMLElement) loadChildren() (core.Value, error) {
	if !el.IsConnected() {
		return values.NewArray(0), nil
	}

	loaded, err := loadNodes(el.logger, el.client, el.broker, el.children)

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Msg("failed to load child nodes")

		return values.None, err
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
	if reply.NodeID != el.id {
		return
	}

	el.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
				Msg("failed to update node")

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
	if reply.NodeID != el.id {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	el.attributes.Write(func(v core.Value, err error) {
		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
				Msg("failed to update node")

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

	if reply.NodeID != el.id {
		return
	}

	node, err := el.client.DOM.DescribeNode(context.Background(), dom.NewDescribeNodeArgs())

	if err != nil {
		el.logger.Error().
			Timestamp().
			Err(err).
			Int("id", int(el.id)).
			Msg("failed to update node")

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

	if reply.ParentNodeID != el.id {
		return
	}

	targetIDx := -1
	prevID := reply.PreviousNodeID
	nextID := reply.Node.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id == prevID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	arr := el.children
	el.children = append(arr[:targetIDx], append([]dom.NodeID{nextID}, arr[targetIDx:]...)...)

	if !el.loadedChildren.Ready() {
		return
	}

	el.loadedChildren.Write(func(v core.Value, err error) {
		loadedArr := v.(*values.Array)
		loadedEl, err := LoadElement(el.logger, el.client, el.broker, nextID)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
				Msg("failed to load an inserted node")

			return
		}

		loadedArr.Insert(values.NewInt(targetIDx), loadedEl)

		newInnerHTML, err := loadInnerHTML(el.client, el.id)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
				Msg("failed to update node")

			return
		}

		el.innerHTML = newInnerHTML
	})
}

func (el *HTMLElement) handleChildDeleted(message interface{}) {
	reply, ok := message.(*dom.ChildNodeRemovedReply)

	if !ok {
		return
	}

	if reply.ParentNodeID != el.id {
		return
	}

	targetIDx := -1
	targetID := reply.NodeID

	el.Lock()
	defer el.Unlock()

	for idx, id := range el.children {
		if id == targetID {
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
				Int("id", int(el.id)).
				Msg("failed to update node")

			return
		}

		loadedArr := v.(*values.Array)
		loadedArr.RemoveAt(values.NewInt(targetIDx))

		newInnerHTML, err := loadInnerHTML(el.client, el.id)

		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("id", int(el.id)).
				Msg("failed to update node")

			return
		}

		el.innerHTML = newInnerHTML
	})
}
