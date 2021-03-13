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

func SetDefaultParams(opts *Options, params *Params) {
	if params.Headers == nil && opts.Headers != nil {
		params.Headers = NewHTTPHeaders()
	}

	// set default headers
	if opts.Headers != nil {
		opts.Headers.ForEach(func(value []string, key string) bool {
			val := params.Headers.Get(key)

			// do not override user's set values
			if val == "" {
				params.Headers.SetArr(key, value)
			}

			return true
		})
	}

	if params.Cookies == nil && opts.Cookies != nil {
		params.Cookies = make(HTTPCookies)
	}

	// set default cookies
	if opts.Cookies != nil {
		for k, v := range opts.Cookies {
			_, exists := params.Cookies[k]

			// do not override user's set values
			if !exists {
				params.Cookies[k] = v
			}
		}
	}

	// set default user agent
	if opts.UserAgent != "" && params.UserAgent == "" {
		params.UserAgent = opts.UserAgent
	}
}
