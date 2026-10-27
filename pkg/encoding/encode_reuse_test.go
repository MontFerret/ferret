package encoding_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCodecConfigurationCopies(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Codec{}, ferretmsgpack.Codec{}} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			var events []string
			pre := func(name string) encoding.PreEncoderHook {
				return func(context.Context, runtime.Value) error {
					events = append(events, name)

					return nil
				}
			}
			post := func(name string) encoding.PostEncoderHook {
				return func(context.Context, runtime.Value, error) error {
					events = append(events, name)

					return nil
				}
			}
			decoderPre := func(name string) encoding.PreDecoderHook {
				return func(context.Context) error {
					events = append(events, name)

					return nil
				}
			}
			base := codec.EncodeWith().PreHook(pre("encode:base")).PostHook(post("post:base")).Encoder().(encoding.Codec)
			base = base.DecodeWith().PreHook(decoderPre("decode:base")).Decoder().(encoding.Codec)
			builder := base.EncodeWith().PreHook(pre("encode:left")).PostHook(post("post:left"))
			left := builder.Encoder().(encoding.Codec)
			latest := builder.PreHook(pre("encode:later")).Encoder().(encoding.Codec)
			right := base.EncodeWith().PreHook(pre("encode:right")).Encoder().(encoding.Codec)
			decodeBuilder := base.DecodeWith().PreHook(decoderPre("decode:left"))
			decodeLeft := decodeBuilder.Decoder().(encoding.Codec)
			decodeLatest := decodeBuilder.PreHook(decoderPre("decode:later")).Decoder().(encoding.Codec)
			decodeRight := base.DecodeWith().PreHook(decoderPre("decode:right")).Decoder().(encoding.Codec)

			for _, test := range []struct {
				codec encoding.Codec
				name  string
				want  []string
			}{
				{codec, "zero", nil},
				{base, "base", []string{"encode:base", "post:base", "decode:base"}},
				{left, "encoder_snapshot", []string{"encode:base", "encode:left", "post:base", "post:left", "decode:base"}},
				{latest, "encoder_extended", []string{"encode:base", "encode:left", "encode:later", "post:base", "post:left", "decode:base"}},
				{right, "encoder_sibling", []string{"encode:base", "encode:right", "post:base", "decode:base"}},
				{decodeLeft, "decoder_snapshot", []string{"encode:base", "post:base", "decode:base", "decode:left"}},
				{decodeLatest, "decoder_extended", []string{"encode:base", "post:base", "decode:base", "decode:left", "decode:later"}},
				{decodeRight, "decoder_sibling", []string{"encode:base", "post:base", "decode:base", "decode:right"}},
			} {
				t.Run(test.name, func(t *testing.T) {
					events = nil
					data, err := encoding.EncodeBytes(t.Context(), test.codec, runtime.Int(7))
					if err != nil {
						t.Fatal(err)
					}

					value, err := encoding.DecodeBytes(t.Context(), test.codec, data)
					if err != nil || value != runtime.Int(7) || !reflect.DeepEqual(events, test.want) {
						t.Fatalf("configuration changed: value=%v events=%v want=%v err=%v", value, events, test.want, err)
					}
				})
			}
		})
	}
}

func TestCodecConcurrentAndReentrantEncoding(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			for _, configured := range []bool{false, true} {
				t.Run(fmt.Sprintf("configured_%t", configured), func(t *testing.T) {
					type invocation struct {
						ctx  context.Context
						pre  int
						post int
					}
					type contextKey struct{}

					enc := encoding.Encoder(codec)
					if configured {
						enc = codec.EncodeWith().PreHook(func(ctx context.Context, _ runtime.Value) error {
							call, _ := ctx.Value(contextKey{}).(*invocation)
							if call == nil || ctx != call.ctx {
								return errors.New("pre-hook received another invocation's context")
							}

							call.pre++

							return nil
						}).PostHook(func(ctx context.Context, _ runtime.Value, err error) error {
							call, _ := ctx.Value(contextKey{}).(*invocation)
							if call == nil || ctx != call.ctx || err != nil {
								return fmt.Errorf("post-hook received another invocation's context or error: %v", err)
							}

							call.post++

							return nil
						}).Encoder()
					}

					value := runtime.NewArrayWith(runtime.True, runtime.None, runtime.String("payload"))
					golden, err := encoding.EncodeBytes(t.Context(), codec, value)
					if err != nil {
						t.Fatal(err)
					}

					scalar := runtime.String("shared scalar")
					scalarGolden, err := encoding.EncodeBytes(t.Context(), codec, scalar)
					if err != nil {
						t.Fatal(err)
					}

					results := make(chan error, 8)
					for range 8 {
						go func() {
							for i := range 32 {
								current, want, visits := runtime.Value(value), golden, 4
								if i%2 == 0 {
									current, want, visits = scalar, scalarGolden, 1
								}

								call := &invocation{}
								call.ctx = context.WithValue(t.Context(), contextKey{}, call)
								data, err := encoding.EncodeBytes(call.ctx, enc, current)
								if err != nil || !bytes.Equal(data, want) || (configured && (call.pre != visits || call.post != visits)) {
									results <- fmt.Errorf("concurrent encode: data=%x pre=%d post=%d err=%v", data, call.pre, call.post, err)

									return
								}
							}

							results <- nil
						}()
					}

					for range 8 {
						if err := <-results; err != nil {
							t.Error(err)
						}
					}

					outer, inner := &invocation{}, &invocation{}
					outer.ctx = context.WithValue(t.Context(), contextKey{}, outer)
					inner.ctx = context.WithValue(t.Context(), contextKey{}, inner)
					var nested []byte
					writer := &encodeReentrantWriter{beforeWrite: func() error {
						var err error
						nested, err = encoding.EncodeBytes(inner.ctx, enc, value)

						return err
					}}
					if err := enc.Encode(outer.ctx, writer, value); err != nil || !bytes.Equal(writer.Bytes(), golden) || !bytes.Equal(nested, golden) {
						t.Fatalf("reentrant encode corrupted output: outer=%x inner=%x err=%v", writer.Bytes(), nested, err)
					}

					if configured && (outer.pre != 4 || outer.post != 4 || inner.pre != 4 || inner.post != 4) {
						t.Fatalf("reentrant hooks crossed operations: outer=%d/%d inner=%d/%d", outer.pre, outer.post, inner.pre, inner.post)
					}
				})
			}
		})
	}
}

func TestCodecRetainsBytesAcrossReuseAndFailures(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			failure := errors.New("nested hook failed")
			value := runtime.NewArrayWith(runtime.String("retained"), runtime.Int(2))
			data, err := encoding.EncodeBytes(t.Context(), codec, value)
			if err != nil {
				t.Fatal(err)
			}

			golden := bytes.Clone(data)
			failing := codec.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
				if value == runtime.Int(2) {
					return failure
				}

				return nil
			}).Encoder()
			prefix, err := encoding.EncodeBytes(t.Context(), failing, value)
			if len(prefix) == 0 || !errors.Is(err, failure) {
				t.Fatalf("missing failed prefix: data=%x err=%v", prefix, err)
			}

			prefixGolden := bytes.Clone(prefix)
			writer := &failingWriter{limit: 0, failure: failure}
			if err := codec.Encode(t.Context(), writer, runtime.True); !errors.Is(err, failure) || writer.calls != 1 || writer.Len() != 0 {
				t.Fatalf("embedded buffer bypassed Write override: calls=%d data=%x err=%v", writer.calls, writer.Bytes(), err)
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			canceled := codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error {
				cancel()

				return nil
			}).Encoder()
			if err := canceled.Encode(ctx, io.Discard, value); !errors.Is(err, context.Canceled) {
				t.Fatalf("missing in-flight cancellation: %v", err)
			}

			for _, size := range []int{1, 4096, 64, 16384, 0} {
				if _, err := encoding.EncodeBytes(t.Context(), codec, runtime.String(strings.Repeat("x", size))); err != nil {
					t.Fatal(err)
				}
			}

			next, err := encoding.EncodeBytes(t.Context(), codec, value)
			if err != nil || !bytes.Equal(next, golden) || !bytes.Equal(data, golden) || !bytes.Equal(prefix, prefixGolden) || writer.calls != 1 {
				t.Fatalf("pooled state or returned bytes reused: data=%x prefix=%x next=%x calls=%d err=%v", data, prefix, next, writer.calls, err)
			}
		})
	}
}
