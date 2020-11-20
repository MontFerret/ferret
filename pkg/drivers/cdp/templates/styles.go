package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

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
