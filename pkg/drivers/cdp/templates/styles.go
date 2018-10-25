package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func StyleRead(name values.String) string {
	n := name.String()
	return fmt.Sprintf(
		`el.style[%s] != "" ? el.style[%s] : null`,
		eval.ParamString(n),
		eval.ParamString(n),
	)
}
