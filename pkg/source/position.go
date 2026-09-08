package source

import apisource "github.com/MontFerret/api/source"

// Native source positions use one-based lines and UTF-8 byte columns. Spans
// use zero-based byte offsets with an exclusive end. The portable coordinate
// values carry these conventions without duplicating their representation.
type (
	Position = apisource.Position
	Span     = apisource.Span
	Location = apisource.Location
	Range    = apisource.Range
)
