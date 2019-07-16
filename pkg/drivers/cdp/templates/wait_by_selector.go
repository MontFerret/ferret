package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func WaitBySelector(selector values.String, when drivers.WaitEvent, value core.Value, check string) string {
	return fmt.Sprintf(
		`
			const el = document.querySelector(%s); // selector
			
			if (el == null) {
				return false;
			}

			const result = %s; // check

			// when value
			if (result %s %s) {
				return true;
			}
			
			// null means we need to repeat
			return null;
`,
		eval.ParamString(selector.String()),
		check,
		WaitEventToEqOperator(when),
		eval.Param(value),
	)
}
