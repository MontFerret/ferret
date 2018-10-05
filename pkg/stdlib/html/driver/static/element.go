package static

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
	"github.com/PuerkitoBio/goquery"
	"hash/fnv"
)

type HtmlElement struct {
	selection *goquery.Selection
	attrs     *values.Object
	children  *values.Array
}

func NewHtmlElement(node *goquery.Selection) (*HtmlElement, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HtmlElement{node, nil, nil}, nil
}

func (el *HtmlElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(el.InnerText().String())
}

func (el *HtmlElement) Type() core.Type {
	return core.HtmlElementType
}

func (el *HtmlElement) String() string {
	return el.InnerHtml().String()
}

func (el *HtmlElement) Compare(other core.Value) int {
	switch other.Type() {
	case core.HtmlElementType:
		// TODO: complete the comparison
		return -1
	default:
		if other.Type() > core.HtmlElementType {
			return -1
		}

		return 1
	}
}

func (el *HtmlElement) Unwrap() interface{} {
	return el.selection
}

func (el *HtmlElement) Hash() uint64 {
	str, err := el.selection.Html()

	if err != nil {
		return 0
	}

	h := fnv.New64a()

	h.Write([]byte(el.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(str))

	return h.Sum64()
}

func (el *HtmlElement) Clone() core.Value {
	c, _ := NewHtmlElement(el.selection.Clone())

	return c
}

func (el *HtmlElement) NodeType() values.Int {
	nodes := el.selection.Nodes

	if len(nodes) == 0 {
		return 0
	}

	return values.NewInt(common.ToHtmlType(nodes[0].Type))
}

func (el *HtmlElement) NodeName() values.String {
	return values.NewString(goquery.NodeName(el.selection))
}

func (el *HtmlElement) Length() values.Int {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Length()
}

func (el *HtmlElement) Value() core.Value {
	val, ok := el.selection.Attr("value")

	if ok {
		return values.NewString(val)
	}

	return values.EmptyString
}

func (el *HtmlElement) InnerText() values.String {
	return values.NewString(el.selection.Text())
}

func (el *HtmlElement) InnerHtml() values.String {
	h, err := el.selection.Html()

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(h)
}

func (el *HtmlElement) GetAttributes() core.Value {
	if el.attrs == nil {
		el.attrs = el.parseAttrs()
	}

	return el.attrs
}

func (el *HtmlElement) GetAttribute(name values.String) core.Value {
	v, ok := el.selection.Attr(name.String())

	if ok {
		return values.NewString(v)
	}

	return values.None
}

func (el *HtmlElement) GetChildNodes() core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children
}

func (el *HtmlElement) GetChildNode(idx values.Int) core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Get(idx)
}

func (el *HtmlElement) QuerySelector(selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	res, err := NewHtmlElement(selection)

	if err != nil {
		return values.None
	}

	return res
}

func (el *HtmlElement) QuerySelectorAll(selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	arr := values.NewArray(selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		el, err := NewHtmlElement(selection)

		if err == nil {
			arr.Push(el)
		}
	})

	return arr
}

func (el *HtmlElement) parseAttrs() *values.Object {
	obj := values.NewObject()

	for _, name := range common.Attributes {
		val, ok := el.selection.Attr(name)

		if ok {
			obj.Set(values.NewString(name), values.NewString(val))
		}
	}

	return obj
}

func (el *HtmlElement) parseChildren() *values.Array {
	children := el.selection.Children()

	arr := values.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		//name := goquery.NodeName(selection)
		//if name != "#text" && name != "#comment" {
		//	child, err := NewHtmlElement(selection)
		//
		//	if err == nil {
		//		arr.Push(child)
		//	}
		//}

		child, err := NewHtmlElement(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}
