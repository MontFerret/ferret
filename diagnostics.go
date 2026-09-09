package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

// FormatError formats a diagnostic error into a human-readable string.
func FormatError(err error) string {
	return diagnostics.Format(err)
}
