package msgpack_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMsgpackAdoptionCallbackIsOperationScoped(t *testing.T) {
	for _, failed := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failure"}[failed], func(t *testing.T) {
			codec := ferretmsgpack.Default.EncodeWith().PreHook(func(context.Context, runtime.Value) error { return nil }).Encoder()
			calls := 0
			ctx := encodingownership.WithValueAdopter(t.Context(), func(value runtime.Value) {
				calls++
				if value != runtime.Int(1) {
					t.Fatalf("adopted wrong value: %v", value)
				}
			})
			dst := io.Writer(io.Discard)
			writeErr := errors.New("write failed")
			if failed {
				dst = &errorWriter{failure: writeErr}
			}

			err := codec.Encode(ctx, dst, &iterOnly{items: []runtime.Value{runtime.Int(1)}})
			if calls != 1 || (failed && !errors.Is(err, writeErr)) || (!failed && err != nil) {
				t.Fatalf("first operation: calls=%d err=%v", calls, err)
			}

			data, err := encoding.EncodeBytes(t.Context(), codec, &iterOnly{items: []runtime.Value{runtime.Int(2)}})
			if err != nil || calls != 1 {
				t.Fatalf("codec reused an earlier adopter: calls=%d err=%v", calls, err)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), ferretmsgpack.Default, data)
			if err != nil || decoded.String() != "[2]" {
				t.Fatalf("second output: value=%v err=%v", decoded, err)
			}
		})
	}
}

func TestMsgpackCachesAdoptionAcrossHostPredicateContexts(t *testing.T) {
	calls, lookups := 0, 0
	ctx := ownershipContext{Context: encodingownership.WithValueAdopter(t.Context(), func(runtime.Value) { calls++ }), lookups: &lookups}
	value := &ownershipMap{Object: runtime.NewObjectWith(map[string]runtime.Value{
		"a": &iterOnly{items: []runtime.Value{runtime.Int(1)}},
		"b": &iterOnly{items: []runtime.Value{runtime.Int(1)}},
	})}

	data, err := encoding.EncodeBytes(ctx, ferretmsgpack.Default, value)
	if err != nil || calls != 2 || lookups != 1 {
		t.Fatalf("operation adoption: calls=%d lookups=%d err=%v", calls, lookups, err)
	}

	decoded, err := encoding.DecodeBytes(t.Context(), ferretmsgpack.Default, data)
	if err != nil || decoded.String() != `{"a":[1],"b":[1]}` {
		t.Fatalf("encoded host map: value=%v err=%v", decoded, err)
	}
}
