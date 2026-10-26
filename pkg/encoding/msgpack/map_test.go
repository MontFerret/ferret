package msgpack_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMsgpackMapMixedIntegerValues(t *testing.T) {
	value := runtime.NewObjectWith(map[string]runtime.Value{
		"minimum": runtime.Int(math.MinInt64),
		"maximum": runtime.Int(math.MaxInt64),
		"custom":  &customIntegerValue{Int: 7},
		"float":   runtime.Float(1.25),
		"string":  runtime.String("value"),
		"null":    runtime.None,
	})
	data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := encoding.DecodeBytes(t.Context(), ferretmsgpack.Default, data)
	if err != nil {
		t.Fatal(err)
	}

	want := runtime.NewObjectWith(map[string]runtime.Value{
		"minimum": runtime.Int(math.MinInt64),
		"maximum": runtime.Int(math.MaxInt64),
		"custom":  runtime.String("custom integer"),
		"float":   runtime.Float(1.25),
		"string":  runtime.String("value"),
		"null":    runtime.None,
	})
	assertValueEqual(t, decoded, want)
}

func TestMsgpackMapIntegerRetainsHooksAndAdoption(t *testing.T) {
	const integer = runtime.Int(-9)
	value := runtime.NewObjectWith(map[string]runtime.Value{"value": integer})
	for _, withHooks := range []bool{false, true} {
		for _, withAdopter := range []bool{false, true} {
			t.Run(fmt.Sprintf("hooks_%t_adopter_%t", withHooks, withAdopter), func(t *testing.T) {
				pre, post, visits := 0, 0, 0
				ctx := t.Context()
				if withAdopter {
					ctx = encodingownership.WithValueAdopter(ctx, func(value runtime.Value, _ bool) {
						if value == integer {
							visits++
						}
					})
				}

				enc := encoding.Encoder(ferretmsgpack.Default)
				if withHooks {
					enc = ferretmsgpack.Default.EncodeWith().PreHook(func(actual context.Context, value runtime.Value) error {
						if actual != ctx {
							t.Error("pre-hook lost caller context")
						}

						if value == integer {
							pre++
						}

						return nil
					}).PostHook(func(actual context.Context, value runtime.Value, err error) error {
						if actual != ctx || err != nil {
							t.Errorf("post-hook lost caller context or failed: %v", err)
						}

						if value == integer {
							post++
						}

						return nil
					}).Encoder()
				}

				data, err := encoding.EncodeBytes(ctx, enc, value)
				if err != nil {
					t.Fatal(err)
				}

				decoded, err := encoding.DecodeBytes(t.Context(), ferretmsgpack.Default, data)
				if err != nil {
					t.Fatal(err)
				}

				assertValueEqual(t, decoded, value)
				if withHooks && (pre != 1 || post != 1) {
					t.Fatalf("integer hooks skipped: pre=%d post=%d", pre, post)
				}

				if withAdopter && ((!withHooks && visits != 2) || (withHooks && visits < 2)) {
					t.Fatalf("integer acquisition or recursive adoption skipped: visits=%d", visits)
				}
			})
		}
	}
}

func TestMsgpackMapIntegerSamplesTraversalCancellation(t *testing.T) {
	properties := make(map[string]runtime.Value, 1000)
	for i := 0; i < 1000; i++ {
		properties[strconv.Itoa(i)] = runtime.Int(1)
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	writer := &mapCancelWriter{cancel: cancel}
	err := ferretmsgpack.Default.Encode(ctx, writer, runtime.NewObjectWith(properties))
	// Cancellation follows the map header. Each entry occupies at most five
	// bytes, so the traversal checkpoint must stop before the 4 KiB I/O limit.
	if !errors.Is(err, context.Canceled) || writer.Len() == 0 || writer.Len() > 1280 {
		t.Fatalf("integer map missed traversal cancellation: length=%d error=%v", writer.Len(), err)
	}
}
