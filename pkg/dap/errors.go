package dap

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func formatError(err error) string {
	if err == nil {
		return ""
	}

	return diagnostics.Format(err)
}

func wrapStateError(message string) error {
	return fmt.Errorf("dap adapter: %s", message)
}
