package ferret

import (
	gooptions "github.com/ziflex/go-options"

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
)

// NewSource creates a new Source instance with the given name and content.
func NewSource(name, content string) Source {
	return source.New(name, content)
}

// NewAnonymousSource creates a new anonymous Source instance with the given content.
func NewAnonymousSource(content string) Source {
	return source.NewAnonymous(content)
}

// WithOptimizationLevel configures optimization for normal query compilation.
// OptimizationFull is used by default; debug compilation always uses OptimizationNone.
func WithOptimizationLevel(level OptimizationLevel) Option {
	return gooptions.New(
		func(config *config, level OptimizationLevel) {
			config.optimizationLevel = level
		},
	).
		Named("optimization level").
		Validators(
			gooptions.OneOf(
				OptimizationNone,
				OptimizationBasic,
				OptimizationFull,
			),
		).
		Value(level).
		Build()
}
