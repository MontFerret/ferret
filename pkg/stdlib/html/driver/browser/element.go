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
	"log"
	"strconv"
	"time"
)

const DefaultTimeout = time.Second * 30

type HtmlElement struct {
	client     *cdp.Client
	id         dom.NodeID
	node       dom.Node
	attributes *values.Object
}

func NewHtmlElement(
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
			SetDepth(-1),
	)

	if err != nil {
		log.Println("ERROR:", err)
		return nil, core.Error(err, strconv.Itoa(int(id)))
	}

	return &HtmlElement{client, id, node.Node, nil}, nil
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
	args.NodeID = &el.node.NodeID

	reply, err := el.client.DOM.GetOuterHTML(ctx, args)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(el.node.NodeID)))
	}

	return json.Marshal(reply.OuterHTML)
}

func (el *HtmlElement) String() string {
	return *el.node.Value
}

func (el *HtmlElement) Compare(other core.Value) int {
	switch other.Type() {
	case core.HtmlDocumentType:
		other := other.(*HtmlElement)

		id := int(el.node.NodeID)
		otherId := int(other.node.NodeID)

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
	return el.node
}

func (el *HtmlElement) Hash() int {
	h := sha512.New()

	out, err := h.Write([]byte(*el.node.Value))

	if err != nil {
		return 0
	}

	return out
}

func (el *HtmlElement) Value() core.Value {
	return values.None
}

func (el *HtmlElement) Length() values.Int {
	if el.node.ChildNodeCount == nil {
		return values.ZeroInt
	}

	return values.NewInt(*el.node.ChildNodeCount)
}

func (el *HtmlElement) NodeType() values.Int {
	return values.NewInt(el.node.NodeType)
}

func (el *HtmlElement) NodeName() values.String {
	return values.NewString(el.node.NodeName)
}

func (el *HtmlElement) GetAttributes() core.Value {
	if el.attributes == nil {
		el.attributes = el.parseAttrs()
	}

	return el.attributes
}

func (el *HtmlElement) GetAttribute(name values.String) core.Value {
	if el.attributes == nil {
		el.attributes = el.parseAttrs()
	}

	val, found := el.attributes.Get(name)

	if !found {
		return values.None
	}

	return val
}

func (el *HtmlElement) GetChildNodes() core.Value {
	arr := values.NewArray(len(el.node.Children))

	for idx := range el.node.Children {
		el := el.GetChildNode(values.NewInt(idx))

		if el != values.None {
			arr.Push(el)
		}
	}

	return arr
}

func (el *HtmlElement) GetChildNode(idx values.Int) core.Value {
	if el.Length() < idx {
		return values.None
	}

	childNode := el.node.Children[idx]

	return &HtmlElement{el.client, childNode.NodeID, childNode, nil}
}

func (el *HtmlElement) QuerySelector(selector values.String) core.Value {
	ctx := context.Background()

	selectorArgs := dom.NewQuerySelectorArgs(el.id, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		el.logErr(err, selector.String())
		return values.None
	}

	res, err := NewHtmlElement(el.client, found.NodeID)

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
		el.logErr(err, selector.String())
		return values.None
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		childEl, err := NewHtmlElement(el.client, id)

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
	ctx, cancelFn := el.createCtx()

	defer cancelFn()

	res, err := el.client.DOM.GetOuterHTML(ctx, dom.NewGetOuterHTMLArgs().SetNodeID(el.id))

	if err != nil {
		el.logErr(err)

		return values.EmptyString
	}

	return values.NewString(res.OuterHTML)
}

func (el *HtmlElement) createCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
}

func (el *HtmlElement) parseAttrs() *values.Object {
	var attr values.String

	res := values.NewObject()

	for _, el := range el.node.Attributes {
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

func (el *HtmlElement) logErr(values ...interface{}) {
	args := make([]interface{}, 0, len(values)+1)
	args = append(args, "ERROR:")
	args = append(args, values...)
	args = append(args, "id:", el.node.NodeID)

	log.Println(args...)
}
