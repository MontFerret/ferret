package base

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func FormatValue(val runtime.Value) string {
	valStr := val.String()

	if val == runtime.None {
		valStr = "none"
	}

	valStr = strings.ReplaceAll(valStr, `\`, `\\`)
	valStr = strings.ReplaceAll(valStr, `'`, `\'`)

	return fmt.Sprintf("%s '%s'", runtime.TypeOf(val), valStr)
}
