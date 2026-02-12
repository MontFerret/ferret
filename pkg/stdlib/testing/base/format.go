package base

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func FormatValue(val runtime.Value) string {
	valStr := val.String()

	if val == runtime.None {
		valStr = "none"
	}

	return fmt.Sprintf("[%s] '%s'", runtime.Reflect(val), valStr)
}
