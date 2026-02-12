package data

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Transformer interface {
	runtime.Value
	runtime.Iterable
	runtime.KeyReadable
	runtime.KeyWritable
	runtime.Measurable
	io.Closer
}
