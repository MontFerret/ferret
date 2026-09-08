package ferret

import "github.com/MontFerret/ferret/v2/pkg/engine"

type (
	// Source retains source text and its name for compilation and diagnostics.
	Source = engine.Source

	// Position identifies a one-based line and byte column in source text.
	Position = engine.Position

	// Span identifies a zero-based, half-open byte interval in source text.
	Span = engine.Span

	// Location identifies a position in a named source.
	Location = engine.Location

	// Range identifies a range of positions in a named source.
	Range = engine.Range
)

// NewSource creates a named source for compilation and diagnostics.
func NewSource(name, content string) Source {
	return engine.NewSource(name, content)
}

// NewAnonymousSource creates a source without a name.
func NewAnonymousSource(content string) Source {
	return engine.NewAnonymousSource(content)
}
