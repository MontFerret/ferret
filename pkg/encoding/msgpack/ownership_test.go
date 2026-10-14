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
			calls, yielded := 0, false
			ctx := encodingownership.WithValueAdopter(t.Context(), func(value runtime.Value, _ bool) {
				calls++
				if value == runtime.Int(1) {
					yielded = true
				}
			})
			dst := io.Writer(io.Discard)
			writeErr := errors.New("write failed")
			if failed {
				dst = &errorWriter{failure: writeErr}
			}

			err := codec.Encode(ctx, dst, &iterOnly{items: []runtime.Value{runtime.Int(1)}})
			if !yielded || (failed && !errors.Is(err, writeErr)) || (!failed && err != nil) {
				t.Fatalf("first operation: calls=%d err=%v", calls, err)
			}

			firstCalls := calls
			data, err := encoding.EncodeBytes(t.Context(), codec, &iterOnly{items: []runtime.Value{runtime.Int(2)}})
			if err != nil || calls != firstCalls {
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
	ctx := ownershipContext{Context: encodingownership.WithValueAdopter(t.Context(), func(runtime.Value, bool) { calls++ }), lookups: &lookups}
	value := &ownershipMap{Object: runtime.NewObjectWith(map[string]runtime.Value{
		"a": &iterOnly{items: []runtime.Value{runtime.Int(1)}},
		"b": &iterOnly{items: []runtime.Value{runtime.Int(1)}},
	})}

	data, err := encoding.EncodeBytes(ctx, ferretmsgpack.Default, value)
	if err != nil || calls == 0 || lookups != 1 {
		t.Fatalf("operation adoption: calls=%d lookups=%d err=%v", calls, lookups, err)
	}

	decoded, err := encoding.DecodeBytes(t.Context(), ferretmsgpack.Default, data)
	if err != nil || decoded.String() != `{"a":[1],"b":[1]}` {
		t.Fatalf("encoded host map: value=%v err=%v", decoded, err)
	}
}
