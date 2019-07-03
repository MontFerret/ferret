package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"

	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func AttributeRead(name values.String) string {
	n := name.String()
	return fmt.Sprintf(`
		el.attributes[%s] != null ? el.attributes[%s].value : null
	`, eval.ParamString(n), eval.ParamString(n))
}
