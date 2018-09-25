package browser

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"strconv"
	"time"
)

const DefaultTimeout = time.Second * 30

type HtmlElement struct {
	client         *cdp.Client
	id             dom.NodeID
	nodeType       values.Int
	nodeName       values.String
	value          string
	attributes     *values.Object
	children       []dom.NodeID
	loadedChildren *values.Array
}

func LoadElement(
	client *cdp.Client,
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

	return NewHtmlElement(client, id, node.Node), nil
}

func NewHtmlElement(
	client *cdp.Client,
	id dom.NodeID,
	node dom.Node,
) *HtmlElement {
	el := new(HtmlElement)
	el.client = client
	el.id = id
	el.nodeType = values.NewInt(node.NodeType)
	el.nodeName = values.NewString(node.NodeName)
	el.value = ""
	el.attributes = parseAttrs(node.Attributes)

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
	// el.client = nil

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
	return el.attributes
}

func (el *HtmlElement) GetAttribute(name values.String) core.Value {
	val, found := el.attributes.Get(name)

	if !found {
		return values.None
	}

	return val
}

func (el *HtmlElement) GetChildNodes() core.Value {
	if el.loadedChildren == nil {
		el.loadedChildren = loadNodes(el.client, el.children)
	}

	return el.loadedChildren
}

func (el *HtmlElement) GetChildNode(idx values.Int) core.Value {
	if el.loadedChildren == nil {
		el.loadedChildren = loadNodes(el.client, el.children)
	}

	return el.loadedChildren.Get(idx)
}

func (el *HtmlElement) QuerySelector(selector values.String) core.Value {
	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorArgs(el.id, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		return values.None
	}

	res, err := LoadElement(el.client, found.NodeID)

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
		childEl, err := LoadElement(el.client, id)

		if err != nil {
			return values.None
		}

		arr.Push(childEl)
	}

	return arr
}

func (el *HtmlElement) InnerText() values.String {
	h := el.InnerHtml()

	if h == values.EmptyString {
		return h
	}

	buff := bytes.NewBuffer([]byte(h))

	parsed, err := goquery.NewDocumentFromReader(buff)

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(parsed.Text())
}

func (el *HtmlElement) InnerHtml() values.String {
	ctx, cancelFn := createCtx()

	defer cancelFn()

	res, err := el.client.DOM.GetOuterHTML(ctx, dom.NewGetOuterHTMLArgs().SetNodeID(el.id))

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(res.OuterHTML)
}

func (el *HtmlElement) Click() (values.Boolean, error) {
	ctx, cancel := createCtx()

	defer cancel()

	return DispatchEvent(ctx, el.client, el.id, "click")
}

func createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func parseAttrs(attrs []string) *values.Object {
	var attr values.String

	res := values.NewObject()

	for _, el := range attrs {
		str := values.NewString(el)

		if common.IsAttribute(el) {
			attr = str
			res.Set(str, values.EmptyString)
		} else {
			current, ok := res.Get(attr)

			if ok {
				res.Set(attr, current.(values.String).Concat(values.SpaceString).Concat(str))
			}
		}
	}

	return res
}

func loadNodes(client *cdp.Client, nodes []dom.NodeID) *values.Array {
	arr := values.NewArray(len(nodes))

	for _, id := range nodes {
		child, err := LoadElement(client, id)

		if err != nil {
			break
		}

		arr.Push(child)
	}

	return arr
}
