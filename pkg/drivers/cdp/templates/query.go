package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

func QuerySelector(selector string) string {
	return fmt.Sprintf(`
		(el) => {
			const found = el.querySelector(%s);
	
			if (found == null) {
				throw new Error(%s);
			}
	
			return found;
		}
	`,
		eval.ParamString(selector),
		eval.ParamString(drivers.ErrNotFound.Error()),
	)
}

func QuerySelectorAll(selector string) string {
	return fmt.Sprintf(`
		(el) => {
			const found = el.querySelectorAll(%s);
	
			if (found == null) {
				throw new Error(%s);
			}
	
			return found;
		}
	`,
		eval.ParamString(selector),
		eval.ParamString(drivers.ErrNotFound.Error()),
	)
}
