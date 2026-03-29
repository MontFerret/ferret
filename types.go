package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/logging"
)

var (
	ParseLogLevel     = logging.ParseLogLevel
	MustParseLogLevel = logging.MustParseLogLevel
	FormatError       = diagnostics.Format
)
