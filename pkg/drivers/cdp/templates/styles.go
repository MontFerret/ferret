package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

func GetStyle(name string) string {
	return fmt.Sprintf(`
	(el) => {
		const out = {};
		const styles = window.getComputedStyle(el);

		return styles[%s];
	}
`, eval.ParamString(name))
}

func SetStyle(name, value string) string {
	return fmt.Sprintf(`
		(el) => {
			el.style[%s] = %s;
		}
`, eval.ParamString(name), eval.ParamString(value))
}

func RemoveStyles(names []values.String) string {
	return fmt.Sprintf(`
		(el) => {
			const style = el.style;
			[%s].forEach((name) => { style[name] = "" })
		}
	`,
		eval.ParamStringList(names),
	)
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

func StyleRead(name values.String) string {
	n := name.String()
	return fmt.Sprintf(`
	((function() {
		const cs = window.getComputedStyle(el);
		const currentValue = cs.getPropertyValue(%s);

		return currentValue || null;
	})())
`, eval.ParamString(n))
}
