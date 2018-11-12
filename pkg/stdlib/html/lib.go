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
		"DOCUMENT":         Document,
		"DOCUMENT_PARSE":   DocumentParse,
		"ELEMENT":          Element,
		"ELEMENTS":         Elements,
		"ELEMENTS_COUNT":   ElementsCount,
		"WAIT_ELEMENT":     WaitElement,
		"WAIT_NAVIGATION":  WaitNavigation,
		"WAIT_CLASS":       WaitClass,
		"WAIT_CLASS_ALL":   WaitClassAll,
		"CLICK":            Click,
		"CLICK_ALL":        ClickAll,
		"NAVIGATE":         Navigate,
		"NAVIGATE_BACK":    NavigateBack,
		"NAVIGATE_FORWARD": NavigateForward,
		"INPUT":            Input,
		"INNER_HTML":       InnerHTML,
		"INNER_HTML_ALL":   InnerHTMLAll,
		"INNER_TEXT":       InnerText,
		"INNER_TEXT_ALL":   InnerTextAll,
		"SELECT":           Select,
		"SCREENSHOT":       Screenshot,
		"PDF":              PDF,
		"DOWNLOAD":         Download,
	}
}
