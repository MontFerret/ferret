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
		"CLICK":             Click,
		"CLICK_ALL":         ClickAll,
		"DOCUMENT":          Document,
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
		"WAIT_ELEMENT":      WaitElement,
		"WAIT_NO_ELEMENT":   WaitNoElement,
		"WAIT_CLASS":        WaitClass,
		"WAIT_NO_CLASS":     WaitNoClass,
		"WAIT_CLASS_ALL":    WaitClassAll,
		"WAIT_NO_CLASS_ALL": WaitNoClassAll,
		"WAIT_NAVIGATION":   WaitNavigation,
	}
}

func ValidateDocument(ctx context.Context, value core.Value) (core.Value, error) {
	err := core.ValidateType(value, drivers.HTMLDocumentType, types.String)
	if err != nil {
		return values.None, err
	}

	var doc drivers.HTMLDocument

	if value.Type() == types.String {
		buf, err := Document(ctx, value, values.NewBoolean(true))

		if err != nil {
			return values.None, err
		}

		doc = buf.(drivers.HTMLDocument)
	} else {
		doc = value.(drivers.HTMLDocument)
	}

	return doc, nil
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
		return value.(drivers.HTMLDocument).DocumentElement(), nil
	} else if vt == drivers.HTMLElementType {
		return value.(drivers.HTMLElement), nil
	}

	return nil, core.TypeError(value.Type(), drivers.HTMLDocumentType, drivers.HTMLElementType)
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
