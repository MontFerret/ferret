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
	return json.Marshal(el.GetInnerText(context.Background()).String())
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) String() string {
	return el.GetInnerHTML(context.Background()).String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		ctx, fn := drivers.WithDefaultTimeout(context.Background())
		defer fn()

		return el.GetInnerHTML(ctx).Compare(other.GetInnerHTML(ctx))
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

func (el *HTMLElement) GetValue(_ context.Context) core.Value {
	val, ok := el.selection.Attr("value")

	if ok {
		return values.NewString(val)
	}

	return values.EmptyString
}

func (el *HTMLElement) SetValue(_ context.Context, value core.Value) error {
	el.selection.SetAttr("value", value.String())

	return nil
}

func (el *HTMLElement) GetInnerText(_ context.Context) values.String {
	return values.NewString(el.selection.Text())
}

func (el *HTMLElement) GetInnerHTML(_ context.Context) values.String {
	h, err := el.selection.Html()

	if err != nil {
		return values.EmptyString
	}

	return values.NewString(h)
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

func (el *HTMLElement) GetAttributes(_ context.Context) *values.Object {
	el.ensureAttrs()

	return el.attrs.Copy().(*values.Object)
}

func (el *HTMLElement) GetAttribute(_ context.Context, name values.String) core.Value {
	el.ensureAttrs()

	return el.attrs.MustGet(name)
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

func (el *HTMLElement) GetChildNodes(_ context.Context) core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children
}

func (el *HTMLElement) GetChildNode(_ context.Context, idx values.Int) core.Value {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Get(idx)
}

func (el *HTMLElement) QuerySelector(_ context.Context, selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.None
	}

	res, err := NewHTMLElement(selection)

	if err != nil {
		return values.None
	}

	return res
}

func (el *HTMLElement) QuerySelectorAll(_ context.Context, selector values.String) core.Value {
	selection := el.selection.Find(selector.String())

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

func (el *HTMLElement) XPath(_ context.Context, _ values.String) (core.Value, error) {
	return values.None, core.ErrNotSupported
}

func (el *HTMLElement) InnerHTMLBySelector(_ context.Context, selector values.String) values.String {
	selection := el.selection.Find(selector.String())

	str, err := selection.Html()

	// TODO: log error
	if err != nil {
		return values.EmptyString
	}

	return values.NewString(str)
}

func (el *HTMLElement) InnerHTMLBySelectorAll(_ context.Context, selector values.String) *values.Array {
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		str, err := selection.Html()

		// TODO: log error
		if err == nil {
			arr.Push(values.NewString(strings.TrimSpace(str)))
		}
	})

	return arr
}

func (el *HTMLElement) InnerTextBySelector(_ context.Context, selector values.String) values.String {
	selection := el.selection.Find(selector.String())

	return values.NewString(selection.Text())
}

func (el *HTMLElement) InnerTextBySelectorAll(_ context.Context, selector values.String) *values.Array {
	selection := el.selection.Find(selector.String())
	arr := values.NewArray(selection.Length())

	selection.Each(func(_ int, selection *goquery.Selection) {
		arr.Push(values.NewString(selection.Text()))
	})

	return arr
}

func (el *HTMLElement) CountBySelector(_ context.Context, selector values.String) values.Int {
	selection := el.selection.Find(selector.String())

	if selection == nil {
		return values.ZeroInt
	}

	return values.NewInt(selection.Size())
}

func (el *HTMLElement) ExistsBySelector(_ context.Context, selector values.String) values.Boolean {
	selection := el.selection.Closest(selector.String())

	if selection == nil {
		return values.False
	}

	return values.True
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

func (el *HTMLElement) Click(_ context.Context) (values.Boolean, error) {
	return false, core.ErrNotSupported
}

func (el *HTMLElement) Input(_ context.Context, _ core.Value, _ values.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Select(_ context.Context, _ *values.Array) (*values.Array, error) {
	return nil, core.ErrNotSupported
}

func (el *HTMLElement) ScrollIntoView(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Hover(_ context.Context) error {
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
	str := el.GetAttribute(ctx, "style")

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
