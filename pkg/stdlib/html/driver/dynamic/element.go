package dynamic

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
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
	connected      bool
	id             dom.NodeID
	nodeType       values.Int
	nodeName       values.String
	innerHtml      values.String
	innerText      *common.LazyValue
	value          string
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
			SetDepth(-1),
	)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	innerHtml, err := client.DOM.GetOuterHTML(
		ctx,
		dom.NewGetOuterHTMLArgs().SetNodeID(id),
	)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	return NewHtmlElement(
		client,
		broker,
		id,
		node.Node,
		values.NewString(innerHtml.OuterHTML),
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
	el.connected = true
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
	el.value = ""
	el.loadedChildren = common.NewLazyValue(func() (core.Value, error) {
		return loadNodes(client, broker, el.children)
	})

	var childCount int

	if node.ChildNodeCount != nil {
		childCount = *node.ChildNodeCount
	}

	if node.Value != nil {
		el.value = *node.Value
	}

	el.children = make([]dom.NodeID, childCount)

	for idx, child := range node.Children {
		el.children[idx] = child.NodeID
	}

	return el
}

func (el *HtmlElement) Close() error {
	return nil
}

func (el *HtmlElement) Type() core.Type {
	return core.HtmlElementType
}

func (el *HtmlElement) MarshalJSON() ([]byte, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), DefaultTimeout)

	defer cancelFn()

	args := dom.NewGetOuterHTMLArgs()
	args.NodeID = &el.id

	reply, err := el.client.DOM.GetOuterHTML(ctx, args)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(el.id)))
	}

	return json.Marshal(reply.OuterHTML)
}

func (el *HtmlElement) String() string {
	return el.value
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

	out, err := h.Write([]byte(el.value))

	if err != nil {
		return 0
	}

	return out
}

func (el *HtmlElement) Value() core.Value {
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

	return val
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
	return el.innerHtml
}

func (el *HtmlElement) Click() (values.Boolean, error) {
	ctx, cancel := contextWithTimeout()

	defer cancel()

	return events.DispatchEvent(ctx, el.client, el.id, "click")
}
