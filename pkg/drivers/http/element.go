package http

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"strings"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
)

type HTMLElement struct {
	selection *goquery.Selection
	attrs     *values.Object
	styles    *values.Object
	children  *values.Array
}

func NewHTMLElement(node *goquery.Selection) (drivers.HTMLElement, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HTMLElement{node, nil, nil, nil}, nil
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(el.String())
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) String() string {
	ih, err := el.GetInnerHTML(context.Background())

	if err != nil {
		return ""
	}

	return ih.String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		return int64(strings.Compare(el.String(), other.String()))
	default:
		return drivers.Compare(el.Type(), other.Type())
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el.selection
}

func (el *HTMLElement) Hash() uint64 {
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

func (el *HTMLElement) Copy() core.Value {
	c, _ := NewHTMLElement(el.selection.Clone())

	return c
}

func (el *HTMLElement) IsDetached() values.Boolean {
	return values.True
}

func (el *HTMLElement) GetNodeType() values.Int {
	nodes := el.selection.Nodes

	if len(nodes) == 0 {
		return 0
	}

	return values.NewInt(common.FromHTMLType(nodes[0].Type))
}

func (el *HTMLElement) Close() error {
	return nil
}

func (el *HTMLElement) GetNodeName() values.String {
	return values.NewString(goquery.NodeName(el.selection))
}

func (el *HTMLElement) Length() values.Int {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Length()
}

func (el *HTMLElement) GetValue(_ context.Context) (core.Value, error) {
	val, ok := el.selection.Attr("value")

	if ok {
		return values.NewString(val), nil
	}

	return values.EmptyString, nil
}

func (el *HTMLElement) SetValue(_ context.Context, value core.Value) error {
	el.selection.SetAttr("value", value.String())

	return nil
}

func (el *HTMLElement) GetInnerText(_ context.Context) (values.String, error) {
	return values.NewString(el.selection.Text()), nil
}

func (el *HTMLElement) SetInnerText(_ context.Context, innerText values.String) error {
	el.selection.SetText(innerText.String())

	return nil
}

func (el *HTMLElement) GetInnerHTML(_ context.Context) (values.String, error) {
	h, err := el.selection.Html()

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(h), nil
}

func (el *HTMLElement) SetInnerHTML(_ context.Context, value values.String) error {
	el.selection.SetHtml(value.String())

	return nil
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*values.Object, error) {
	if err := el.ensureStyles(ctx); err != nil {
		return values.NewObject(), err
	}

	return el.styles.Copy().(*values.Object), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name values.String) (core.Value, error) {
	if err := el.ensureStyles(ctx); err != nil {
		return values.None, err
	}

	return el.styles.MustGet(name), nil
}

func (el *HTMLElement) SetStyle(ctx context.Context, name values.String, value core.Value) error {
	if err := el.ensureStyles(ctx); err != nil {
		return err
	}

	el.styles.Set(name, value)

	str := common.SerializeStyles(ctx, el.styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) SetStyles(ctx context.Context, newStyles *values.Object) error {
	if newStyles == nil {
		return nil
	}

	if err := el.ensureStyles(ctx); err != nil {
		return err
	}

	newStyles.ForEach(func(i core.Value, key string) bool {
		el.styles.Set(values.NewString(key), i)

		return true
	})

	str := common.SerializeStyles(ctx, el.styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, name ...values.String) error {
	if len(name) == 0 {
		return nil
	}

	if err := el.ensureStyles(ctx); err != nil {
		return err
	}

	for _, s := range name {
		el.styles.Remove(s)
	}

	str := common.SerializeStyles(ctx, el.styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *values.Object) error {
	if attrs == nil {
		return nil
	}

	el.ensureAttrs()

	var err error

	attrs.ForEach(func(value core.Value, key string) bool {
		err = el.SetAttribute(ctx, values.NewString(key), values.NewString(value.String()))

		return err == nil
	})

	return err
}

func (el *HTMLElement) GetAttributes(_ context.Context) (*values.Object, error) {
	el.ensureAttrs()

	return el.attrs.Copy().(*values.Object), nil
}

func (el *HTMLElement) GetAttribute(_ context.Context, name values.String) (core.Value, error) {
	el.ensureAttrs()

	return el.attrs.MustGet(name), nil
}

func (el *HTMLElement) SetAttribute(_ context.Context, name, value values.String) error {
	el.ensureAttrs()

	el.attrs.Set(name, value)
	el.selection.SetAttr(string(name), string(value))

	return nil
}

func (el *HTMLElement) RemoveAttribute(_ context.Context, name ...values.String) error {
	el.ensureAttrs()

	for _, attr := range name {
		el.attrs.Remove(attr)
		el.selection.RemoveAttr(attr.String())
	}

	return nil
}

func (el *HTMLElement) GetChildNodes(_ context.Context) (*values.Array, error) {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Copy().(*values.Array), nil
}

func (el *HTMLElement) GetChildNode(_ context.Context, idx values.Int) (core.Value, error) {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Get(idx), nil
}

func (el *HTMLElement) QuerySelector(_ context.Context, selector values.String) (core.Value, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.None, nil
	}

	res, err := NewHTMLElement(selection)

	if err != nil {
		return values.None, err
	}

	return res, nil
}

func (el *HTMLElement) QuerySelectorAll(_ context.Context, selector values.String) (*values.Array, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.NewArray(0), nil
	}

	arr := values.NewArray(selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		el, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(el)
		}
	})

	return arr, nil
}

func (el *HTMLElement) XPath(_ context.Context, expression values.String) (core.Value, error) {
	h, err := outerHTML(el.selection)

	if err != nil {
		return values.None, err
	}

	exp, err := xpath.Compile(expression.String())

	if err != nil {
		return values.None, err
	}

	rootNode, err := htmlquery.Parse(strings.NewReader(h))

	if err != nil {
		return values.None, err
	}

	out := exp.Evaluate(htmlquery.CreateXPathNavigator(rootNode))

	switch res := out.(type) {
	case *xpath.NodeIterator:
		items := values.NewArray(10)

		for {
			if !res.MoveNext() {
				break
			}

			item, err := parseXPathNode(res.Current().(*htmlquery.NodeNavigator))

			if err != nil {
				return values.None, err
			}

			items.Push(item)
		}

		return items, nil
	default:
		return values.Parse(res), nil
	}
}

func (el *HTMLElement) SetInnerHTMLBySelector(_ context.Context, selector, innerHTML values.String) error {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return drivers.ErrNotFound
	}

	selection.SetHtml(innerHTML.String())

	return nil
}

func (el *HTMLElement) GetInnerHTMLBySelector(_ context.Context, selector values.String) (values.String, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.EmptyString, drivers.ErrNotFound
	}

	str, err := selection.Html()

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(str), nil
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(_ context.Context, selector values.String) (*values.Array, error) {
	var err error
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.EachWithBreak(func(_ int, selection *goquery.Selection) bool {
		str, e := selection.Html()

		if e != nil {
			err = e
			return false
		}

		arr.Push(values.NewString(strings.TrimSpace(str)))

		return true
	})

	if err != nil {
		return values.NewArray(0), err
	}

	return arr, nil
}

func (el *HTMLElement) GetInnerTextBySelector(_ context.Context, selector values.String) (values.String, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.EmptyString, drivers.ErrNotFound
	}

	return values.NewString(selection.Text()), nil
}

func (el *HTMLElement) SetInnerTextBySelector(_ context.Context, selector, innerText values.String) error {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return drivers.ErrNotFound
	}

	selection.SetHtml(innerText.String())

	return nil
}

func (el *HTMLElement) GetInnerTextBySelectorAll(_ context.Context, selector values.String) (*values.Array, error) {
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr, nil
}

func (el *HTMLElement) CountBySelector(_ context.Context, selector values.String) (values.Int, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.ZeroInt, nil
	}

	return values.NewInt(selection.Size()), nil
}

func (el *HTMLElement) ExistsBySelector(_ context.Context, selector values.String) (values.Boolean, error) {
	selection := el.selection.Find(selector.String())

	if selection.Length() == 0 {
		return values.False, nil
	}

	return values.True, nil
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInElement(ctx, el, path)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInElement(ctx, el, path, value)
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) Click(_ context.Context, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClickBySelector(_ context.Context, _ values.String, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClickBySelectorAll(_ context.Context, _ values.String, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Clear(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClearBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Input(_ context.Context, _ core.Value, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) InputBySelector(_ context.Context, _ values.String, _ core.Value, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Select(_ context.Context, _ *values.Array) (*values.Array, error) {
	return nil, core.ErrNotSupported
}

func (el *HTMLElement) SelectBySelector(_ context.Context, _ values.String, _ *values.Array) (*values.Array, error) {
	return nil, core.ErrNotSupported
}

func (el *HTMLElement) ScrollIntoView(_ context.Context, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Focus(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) FocusBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Blur(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) BlurBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Hover(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) HoverBySelector(_ context.Context, _ values.String) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForClass(_ context.Context, _ values.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForAttribute(_ context.Context, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForStyle(_ context.Context, _ values.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ensureStyles(ctx context.Context) error {
	if el.styles == nil {
		styles, err := el.parseStyles(ctx)

		if err != nil {
			return err
		}

		el.styles = styles
	}

	return nil
}

func (el *HTMLElement) parseStyles(ctx context.Context) (*values.Object, error) {
	str, err := el.GetAttribute(ctx, "style")

	if err != nil {
		return values.NewObject(), err
	}

	if str == values.None {
		return values.NewObject(), nil
	}

	styles, err := common.DeserializeStyles(values.NewString(str.String()))

	if err != nil {
		return nil, err
	}

	return styles, nil
}

func (el *HTMLElement) ensureAttrs() {
	if el.attrs == nil {
		el.attrs = el.parseAttrs()
	}
}

func (el *HTMLElement) parseAttrs() *values.Object {
	obj := values.NewObject()

	for _, name := range common.Attributes {
		val, ok := el.selection.Attr(name)

		if ok {
			obj.Set(values.NewString(name), values.NewString(val))
		}
	}

	return obj
}

func (el *HTMLElement) parseChildren() *values.Array {
	children := el.selection.Children()

	arr := values.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		child, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}
