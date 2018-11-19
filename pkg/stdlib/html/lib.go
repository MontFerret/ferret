package html

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
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
		"DOCUMENT_PARSE":   DocumentParse,
		"DOWNLOAD":         Download,
		"ELEMENT":          Element,
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
