package base

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func FormatValue(val core.Value) string {
	valStr := val.String()

	if val == core.None {
		valStr = "none"
	}

	return fmt.Sprintf("[%s] '%s'", core.Reflect(val), valStr)
}
