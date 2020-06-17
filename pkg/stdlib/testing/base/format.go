package base

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func FormatValue(val core.Value) string {
	valStr := val.String()

	if val == values.None {
		valStr = "none"
	}

	return fmt.Sprintf("[%s] '%s'", val.Type(), valStr)
}
