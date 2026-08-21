package testing

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func formatValue(value runtime.Value) string {
	valueString := value.String()
	if value == runtime.None {
		valueString = "none"
	}

	valueString = strings.ReplaceAll(valueString, `\`, `\\`)
	valueString = strings.ReplaceAll(valueString, `'`, `\'`)

	return fmt.Sprintf("%s '%s'", runtime.TypeOf(value), valueString)
}
