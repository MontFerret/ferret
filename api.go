package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	Source = source.Source

	// Output is the encoded result returned from session or engine execution.
	Output = encoding.Output
)
