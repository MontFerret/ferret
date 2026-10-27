package encoding_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCodecScalarBufferMatchesWriter(t *testing.T) {
	regexp, err := runtime.NewRegexp(`^<item>&[a-z]+>$`)
	if err != nil {
		t.Fatal(err)
	}

	for _, family := range []struct {
		codec   encoding.Codec
		zero    encoding.Encoder
		pointer encoding.Encoder
	}{
		{ferretjson.Default, ferretjson.Codec{}, &ferretjson.Default},
		{ferretmsgpack.Default, ferretmsgpack.Codec{}, &ferretmsgpack.Default},
	} {
		t.Run(family.codec.ContentType(), func(t *testing.T) {
			decoderOnly := family.codec.DecodeWith().PreHook(func(context.Context) error {
				t.Error("encoding called a decoder hook")

				return nil
			}).Decoder().(encoding.Codec)
			pre, post := 0, 0
			var expectedContext context.Context
			hooked := family.codec.EncodeWith().PreHook(func(ctx context.Context, _ runtime.Value) error {
				if ctx != expectedContext {
					t.Error("pre-hook lost the caller's context")
				}

				pre++

				return nil
			}).PostHook(func(ctx context.Context, _ runtime.Value, err error) error {
				if ctx != expectedContext || err != nil {
					t.Errorf("post-hook lost context or received an error: %v", err)
				}

				post++

				return nil
			}).Encoder()

			for _, mode := range []struct {
				encoder encoding.Encoder
				name    string
				adopt   bool
				hooks   bool
			}{
				{encoder: family.codec, name: "default"},
				{encoder: family.zero, name: "zero"},
				{encoder: family.pointer, name: "pointer"},
				{encoder: decoderOnly, name: "decoder_only"},
				{encoder: hooked, name: "hooked", hooks: true},
				{encoder: family.codec, name: "background_adopter", adopt: true},
			} {
				t.Run(mode.name, func(t *testing.T) {
					var retained, snapshots [][]byte
					for _, test := range []struct {
						value runtime.Value
						name  string
					}{
						{name: "duration", value: runtime.Duration(1500 * time.Millisecond)},
						{name: "duration_unicode", value: runtime.Duration(1500 * time.Nanosecond)},
						{name: "regexp", value: regexp},
						{name: "empty", value: runtime.String("")},
						{name: "string_31", value: runtime.String(strings.Repeat("x", 31))},
						{name: "string_32", value: runtime.String(strings.Repeat("x", 32))},
						{name: "escaped", value: runtime.String("\"\\\n\t\x00\x1f")},
						{name: "unicode", value: runtime.String("<>&\u2028\u2029λ")},
						{name: "invalid_utf8", value: runtime.String("a\xff\xc0b")},
					} {
						t.Run(test.name, func(t *testing.T) {
							writer := &failingWriter{limit: 1 << 20}
							if err := family.codec.Encode(t.Context(), writer, test.value); err != nil || writer.calls == 0 {
								t.Fatalf("generic writer failed: calls=%d err=%v", writer.calls, err)
							}

							golden := bytes.Clone(writer.Bytes())
							for _, spare := range []int{0, 1, 31, 32, 33, 64} {
								const consumed, prefix = "consumed|", "prefix:"
								seed := make([]byte, len(consumed)+len(prefix), len(consumed)+len(prefix)+spare)
								copy(seed, consumed+prefix)
								buffer := bytes.NewBuffer(seed)
								buffer.Next(len(consumed))
								ctx := t.Context()
								var discoveries []bool
								lookups := 0
								if mode.adopt {
									ctx = encodingownership.WithValueAdopter(context.Background(), func(value runtime.Value, stored bool) {
										if value != test.value {
											t.Errorf("adopted another value: %v", value)
										}

										discoveries = append(discoveries, stored)
									})
									ctx = discoveryContext{Context: ctx, lookups: &lookups}
								}

								pre, post = 0, 0
								expectedContext = ctx
								err := mode.encoder.Encode(ctx, buffer, test.value)
								want := append([]byte(prefix), golden...)
								if err != nil || !bytes.Equal(buffer.Bytes(), want) {
									t.Fatalf("consumed buffer with %d spare bytes: output=%x want=%x err=%v", spare, buffer.Bytes(), want, err)
								}

								if mode.hooks && (pre != 1 || post != 1) {
									t.Fatalf("scalar hooks skipped or repeated: pre=%d post=%d", pre, post)
								}

								if mode.adopt && (lookups != 1 || !reflect.DeepEqual(discoveries, []bool{true, false})) {
									t.Fatalf("scalar ownership changed: lookups=%d discoveries=%v", lookups, discoveries)
								}
							}

							expectedContext = t.Context()
							data, err := encoding.EncodeBytes(expectedContext, mode.encoder, test.value)
							if err != nil || !bytes.Equal(data, golden) {
								t.Fatalf("fresh collected buffer: output=%x want=%x err=%v", data, golden, err)
							}

							retained = append(retained, data)
							snapshots = append(snapshots, golden)
						})
					}

					for i, data := range retained {
						if !bytes.Equal(data, snapshots[i]) {
							t.Fatalf("later scalar encoding reused output %d", i)
						}
					}
				})
			}
		})
	}
}

func TestCodecScalarBufferValidationAndDiscovery(t *testing.T) {
	for _, codec := range []encoding.Codec{ferretjson.Default, ferretmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			const prefix = "untouched"
			value := runtime.Duration(1500 * time.Millisecond)
			buffer := bytes.NewBufferString(prefix)
			if err := codec.Encode(nil, buffer, value); !errors.Is(err, runtime.ErrInvalidArgument) || buffer.String() != prefix {
				t.Fatalf("nil context reached scalar buffer: data=%q err=%v", buffer.String(), err)
			}

			for _, duringDiscovery := range []bool{false, true} {
				ctx, cancel := context.WithCancel(t.Context())
				t.Cleanup(cancel)
				if !duringDiscovery {
					cancel()
				}

				callbacks := 0
				ctx = encodingownership.WithValueAdopter(ctx, func(item runtime.Value, stored bool) {
					callbacks++
					if item != value || !stored {
						t.Errorf("unexpected canceled-root discovery: value=%v stored=%t", item, stored)
					}

					cancel()
				})
				if err := codec.Encode(ctx, (*bytes.Buffer)(nil), value); !errors.Is(err, runtime.ErrInvalidArgument) || errors.Is(err, context.Canceled) || callbacks != 0 {
					t.Fatalf("typed-nil buffer validation lost precedence: callbacks=%d err=%v", callbacks, err)
				}

				if err := codec.Encode(ctx, buffer, value); !errors.Is(err, context.Canceled) || callbacks != 1 || buffer.String() != prefix {
					t.Fatalf("scalar discovery/cancellation ordering changed: during=%t callbacks=%d data=%q err=%v", duringDiscovery, callbacks, buffer.String(), err)
				}

				cancel()
			}

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			if err := codec.Encode(ctx, buffer, value); !errors.Is(err, context.Canceled) || buffer.String() != prefix {
				t.Fatalf("canceled scalar without adopter wrote output: data=%q err=%v", buffer.String(), err)
			}

			failure := errors.New("wrapper write failed")
			for _, cause := range []error{failure, nil} {
				writer := &failingWriter{limit: 0, failure: cause}
				want := cause
				if want == nil {
					want = io.ErrShortWrite
				}

				if err := codec.Encode(t.Context(), writer, value); !errors.Is(err, want) || writer.calls != 1 || writer.Len() != 0 || writer.closed {
					t.Fatalf("scalar bypassed buffer wrapper: calls=%d data=%x closed=%t err=%v", writer.calls, writer.Bytes(), writer.closed, err)
				}
			}
		})
	}
}
