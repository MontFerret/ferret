package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/logging"
)

var (
	// ParseLogLevel parses a log level string and returns the corresponding LogLevel and an error if the string is invalid.
	ParseLogLevel = logging.ParseLogLevel
	// MustParseLogLevel parses a log level string and panics if the string is not a valid log level.
	MustParseLogLevel = logging.MustParseLogLevel
	// FormatError formats a diagnostic error into a human-readable string.
	FormatError = diagnostics.Format
)
