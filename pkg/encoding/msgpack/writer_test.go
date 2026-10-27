package msgpack_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

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

func TestMsgpackCustomEncoderCanEncodeRecursively(t *testing.T) {
	value := &reentrantArrayEncoder{Array: runtime.NewArrayWith(runtime.Int(99)), ctx: t.Context()}
	for i := 0; i < 3; i++ {
		data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
		if err != nil || !bytes.Equal(data, []byte{0x92, 1, 2}) {
			t.Fatalf("custom list encoding: %x %v", data, err)
		}
	}
}

func TestMsgpackPayloadChunkFailureAndCancellation(t *testing.T) {
	const chunkSize = 4096
	payload := runtime.Binary(bytes.Repeat([]byte{0xab}, 3*chunkSize+17))
	value := runtime.NewArrayWith(payload)
	golden, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
	if err != nil {
		t.Fatal(err)
	}

	// The one-element array's headers precede its binary payload; no suffix follows it.
	headerSize := len(golden) - len(payload)
	failure := errors.New("payload write failed")
	for _, test := range []struct {
		failure error
		want    error
		name    string
		cancel  bool
	}{
		{name: "second_chunk_error", failure: failure, want: failure},
		{name: "second_chunk_short_write", want: io.ErrShortWrite},
		{name: "cancel_after_first_chunk", want: context.Canceled, cancel: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			writer := &payloadBoundaryWriter{failure: test.failure, stopAt: headerSize + chunkSize + 7}
			if test.cancel {
				writer.stopAt = headerSize + chunkSize
				writer.cancel = cancel
			}

			var ancestorErr error
			ancestorPosts := 0
			enc := ferretmsgpack.Default.EncodeWith().PostHook(func(actual context.Context, item runtime.Value, err error) error {
				if actual != ctx {
					t.Error("post-hook lost the caller's context")
				}

				if item == value {
					ancestorErr = err
					ancestorPosts++
				}

				return nil
			}).Encoder()
			err := enc.Encode(ctx, writer, value)
			if !errors.Is(err, test.want) || !errors.Is(ancestorErr, test.want) || ancestorPosts != 1 {
				t.Fatalf("lost payload failure: result=%v ancestor=%v posts=%d want=%v", err, ancestorErr, ancestorPosts, test.want)
			}

			if !writer.stopped || writer.writesAfterStop != 0 || !bytes.Equal(writer.Bytes(), golden[:writer.stopAt]) {
				t.Fatalf("payload continued or prefix changed: stopped=%t late writes=%d length=%d want=%d", writer.stopped, writer.writesAfterStop, writer.Len(), writer.stopAt)
			}

			if test.cancel && ctx.Err() != context.Canceled {
				t.Fatalf("caller context was not canceled: %v", ctx.Err())
			}
		})
	}
}
