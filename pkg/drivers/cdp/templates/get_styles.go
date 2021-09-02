package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

var getStylesTemplate = `
	(el) => {
		const out = {};
		const styles = window.getComputedStyle(el);
	
		Object.keys(styles).forEach((key) => {
			if (!isNaN(parseFloat(key))) {
				const name = styles[key];
				const value = styles.getPropertyValue(name);
				out[name] = value;
			}
		});

		return out;
	}
`

func GetStyles() string {
	return getStylesTemplate
}

func WaitForStyle(name, value string, when drivers.WaitEvent) string {
	return fmt.Sprintf(`
	(el) => {
		const styles = window.getComputedStyle(el);
		const actual = styles[%s];
		const expected = %s;

		// null means we need to repeat
		return actual %s expected ? true : null ;
	}
`, eval.ParamString(name), eval.ParamString(value), WaitEventToEqOperator(when))
}
