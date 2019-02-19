package http

import (
	"context"
	"encoding/json"
	"hash/fnv"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
)

type HTMLElement struct {
	selection *goquery.Selection
	attrs     *values.Object
	children  *values.Array
}

func NewHTMLElement(node *goquery.Selection) (drivers.HTMLElement, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HTMLElement{node, nil, nil}, nil
}

func (nd *HTMLElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(nd.InnerText().String())
}

func (nd *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (nd *HTMLElement) String() string {
	return nd.InnerHTML().String()
}

func (nd *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		return nd.InnerHTML().Compare(other.InnerHTML())
	default:
		return drivers.Compare(nd.Type(), other.Type())
	}
}

func (nd *HTMLElement) Unwrap() interface{} {
	return nd.selection
}

func (nd *HTMLElement) Hash() uint64 {
	str, err := nd.selection.Html()

	if err != nil {
		return 0
	}

	h := fnv.New64a()

	h.Write([]byte(nd.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(str))

	return h.Sum64()
}

func (nd *HTMLElement) Copy() core.Value {
	c, _ := NewHTMLElement(nd.selection.Clone())

	return c
}

func (nd *HTMLElement) NodeType() values.Int {
	nodes := nd.selection.Nodes

	if len(nodes) == 0 {
		return 0
	}

	return values.NewInt(common.ToHTMLType(nodes[0].Type))
}

func (nd *HTMLElement) Close() error {
	return nil
}

func (nd *HTMLElement) NodeName() values.String {
	return values.NewString(goquery.NodeName(nd.selection))
}

func (nd *HTMLElement) Length() values.Int {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children.Length()
}

func (nd *HTMLElement) GetValue() core.Value {
	val, ok := nd.selection.Attr("value")

	if ok {
		return values.NewString(val)
	}

	return values.EmptyString
}

func (nd *HTMLElement) SetValue(value core.Value) error {
	nd.selection.SetAttr("value", value.String())

	return nil
}

func (nd *HTMLElement) InnerText() values.String {
	return values.NewString(nd.selection.Text())
}

func (nd *HTMLElement) InnerHTML() values.String {
	h, err := nd.selection.Html()

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(h)
}

func (nd *HTMLElement) GetAttributes() core.Value {
	if nd.attrs == nil {
		nd.attrs = nd.parseAttrs()
	}

	return nd.attrs
}

func (nd *HTMLElement) GetAttribute(name values.String) core.Value {
	v, ok := nd.selection.Attr(name.String())

	if ok {
		return values.NewString(v)
	}

	return values.None
}

func (nd *HTMLElement) SetAttribute(name, value values.String) error {
	nd.selection.SetAttr(string(name), string(value))

	return nil
}

func (nd *HTMLElement) GetChildNodes() core.Value {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children
}

func (nd *HTMLElement) GetChildNode(idx values.Int) core.Value {
	if nd.children == nil {
		nd.children = nd.parseChildren()
	}

	return nd.children.Get(idx)
}

func (nd *HTMLElement) QuerySelector(selector values.String) core.Value {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	res, err := NewHTMLElement(selection)

	if err != nil {
		return values.None
	}

	return res
}

func (nd *HTMLElement) QuerySelectorAll(selector values.String) core.Value {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	arr := values.NewArray(selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		el, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(el)
		}
	})

	return arr
}

func (nd *HTMLElement) InnerHTMLBySelector(selector values.String) values.String {
	selection := nd.selection.Find(selector.String())

	str, err := selection.Html()

	// TODO: log error
	if err != nil {
		return values.EmptyString
	}

	return values.NewString(str)
}

func (nd *HTMLElement) InnerHTMLBySelectorAll(selector values.String) *values.Array {
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

func (nd *HTMLElement) InnerTextBySelector(selector values.String) values.String {
	selection := nd.selection.Find(selector.String())

	return values.NewString(selection.Text())
}

func (nd *HTMLElement) InnerTextBySelectorAll(selector values.String) *values.Array {
	selection := nd.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr
}

func (nd *HTMLElement) CountBySelector(selector values.String) values.Int {
	selection := nd.selection.Find(selector.String())

	if selection == nil {
		return values.ZeroInt
	}

	return values.NewInt(selection.Size())
}

func (nd *HTMLElement) ExistsBySelector(selector values.String) values.Boolean {
	selection := nd.selection.Closest(selector.String())

	if selection == nil {
		return values.False
	}

	return values.True
}

func (nd *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInElement(ctx, nd, path)
}

func (nd *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInElement(ctx, nd, path, value)
}

func (nd *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(nd)
}

func (nd *HTMLElement) Click() (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (nd *HTMLElement) Input(_ core.Value, _ values.Int) error {
	return core.ErrNotSupported
}

func (nd *HTMLElement) Select(_ *values.Array) (*values.Array, error) {
	return nil, core.ErrNotSupported
}

func (nd *HTMLElement) ScrollIntoView() error {
	return core.ErrNotSupported
}

func (nd *HTMLElement) Hover() error {
	return core.ErrNotSupported
}

func (nd *HTMLElement) WaitForClass(_ values.String, _ values.Int) error {
	return core.ErrNotSupported
}

func (nd *HTMLElement) parseAttrs() *values.Object {
	obj := values.NewObject()

	for _, name := range common.Attributes {
		val, ok := nd.selection.Attr(name)

		if ok {
			obj.Set(values.NewString(name), values.NewString(val))
		}
	}

	return obj
}

func (nd *HTMLElement) parseChildren() *values.Array {
	children := nd.selection.Children()

	arr := values.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		child, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}
