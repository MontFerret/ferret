package encoding_test

import (
	"context"
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCodecBorrowsSelfIterators(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			iterator := &selfCodecIterator{codecIterator: codecIterator{
				items:    []runtime.Value{runtime.Int(1)},
				closeErr: errors.New("caller-owned iterator was closed"),
			}}

			data, err := encoding.EncodeBytes(t.Context(), codec, iterator)
			if err != nil || iterator.closes != 0 {
				t.Fatalf("borrowed iterator: closes=%d err=%v", iterator.closes, err)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), codec, data)
			if err != nil || decoded.String() != "[1]" {
				t.Fatalf("output: value=%v err=%v", decoded, err)
			}

			iterator.index = 0
			writeErr := errors.New("write failed")
			writer := &failingWriter{failure: writeErr, limit: 1}
			if err := codec.Encode(t.Context(), writer, iterator); !errors.Is(err, writeErr) || iterator.closes != 0 {
				t.Fatalf("borrowed iterator on failure: closes=%d err=%v", iterator.closes, err)
			}
		})
	}
}

func TestCodecIteratorResourceIdentity(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			for _, id := range []uint64{1, 2} {
				iterator := &resourceCodecIterator{codecIterator: &codecIterator{items: []runtime.Value{runtime.Int(1)}}, id: id}
				value := &resourceCodecIterable{codecIterable: &codecIterable{iter: iterator}, id: 1}

				if err := codec.Encode(t.Context(), io.Discard, value); err != nil {
					t.Fatal(err)
				}

				wantCloses := 0
				if id != value.id {
					wantCloses = 1
				}

				if iterator.closes != wantCloses {
					t.Fatalf("resource IDs %d/%d: closes=%d, want %d", value.id, id, iterator.closes, wantCloses)
				}
			}
		})
	}
}

func TestCodecIteratorHooksAndClosureOrder(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			var events []string
			closeErr := errors.New("iterator close failed")
			iterator := &codecIterator{
				items:    []runtime.Value{runtime.NewArrayWith(runtime.Int(1), runtime.Int(2))},
				closeErr: closeErr,
				events:   &events,
			}
			value := &codecIterable{iter: iterator}
			label := func(value runtime.Value) string {
				if _, ok := value.(runtime.List); ok {
					return "list"
				}

				return value.String()
			}
			var rootErr error
			enc := codec.EncodeWith().PreHook(func(ctx context.Context, value runtime.Value) error {
				if ctx != t.Context() {
					t.Fatal("pre-hook lost context")
				}

				events = append(events, "pre "+label(value))

				return nil
			}).PostHook(func(ctx context.Context, item runtime.Value, err error) error {
				if ctx != t.Context() {
					t.Fatal("post-hook lost context")
				}

				if item == value {
					rootErr = err
				} else if err != nil {
					t.Fatalf("iterator closure reached child post-hook: %v", err)
				}

				events = append(events, "post "+label(item))

				return nil
			}).Encoder()

			data, err := encoding.EncodeBytes(t.Context(), enc, value)
			if !errors.Is(err, closeErr) || !errors.Is(rootErr, closeErr) || iterator.closes != 1 {
				t.Fatalf("iterator closure: closes=%d result=%v root=%v", iterator.closes, err, rootErr)
			}

			want := []string{"pre iterable", "pre list", "pre 1", "post 1", "pre 2", "post 2", "post list", "close iterator", "post iterable"}
			if !slices.Equal(events, want) {
				t.Fatalf("depth-first hooks and cleanup: %v, want %v", events, want)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), codec, data)
			if err != nil || decoded.String() != "[[1,2]]" {
				t.Fatalf("encoded children before closing: value=%v err=%v", decoded, err)
			}
		})
	}
}

func TestCodecIteratorSettlesFailureAndCancellation(t *testing.T) {
	primary := errors.New("operation failed")
	cleanup := errors.New("iterator close failed")
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			for _, stage := range []string{"next", "writer", "pre-hook", "cancel"} {
				t.Run(stage, func(t *testing.T) {
					iterator := &codecIterator{items: []runtime.Value{runtime.Int(1)}, closeErr: cleanup}
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()

					dst := io.Writer(io.Discard)
					var postErr error
					value := &codecIterable{iter: iterator}
					config := codec.EncodeWith().PostHook(func(_ context.Context, item runtime.Value, err error) error {
						if item == value {
							postErr = err
						}

						return nil
					})

					wantErr := primary
					switch stage {
					case "next":
						iterator.nextErr = primary
					case "writer":
						dst = &failingWriter{failure: primary, limit: 1}
					case "pre-hook":
						config.PreHook(func(_ context.Context, item runtime.Value) error {
							if item == runtime.Int(1) {
								return primary
							}

							return nil
						})
					case "cancel":
						iterator.items = make([]runtime.Value, 1024)
						for i := range iterator.items {
							iterator.items[i] = runtime.Int(1)
						}

						iterator.onNext = cancel
						wantErr = context.Canceled
					}

					err := config.Encoder().Encode(ctx, dst, value)
					if !errors.Is(err, wantErr) || !errors.Is(err, cleanup) || !errors.Is(postErr, wantErr) || !errors.Is(postErr, cleanup) || iterator.closes != 1 {
						t.Fatalf("settlement: closes=%d result=%v post=%v", iterator.closes, err, postErr)
					}

					if stage == "cancel" && iterator.index > 256 {
						t.Fatalf("cancellation was not bounded: visited=%d", iterator.index)
					}
				})
			}
		})
	}
}
