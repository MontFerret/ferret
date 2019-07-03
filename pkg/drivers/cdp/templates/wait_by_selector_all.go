package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func WaitBySelectorAll(selector values.String, when drivers.WaitEvent, value core.Value, check string) string {
	return fmt.Sprintf(`
			var elements = document.querySelectorAll(%s); // selector
			
			if (elements == null || elements.length === 0) {
				return false;
			}

			var resultCount = 0;
			
			elements.forEach((el) => {
				var result = %s; // check

				// when
				if (result %s %s) {
					resultCount++;
				}
			});
	
			if (resultCount === elements.length) {
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
