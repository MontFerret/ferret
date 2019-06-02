package html

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

const defaultTimeout = 5000

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"ATTR_GET":          AttributeGet,
		"ATTR_REMOVE":       AttributeRemove,
		"ATTR_SET":          AttributeSet,
		"COOKIE_DEL":        CookieDel,
		"COOKIE_GET":        CookieGet,
		"COOKIE_SET":        CookieSet,
		"CLICK":             Click,
		"CLICK_ALL":         ClickAll,
		"DOCUMENT":          Open,
		"DOWNLOAD":          Download,
		"ELEMENT":           Element,
		"ELEMENT_EXISTS":    ElementExists,
		"ELEMENTS":          Elements,
		"ELEMENTS_COUNT":    ElementsCount,
		"HOVER":             Hover,
		"INNER_HTML":        InnerHTML,
		"INNER_HTML_ALL":    InnerHTMLAll,
		"INNER_TEXT":        InnerText,
		"INNER_TEXT_ALL":    InnerTextAll,
		"INPUT":             Input,
		"MOUSE":             MouseMoveXY,
		"NAVIGATE":          Navigate,
		"NAVIGATE_BACK":     NavigateBack,
		"NAVIGATE_FORWARD":  NavigateForward,
		"PAGINATION":        Pagination,
		"PDF":               PDF,
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
	}
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

func resolveElement(value core.Value) (drivers.HTMLElement, error) {
	vt := value.Type()

	if vt == drivers.HTMLDocumentType {
		return value.(drivers.HTMLDocument).Element(), nil
	} else if vt == drivers.HTMLElementType {
		return value.(drivers.HTMLElement), nil
	}

	return nil, core.TypeError(value.Type(), drivers.HTMLDocumentType, drivers.HTMLElementType)
}

func toPage(value core.Value) (drivers.HTMLPage, error) {
	err := core.ValidateType(value, drivers.HTMLPageType)

	if err != nil {
		return nil, err
	}

	return value.(drivers.HTMLPage), nil
}

func toDocument(value core.Value) (drivers.HTMLDocument, error) {
	err := core.ValidateType(value, drivers.HTMLDocumentType)

	if err != nil {
		return nil, err
	}

	return value.(drivers.HTMLDocument), nil
}

func toElement(value core.Value) (drivers.HTMLElement, error) {
	err := core.ValidateType(value, drivers.HTMLElementType)

	if err != nil {
		return nil, err
	}

	return value.(drivers.HTMLElement), nil
}
