package ferret

import "github.com/MontFerret/ferret/v2/pkg/source"

type (
	// Source retains source text and its name for compilation and diagnostics.
	Source = source.Source

	// Position identifies a one-based line and byte column in source text.
	Position = source.Position

	// Span identifies a zero-based, half-open byte interval in source text.
	Span = source.Span

	// Location identifies a position in a named source.
	Location = source.Location

	// Range identifies a range of positions in a named source.
	Range = source.Range
)

// NewSource creates a named source for compilation and diagnostics.
func NewSource(name, content string) Source {
	return source.New(name, content)
}

// NewAnonymousSource creates a source without a name.
func NewAnonymousSource(content string) Source {
	return source.NewAnonymous(content)
}
