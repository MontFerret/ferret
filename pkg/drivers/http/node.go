package http

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HTMLNode struct {
	selection *goquery.Selection
	attrs     *values.Object
	children  *values.Array
}

func NewHTMLNode(node *goquery.Selection) (*HTMLNode, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HTMLNode{node, nil, nil}, nil
}

func (nd *HTMLNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(nd.InnerText().String())
}

func (nd *HTMLNode) Type() core.Type {
	return drivers.HTMLNodeType
}

func (nd *HTMLNode) String() string {
	return nd.InnerHTML().String()
}

func (nd *HTMLNode) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLNodeType:
		other := other.(drivers.HTMLNode)

		return nd.InnerHTML().Compare(other.InnerHTML())
	default:
		return drivers.Compare(nd.Type(), other.Type())
	}
}

func (nd *HTMLNode) Unwrap() interface{} {
	return nd.selection
}

func (nd *HTMLNode) Hash() uint64 {
	str, err := nd.selection.Html()

	if err != nil {
		return 0
	}

	h := fnv.New64a()

	h.Write([]byte(nd.Type().Name()))
	h.Write([]byte(":"))
	h.Write([]byte(str))

	return h.Sum64()
}

func (nd *HTMLNode) Copy() core.Value {
	c, _ := NewHTMLNode(nd.selection.Clone())

	return c
}

func (nd *HTMLNode) NodeType() values.Int {
	nodes := nd.selection.Nodes

	if len(nodes) == 0 {
		return 0
	}

	return values.NewInt(common.ToHTMLType(nodes[0].Type))
}

func (nd *HTMLNode) NodeName() values.String {
	return values.NewString(goquery.NodeName(nd.selection))
}

func (nd *HTMLNode) Length() values.Int {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children.Length()
}

func (nd *HTMLNode) GetValue() core.Value {
	val, ok := nd.selection.Attr("value")

	if ok {
		return values.NewString(val)
	}

	return values.EmptyString
}

func (nd *HTMLNode) SetValue(value core.Value) error {
	nd.selection.SetAttr("value", value.String())

	return nil
}

func (nd *HTMLNode) InnerText() values.String {
	return values.NewString(nd.selection.Text())
}

func (nd *HTMLNode) InnerHTML() values.String {
	h, err := nd.selection.Html()

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(h)
}

func (nd *HTMLNode) GetAttributes() core.Value {
	if nd.attrs == nil {
		nd.attrs = nd.parseAttrs()
	}

	return nd.attrs
}

func (nd *HTMLNode) GetAttribute(name values.String) core.Value {
	v, ok := nd.selection.Attr(name.String())

	if ok {
		return values.NewString(v)
	}

	return values.None
}

func (nd *HTMLNode) SetAttribute(name, value values.String) error {
	nd.selection.SetAttr(string(name), string(value))

	return nil
}

func (nd *HTMLNode) GetChildNodes() core.Value {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children
}

func (nd *HTMLNode) GetChildNode(idx values.Int) core.Value {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children.Get(idx)
}

func (nd *HTMLNode) QuerySelector(selector values.String) core.Value {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	res, err := NewHTMLNode(selection)

	if err != nil {
		return values.None
	}

	return res
}

func (nd *HTMLNode) QuerySelectorAll(selector values.String) core.Value {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	arr := values.NewArray(selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		el, err := NewHTMLNode(selection)

		if err == nil {
			arr.Push(el)
		}
	})

	return arr
}

func (nd *HTMLNode) InnerHTMLBySelector(selector values.String) values.String {
	selection := nd.selection.Find(selector.String())

	str, err := selection.Html()

	// TODO: log error
	if err != nil {
		return values.EmptyString
	}

	return values.NewString(str)
}

func (nd *HTMLNode) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	selection := nd.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		str, err := selection.Html()

		// TODO: log error
		if err == nil {
			arr.Push(values.NewString(str))
		}
	})

	return arr
}

func (nd *HTMLNode) InnerTextBySelector(selector values.String) values.String {
	selection := nd.selection.Find(selector.String())

	return values.NewString(selection.Text())
}

func (nd *HTMLNode) InnerTextBySelectorAll(selector values.String) *values.Array {
	selection := nd.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr
}

func (nd *HTMLNode) CountBySelector(selector values.String) values.Int {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.ZeroInt
	}

	return values.NewInt(selection.Size())
}

func (nd *HTMLNode) ExistsBySelector(selector values.String) values.Boolean {
	selection := nd.selection.Closest(selector.String())

	if selection == nil {
		return values.False
	}

	return values.True
}

func (nd *HTMLNode) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetIn(ctx, nd, path)
}

func (nd *HTMLNode) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetIn(ctx, nd, path, value)
}

func (nd *HTMLNode) Iterate(_ context.Context) (collections.CollectionIterator, error) {
	return common.NewIterator(nd)
}

func (nd *HTMLNode) parseAttrs() *values.Object {
	obj := values.NewObject()

	for _, name := range common.Attributes {
		val, ok := nd.selection.Attr(name)

		if ok {
			obj.Set(values.NewString(name), values.NewString(val))
		}
	}

	return obj
}

func (nd *HTMLNode) parseChildren() *values.Array {
	children := nd.selection.Children()

	arr := values.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		child, err := NewHTMLNode(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}
