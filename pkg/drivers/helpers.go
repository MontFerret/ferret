package drivers

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func ToPage(value core.Value) (HTMLPage, error) {
	err := core.ValidateType(value, HTMLPageType)

	if err != nil {
		return nil, err
	}

	return value.(HTMLPage), nil
}

func ToDocument(value core.Value) (HTMLDocument, error) {
	switch v := value.(type) {
	case HTMLPage:
		return v.GetMainFrame(), nil
	case HTMLDocument:
		return v, nil
	default:
		return nil, core.TypeError(
			value.Type(),
			HTMLPageType,
			HTMLDocumentType,
		)
	}
}

func ToElement(value core.Value) (HTMLElement, error) {
	switch v := value.(type) {
	case HTMLPage:
		return v.GetMainFrame().GetElement(), nil
	case HTMLDocument:
		return v.GetElement(), nil
	case HTMLElement:
		return v, nil
	default:
		return nil, core.TypeError(
			value.Type(),
			HTMLPageType,
			HTMLDocumentType,
			HTMLElementType,
		)
	}
}
