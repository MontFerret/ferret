package data

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type Transformer interface {
	runtime.Value
	runtime.Iterable
	runtime.Measurable
	runtime.KeyReadable
	runtime.KeyWritable
	io.Closer
}
