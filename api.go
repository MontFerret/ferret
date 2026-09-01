package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	// Source represents the input data for a Ferret query.
	Source = source.Source

	// Position represents a specific point in a source file, defined by line and column numbers.
	Position = source.Position

	// Span represents a range of characters in a source file, defined by start and end positions.
	Span = source.Span

	// Location represents the location of a specific point in a source file, including the file name and position.
	Location = source.Location

	// Range represents a range of characters in a source file, including the location and span.
	Range = source.Range

	// Output is the encoded result returned from session or engine execution.
	Output = encoding.Output

	// OptimizationLevel represents the level of optimization applied during query compilation.
	OptimizationLevel = compiler.OptimizationLevel
)

// NewSource creates a new Source instance with the given name and content.
func NewSource(name, content string) Source {
	return source.New(name, content)
}

// NewAnonymousSource creates a new anonymous Source instance with the given content.
func NewAnonymousSource(content string) Source {
	return source.NewAnonymous(content)
}

// WithOptimizationLevel returns a compiler option that sets the optimization level for query compilation.
func WithOptimizationLevel(level OptimizationLevel) compiler.Option {
	return compiler.WithOptimizationLevel(level)
}
