package encoding

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type (
	// Encoder converts runtime values into bytes.
	Encoder interface {
		Encode(value runtime.Value) ([]byte, error)
	}

	// Decoder converts bytes into runtime values.
	Decoder interface {
		Decode(data []byte) (runtime.Value, error)
	}

	// Codec combines encoder and decoder capabilities.
	Codec interface {
		Encoder
		Decoder
	}
)
