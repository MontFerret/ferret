package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"ATTR_GET":          AttributeGet,
			"ATTR_QUERY":        AttributeQuery,
			"ATTR_REMOVE":       AttributeRemove,
			"ATTR_SET":          AttributeSet,
			"BLUR":              Blur,
			"COOKIE_DEL":        CookieDel,
			"COOKIE_GET":        CookieGet,
			"COOKIE_SET":        CookieSet,
			"CLICK":             Click,
			"CLICK_ALL":         ClickAll,
			"DOCUMENT":          Open,
			"DOCUMENT_EXISTS":   DocumentExists,
			"DOWNLOAD":          Download,
			"ELEMENT":           Element,
			"ELEMENT_EXISTS":    ElementExists,
			"ELEMENTS":          Elements,
			"ELEMENTS_COUNT":    ElementsCount,
			"FRAMES":            Frames,
			"FOCUS":             Focus,
			"HOVER":             Hover,
			"INNER_HTML":        GetInnerHTML,
			"INNER_HTML_SET":    SetInnerHTML,
			"INNER_HTML_ALL":    GetInnerHTMLAll,
			"INNER_TEXT":        GetInnerText,
			"INNER_TEXT_SET":    SetInnerText,
			"INNER_TEXT_ALL":    GetInnerTextAll,
			"INPUT":             Input,
			"INPUT_CLEAR":       InputClear,
			"MOUSE":             MouseMoveXY,
			"NAVIGATE":          Navigate,
			"NAVIGATE_BACK":     NavigateBack,
			"NAVIGATE_FORWARD":  NavigateForward,
			"PAGINATION":        Pagination,
			"PARSE":             Parse,
			"PDF":               PDF,
			"PRESS":             Press,
			"PRESS_SELECTOR":    PressSelector,
			"SCREENSHOT":        Screenshot,
			"SCROLL":            ScrollXY,
			"SCROLL_BOTTOM":     ScrollBottom,
			"SCROLL_ELEMENT":    ScrollInto,
			"SCROLL_TOP":        ScrollTop,
			"SELECT":            Select,
			"STYLE_GET":         StyleGet,
			"STYLE_REMOVE":      StyleRemove,
			"STYLE_SET":         StyleSet,
			"WAIT_ATTR":         WaitAttribute,
			"WAIT_NO_ATTR":      WaitNoAttribute,
			"WAIT_ATTR_ALL":     WaitAttributeAll,
			"WAIT_NO_ATTR_ALL":  WaitNoAttributeAll,
			"WAIT_ELEMENT":      WaitElement,
			"WAIT_NO_ELEMENT":   WaitNoElement,
			"WAIT_CLASS":        WaitClass,
			"WAIT_NO_CLASS":     WaitNoClass,
			"WAIT_CLASS_ALL":    WaitClassAll,
			"WAIT_NO_CLASS_ALL": WaitNoClassAll,
			"WAIT_STYLE":        WaitStyle,
			"WAIT_NO_STYLE":     WaitNoStyle,
			"WAIT_STYLE_ALL":    WaitStyleAll,
			"WAIT_NO_STYLE_ALL": WaitNoStyleAll,
			"WAIT_NAVIGATION":   WaitNavigation,
			"XPATH":             XPath,
			"X":                 XPathSelector,
		}))
}

func OpenOrCastPage(ctx context.Context, value core.Value) (drivers.HTMLPage, bool, error) {
	err := core.ValidateType(value, drivers.HTMLPageType, types.String)

	if err != nil {
		return nil, false, err
	}

	var page drivers.HTMLPage
	var closeAfter bool

	if value.Type() == types.String {
		buf, err := Open(ctx, value, values.NewBoolean(true))

		if err != nil {
			return nil, false, err
		}

		page = buf.(drivers.HTMLPage)
		closeAfter = true
	} else {
		page = value.(drivers.HTMLPage)
	}

	return page, closeAfter, nil
}

func waitTimeout(ctx context.Context, value values.Int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		ctx,
		time.Duration(value)*time.Millisecond,
	)
}

func toScrollOptions(value core.Value) (drivers.ScrollOptions, error) {
	result := drivers.ScrollOptions{}

	err := core.ValidateType(value, types.Object)

	if err != nil {
		return result, err
	}

	obj := value.(*values.Object)

	behavior, exists := obj.Get("behavior")

	if exists {
		result.Behavior = drivers.NewScrollBehavior(behavior.String())
	}

	block, exists := obj.Get("block")

	if exists {
		result.Block = drivers.NewScrollVerticalAlignment(block.String())
	}

	inline, exists := obj.Get("inline")

	if exists {
		result.Inline = drivers.NewScrollHorizontalAlignment(inline.String())
	}

	return result, nil
}
