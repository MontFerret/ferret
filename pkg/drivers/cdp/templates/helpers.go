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

func ParamErr(err error) string {
	return EscapeString(err.Error())
}

func EscapeString(value string) string {
	return "`" + value + "`"
}

func toFunction(selector drivers.QuerySelector, cssTmpl, xPathTmpl string) *eval.Function {
	if selector.Kind() == drivers.CSSSelector {
		return eval.F(cssTmpl)
	}

	return eval.F(xPathTmpl)
}
