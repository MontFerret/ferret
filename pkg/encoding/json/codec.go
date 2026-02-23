package json

import (
	"bytes"
	stdjson "encoding/json"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Codec encodes and decodes runtime values as JSON.
type Codec struct{}

// Default is the default JSON codec.
var Default Codec = Codec{}

func (Codec) Encode(value runtime.Value) ([]byte, error) {
	if value == nil {
		value = runtime.None
	}

	var buf bytes.Buffer

	encoder := stdjson.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	out := buf.Bytes()
	if size := len(out); size > 0 && out[size-1] == '\n' {
		out = out[:size-1]
	}

	return out, nil
}

func (Codec) Decode(data []byte) (runtime.Value, error) {
	return runtime.Unmarshal(data)
}
