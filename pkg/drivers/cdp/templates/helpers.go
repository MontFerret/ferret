package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

var (
	notFoundErrorFragment = fmt.Sprintf(`
		if (found == null) {
			throw new Error(%s);
		}
`, ParamErr(drivers.ErrNotFound))
)

const (
	toElementFragment = `
(input) => {
	let result = null;

	if (input instanceof Element) {
		result = input;
	} else if (Array.isArray(input) && input[0] instanceof Element) {
		result = input[0]
	} else if (input instanceof NodeList) {
		result = input[0]
	}

	return result;
}
`
	toElementArrayFragment = `
(input) => {
	if (input == null) {
		return [];
	}

	if (Array.isArray(input)) {
		return input.filter(i => i instanceof Element);
	}

	if (input instanceof NodeList) {
		return Array.from(input)
	}
	
	return [input];
}
`
)

func ParamErr(err error) string {
	return EscapeString(err.Error())
}

func EscapeString(value string) string {
	return "`" + value + "`"
}

func toFunction(selector drivers.QuerySelector, cssTmpl, xPathTmpl string) *eval.Function {
	if selector.Variant() == drivers.CSSSelector {
		return eval.F(cssTmpl)
	}

	return eval.F(xPathTmpl)
}
