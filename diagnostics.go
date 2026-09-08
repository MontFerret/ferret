package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

// FormatError formats a diagnostic error into a human-readable string.
func FormatError(err error) string {
	return engine.FormatError(err)
}
