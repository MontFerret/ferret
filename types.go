package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	ParseLogLevel     = runtime.ParseLogLevel
	MustParseLogLevel = runtime.MustParseLogLevel
	FormatError       = diagnostics.Format
)
