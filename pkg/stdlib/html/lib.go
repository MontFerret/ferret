package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

const defaultTimeout = 5000

var (
	ErrNotDynamic = errors.New("expected dynamic document or element")
)

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"CLICK":            Click,
		"CLICK_ALL":        ClickAll,
		"DOCUMENT":         Document,
		"DOWNLOAD":         Download,
		"ELEMENT":          Element,
		"ELEMENT_EXISTS":   ElementExists,
		"ELEMENTS":         Elements,
		"ELEMENTS_COUNT":   ElementsCount,
		"HOVER":            Hover,
		"INNER_HTML":       InnerHTML,
		"INNER_HTML_ALL":   InnerHTMLAll,
		"INNER_TEXT":       InnerText,
		"INNER_TEXT_ALL":   InnerTextAll,
		"INPUT":            Input,
		"NAVIGATE":         Navigate,
		"NAVIGATE_BACK":    NavigateBack,
		"NAVIGATE_FORWARD": NavigateForward,
		"PAGINATION":       Pagination,
		"PDF":              PDF,
		"SCREENSHOT":       Screenshot,
		"SCROLL_BOTTOM":    ScrollBottom,
		"SCROLL_ELEMENT":   ScrollInto,
		"SCROLL_TOP":       ScrollTop,
		"SELECT":           Select,
		"WAIT_ELEMENT":     WaitElement,
		"WAIT_CLASS":       WaitClass,
		"WAIT_CLASS_ALL":   WaitClassAll,
		"WAIT_NAVIGATION":  WaitNavigation,
	}
}

func ValidateDocument(ctx context.Context, value core.Value) (core.Value, error) {
	err := core.ValidateType(value, drivers.HTMLDocumentType, types.String)
	if err != nil {
		return values.None, err
	}

	var doc drivers.HTMLDocument
	var ok bool

	if value.Type() == types.String {
		buf, err := Document(ctx, value, values.NewBoolean(true))

		if err != nil {
			return values.None, err
		}

		doc, ok = buf.(drivers.HTMLDocument)
	} else {
		doc, ok = value.(drivers.HTMLDocument)
	}

	if !ok {
		return nil, ErrNotDynamic
	}

	return doc, nil
}
