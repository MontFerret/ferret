package msgpack_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type ignoringWriterValue struct {
	marshalerValue
	failure error
}

func (v *ignoringWriterValue) EncodeMsgpack(enc *vmmsgpack.Encoder) error {
	enc.UseCompactInts(true)
	_, _ = enc.Writer().Write([]byte{1, 2, 3})
	return v.failure
}

type errorWriter struct{ failure error }

func (w errorWriter) Write([]byte) (int, error) { return 0, w.failure }

func TestMsgpackCustomEncoderPreservesIgnoredWriterErrors(t *testing.T) {
	writeErr, customErr := errors.New("write failed"), errors.New("custom encoding failed")
	for _, returned := range []error{nil, customErr} {
		var postErr error
		enc := ferretmsgpack.Default.EncodeWith().PostHook(func(_ context.Context, _ runtime.Value, err error) error {
			postErr = err
			return nil
		}).Encoder()
		err := enc.Encode(t.Context(), errorWriter{failure: writeErr}, &ignoringWriterValue{failure: returned})
		if !errors.Is(err, writeErr) || !errors.Is(postErr, writeErr) || (returned != nil && !errors.Is(err, returned)) {
			t.Fatalf("lost error: result=%v post=%v", err, postErr)
		}
	}
	// Flags changed by custom encoding do not leak into the next pooled operation.
	data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, &unwrapValue{value: int64(1)})
	if err != nil || !bytes.Equal(data, []byte{0xd3, 0, 0, 0, 0, 0, 0, 0, 1}) {
		t.Fatalf("pool reset: %x %v", data, err)
	}
}
