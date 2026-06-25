package artifact

import programformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format"

type (
	// FormatID identifies a serialized program payload format in an artifact
	// header.
	FormatID uint8

	// options controls artifact marshaling.
	options struct {
		Format FormatID
	}

	// Option controls artifact marshaling.
	Option func(*options)

	// RegisteredFormat associates an artifact format id with a payload format
	// implementation.
	RegisteredFormat struct {
		Format programformat.Format
		ID     FormatID
	}
)

const (
	FormatJSON    FormatID = 1
	FormatMsgPack FormatID = 2
	DefaultFormat          = FormatMsgPack
)

func newOptions(setters ...Option) options {
	opts := options{
		Format: DefaultFormat,
	}

	for _, setter := range setters {
		setter(&opts)
	}

	return opts
}

func WithFormat(format FormatID) Option {
	return func(o *options) {
		o.Format = format
	}
}
