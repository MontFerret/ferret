package http

import (
	"context"
	"hash/fnv"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type HTMLElement struct {
	selection *goquery.Selection
	attrs     *internal.Object
	styles    *internal.Object
	children  *internal.Array
}

func NewHTMLElement(node *goquery.Selection) (drivers.HTMLElement, error) {
	if node == nil {
		return nil, core.Error(core.ErrMissedArgument, "element selection")
	}

	return &HTMLElement{node, nil, nil, nil}, nil
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(el.String(), jettison.NoHTMLEscaping())
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
		return drivers.CompareTypes(el.Type(), other.Type())
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

func (el *HTMLElement) GetNodeType(_ context.Context) (core.Int, error) {
	nodes := el.selection.Nodes

	if len(nodes) == 0 {
		return 0, nil
	}

	return core.NewInt(common.FromHTMLType(nodes[0].Type)), nil
}

func (el *HTMLElement) Close() error {
	return nil
}

func (el *HTMLElement) GetNodeName(_ context.Context) (core.String, error) {
	return core.NewString(goquery.NodeName(el.selection)), nil
}

func (el *HTMLElement) Length() core.Int {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Length()
}

func (el *HTMLElement) GetValue(_ context.Context) (core.Value, error) {
	val, ok := el.selection.Attr("value")

	if ok {
		return core.NewString(val), nil
	}

	return core.EmptyString, nil
}

func (el *HTMLElement) SetValue(_ context.Context, value core.Value) error {
	el.selection.SetAttr("value", value.String())

	return nil
}

func (el *HTMLElement) GetInnerText(_ context.Context) (core.String, error) {
	return core.NewString(el.selection.Text()), nil
}

func (el *HTMLElement) SetInnerText(_ context.Context, innerText core.String) error {
	el.selection.SetText(innerText.String())

	return nil
}

func (el *HTMLElement) GetInnerHTML(_ context.Context) (core.String, error) {
	h, err := el.selection.Html()

	if err != nil {
		return core.EmptyString, err
	}

	return core.NewString(h), nil
}

func (el *HTMLElement) SetInnerHTML(_ context.Context, value core.String) error {
	el.selection.SetHtml(value.String())

	return nil
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*internal.Object, error) {
	if err := el.ensureStyles(ctx); err != nil {
		return internal.NewObject(), err
	}

	return el.styles.Copy().(*internal.Object), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name core.String) (core.Value, error) {
	if err := el.ensureStyles(ctx); err != nil {
		return core.None, err
	}

	return el.styles.MustGet(name), nil
}

func (el *HTMLElement) SetStyle(ctx context.Context, name, value core.String) error {
	if err := el.ensureStyles(ctx); err != nil {
		return err
	}

	el.styles.Set(name, value)

	str := common.SerializeStyles(ctx, el.styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) SetStyles(ctx context.Context, newStyles *internal.Object) error {
	if newStyles == nil {
		return nil
	}

	if err := el.ensureStyles(ctx); err != nil {
		return err
	}

	newStyles.ForEach(func(i core.Value, key string) bool {
		el.styles.Set(core.NewString(key), i)

		return true
	})

	str := common.SerializeStyles(ctx, el.styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, name ...core.String) error {
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

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *internal.Object) error {
	if attrs == nil {
		return nil
	}

	el.ensureAttrs()

	var err error

	attrs.ForEach(func(value core.Value, key string) bool {
		err = el.SetAttribute(ctx, core.NewString(key), core.NewString(value.String()))

		return err == nil
	})

	return err
}

func (el *HTMLElement) GetAttributes(_ context.Context) (*internal.Object, error) {
	el.ensureAttrs()

	return el.attrs.Copy().(*internal.Object), nil
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name core.String) (core.Value, error) {
	el.ensureAttrs()

	if name == common.AttrNameStyle {
		return el.GetStyles(ctx)
	}

	return el.attrs.MustGet(name), nil
}

func (el *HTMLElement) SetAttribute(_ context.Context, name, value core.String) error {
	el.ensureAttrs()

	if name == common.AttrNameStyle {
		el.styles = nil
	}

	el.attrs.Set(name, value)
	el.selection.SetAttr(string(name), string(value))

	return nil
}

func (el *HTMLElement) RemoveAttribute(_ context.Context, name ...core.String) error {
	el.ensureAttrs()

	for _, attr := range name {
		el.attrs.Remove(attr)
		el.selection.RemoveAttr(attr.String())
	}

	return nil
}

func (el *HTMLElement) GetChildNodes(_ context.Context) (*internal.Array, error) {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Copy().(*internal.Array), nil
}

func (el *HTMLElement) GetChildNode(_ context.Context, idx core.Int) (core.Value, error) {
	if el.children == nil {
		el.children = el.parseChildren()
	}

	return el.children.Get(idx), nil
}

func (el *HTMLElement) QuerySelector(_ context.Context, selector drivers.QuerySelector) (core.Value, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return core.None, drivers.ErrNotFound
		}

		res, err := NewHTMLElement(selection)

		if err != nil {
			return core.None, err
		}

		return res, nil
	}

	found, err := EvalXPathToNode(el.selection, selector.String())

	if err != nil {
		return core.None, err
	}

	if found == nil {
		return core.None, drivers.ErrNotFound
	}

	return found, nil
}

func (el *HTMLElement) QuerySelectorAll(_ context.Context, selector drivers.QuerySelector) (*internal.Array, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return internal.NewArray(0), nil
		}

		arr := internal.NewArray(selection.Length())

		selection.Each(func(i int, selection *goquery.Selection) {
			el, err := NewHTMLElement(selection)

			if err == nil {
				arr.Push(el)
			}
		})

		return arr, nil
	}

	return EvalXPathToNodes(el.selection, selector.String())
}

func (el *HTMLElement) XPath(_ context.Context, expression core.String) (core.Value, error) {
	return EvalXPathTo(el.selection, expression.String())
}

func (el *HTMLElement) SetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector, innerHTML core.String) error {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return drivers.ErrNotFound
		}

		selection.SetHtml(innerHTML.String())
	}

	found, err := EvalXPathToElement(el.selection, selector.String())

	if err != nil {
		return err
	}

	if found == nil {
		return drivers.ErrNotFound
	}

	return found.SetInnerHTML(ctx, innerHTML)
}

func (el *HTMLElement) GetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector) (core.String, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return core.EmptyString, drivers.ErrNotFound
		}

		str, err := selection.Html()

		if err != nil {
			return core.EmptyString, err
		}

		return core.NewString(str), nil
	}

	found, err := EvalXPathToElement(el.selection, selector.String())

	if err != nil {
		return core.EmptyString, err
	}

	if found == nil {
		return core.EmptyString, drivers.ErrNotFound
	}

	return found.GetInnerHTML(ctx)
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*internal.Array, error) {
	if selector.Kind() == drivers.CSSSelector {
		var err error
		selection := el.selection.Find(selector.String())
		arr := internal.NewArray(selection.Length())

		selection.EachWithBreak(func(_ int, selection *goquery.Selection) bool {
			str, e := selection.Html()

			if e != nil {
				err = e
				return false
			}

			arr.Push(core.NewString(strings.TrimSpace(str)))

			return true
		})

		if err != nil {
			return internal.NewArray(0), err
		}

		return arr, nil
	}

	return EvalXPathToNodesWith(el.selection, selector.String(), func(node *html.Node) (core.Value, error) {
		n, err := parseXPathNode(node)

		if err != nil {
			return core.None, err
		}

		found, err := drivers.ToElement(n)

		if err != nil {
			return core.None, err
		}

		return found.GetInnerHTML(ctx)
	})
}

func (el *HTMLElement) GetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector) (core.String, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return core.EmptyString, drivers.ErrNotFound
		}

		return core.NewString(selection.Text()), nil
	}

	found, err := EvalXPathToElement(el.selection, selector.String())

	if err != nil {
		return core.EmptyString, err
	}

	if found == nil {
		return core.EmptyString, drivers.ErrNotFound
	}

	return found.GetInnerText(ctx)
}

func (el *HTMLElement) SetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector, innerText core.String) error {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return drivers.ErrNotFound
		}

		selection.SetHtml(innerText.String())

		return nil
	}

	found, err := EvalXPathToElement(el.selection, selector.String())

	if err != nil {
		return err
	}

	if found == nil {
		return drivers.ErrNotFound
	}

	return found.SetInnerText(ctx, innerText)
}

func (el *HTMLElement) GetInnerTextBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*internal.Array, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())
		arr := internal.NewArray(selection.Length())

		selection.Each(func(_ int, selection *goquery.Selection) {
			arr.Push(core.NewString(selection.Text()))
		})

		return arr, nil
	}

	return EvalXPathToNodesWith(el.selection, selector.String(), func(node *html.Node) (core.Value, error) {
		n, err := parseXPathNode(node)

		if err != nil {
			return core.None, err
		}

		found, err := drivers.ToElement(n)

		if err != nil {
			return core.None, err
		}

		return found.GetInnerText(ctx)
	})
}

func (el *HTMLElement) CountBySelector(_ context.Context, selector drivers.QuerySelector) (core.Int, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		return core.NewInt(selection.Length()), nil
	}

	arr, err := EvalXPathToNodesWith(el.selection, selector.String(), func(node *html.Node) (core.Value, error) {
		return core.None, nil
	})

	if err != nil {
		return core.ZeroInt, err
	}

	return arr.Length(), nil
}

func (el *HTMLElement) ExistsBySelector(_ context.Context, selector drivers.QuerySelector) (core.Boolean, error) {
	if selector.Kind() == drivers.CSSSelector {
		selection := el.selection.Find(selector.String())

		if selection.Length() == 0 {
			return core.False, nil
		}

		return core.True, nil
	}

	found, err := EvalXPathToNode(el.selection, selector.String())

	if err != nil {
		return core.False, err
	}

	return core.NewBoolean(found != nil), nil
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInElement(ctx, path, el)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInElement(ctx, path, el, value)
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) GetParentElement(_ context.Context) (core.Value, error) {
	parent := el.selection.Parent()

	if parent == nil {
		return core.None, nil
	}

	return NewHTMLElement(parent)
}

func (el *HTMLElement) GetPreviousElementSibling(_ context.Context) (core.Value, error) {
	sibling := el.selection.Prev()

	if sibling == nil {
		return core.None, nil
	}

	return NewHTMLElement(sibling)
}

func (el *HTMLElement) GetNextElementSibling(_ context.Context) (core.Value, error) {
	sibling := el.selection.Next()

	if sibling == nil {
		return core.None, nil
	}

	return NewHTMLElement(sibling)
}

func (el *HTMLElement) Click(_ context.Context, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClickBySelector(_ context.Context, _ drivers.QuerySelector, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClickBySelectorAll(_ context.Context, _ drivers.QuerySelector, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Clear(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) ClearBySelector(_ context.Context, _ drivers.QuerySelector) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Input(_ context.Context, _ core.Value, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) InputBySelector(_ context.Context, _ drivers.QuerySelector, _ core.Value, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Press(_ context.Context, _ []core.String, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) PressBySelector(_ context.Context, _ drivers.QuerySelector, _ []core.String, _ core.Int) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Select(_ context.Context, _ *internal.Array) (*internal.Array, error) {
	return nil, core.ErrNotSupported
}

func (el *HTMLElement) SelectBySelector(_ context.Context, _ drivers.QuerySelector, _ *internal.Array) (*internal.Array, error) {
	return nil, core.ErrNotSupported
}

func (el *HTMLElement) ScrollIntoView(_ context.Context, _ drivers.ScrollOptions) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Focus(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) FocusBySelector(_ context.Context, _ drivers.QuerySelector) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Blur(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) BlurBySelector(_ context.Context, _ drivers.QuerySelector) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) Hover(_ context.Context) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) HoverBySelector(_ context.Context, _ drivers.QuerySelector) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForClass(_ context.Context, _ core.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForElement(_ context.Context, _ drivers.QuerySelector, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForElementAll(_ context.Context, _ drivers.QuerySelector, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForAttribute(_ context.Context, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForAttributeBySelector(_ context.Context, _ drivers.QuerySelector, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForAttributeBySelectorAll(_ context.Context, _ drivers.QuerySelector, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForStyle(_ context.Context, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForStyleBySelector(_ context.Context, _ drivers.QuerySelector, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForStyleBySelectorAll(_ context.Context, _ drivers.QuerySelector, _ core.String, _ core.Value, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForClassBySelector(_ context.Context, _ drivers.QuerySelector, _ core.String, _ drivers.WaitEvent) error {
	return core.ErrNotSupported
}

func (el *HTMLElement) WaitForClassBySelectorAll(_ context.Context, _ drivers.QuerySelector, _ core.String, _ drivers.WaitEvent) error {
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

func (el *HTMLElement) parseStyles(ctx context.Context) (*internal.Object, error) {
	str, err := el.GetAttribute(ctx, "style")

	if err != nil {
		return internal.NewObject(), err
	}

	if str == core.None {
		return internal.NewObject(), nil
	}

	styles, err := common.DeserializeStyles(core.NewString(str.String()))

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

func (el *HTMLElement) parseAttrs() *internal.Object {
	obj := internal.NewObject()

	if len(el.selection.Nodes) == 0 {
		return obj
	}

	node := el.selection.Nodes[0]

	for _, attr := range node.Attr {
		obj.Set(core.NewString(attr.Key), core.NewString(attr.Val))
	}

	return obj
}

func (el *HTMLElement) parseChildren() *internal.Array {
	children := el.selection.Children()

	arr := internal.NewArray(10)

	children.Each(func(i int, selection *goquery.Selection) {
		child, err := NewHTMLElement(selection)

		if err == nil {
			arr.Push(child)
		}
	})

	return arr
}
