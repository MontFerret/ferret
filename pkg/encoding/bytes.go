package encoding

import (
	"bytes"
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// EncodeBytes collects an encoder's writes into caller-owned bytes, including
// any prefix written before an error. Publish the bytes only on success when
// whole-document output must be atomic.
func EncodeBytes(ctx context.Context, enc Encoder, value runtime.Value) ([]byte, error) {
	var dst bytes.Buffer
	err := enc.Encode(ctx, &dst, value)

	return dst.Bytes(), err
}

// DecodeBytes decodes a complete byte document through the reader-based API.
func DecodeBytes(ctx context.Context, dec Decoder, data []byte) (runtime.Value, error) {
	return dec.Decode(ctx, bytes.NewReader(data))
}
