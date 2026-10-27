package msgpack_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	ferretmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMsgpackRootStringHeaderBoundaries(t *testing.T) {
	for _, text := range []string{"", "a", strings.Repeat("a", 31), strings.Repeat("a", 32), strings.Repeat("é", 15), strings.Repeat("é", 16), "\xff"} {
		t.Run(fmt.Sprintf("string_%x", text), func(t *testing.T) {
			want := mustMarshalNative(t, text)
			for _, buffer := range []bool{false, true} {
				var direct bytes.Buffer
				indirect := &opaqueBufferWriter{}
				var dst io.Writer = indirect
				if buffer {
					dst = &direct
				}

				if err := ferretmsgpack.Default.Encode(t.Context(), dst, runtime.String(text)); err != nil {
					t.Fatal(err)
				}

				got := indirect.Bytes()
				if buffer {
					got = direct.Bytes()
				}

				if !bytes.Equal(got, want) {
					t.Fatalf("buffer=%t encoded %x; want %x", buffer, got, want)
				}
			}
		})
	}

	for _, size := range []int{0, 31, 32} {
		value, err := runtime.NewRegexp(strings.Repeat("a", size))
		if err != nil {
			t.Fatal(err)
		}

		data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
		if err != nil || !bytes.Equal(data, mustMarshalNative(t, value.String())) {
			t.Fatalf("regexp length=%d encoded %x: %v", size, data, err)
		}
	}

	for _, value := range []runtime.Duration{0, runtime.Duration(-time.Hour), runtime.Duration(time.Duration(1<<63 - 1))} {
		data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
		if err != nil || !bytes.Equal(data, mustMarshalNative(t, value.String())) {
			t.Fatalf("duration %s encoded %x: %v", value, data, err)
		}
	}
}

func TestMsgpackBufferWriterPreservesCustomStringSettings(t *testing.T) {
	for _, dictionary := range []bool{false, true} {
		t.Run(fmt.Sprintf("dictionary_%t", dictionary), func(t *testing.T) {
			settings := &stringSettingsValue{dictionary: dictionary}
			value := runtime.NewArrayWith(settings, runtime.String("repeat"), runtime.String("repeat"), runtime.NewObjectWith(map[string]runtime.Value{"repeat": runtime.String("repeat")}))
			want := mustEncodeRaw(t, func(enc *vmmsgpack.Encoder) error {
				if err := enc.EncodeArrayLen(4); err != nil {
					return err
				}

				if err := settings.EncodeMsgpack(enc); err != nil {
					return err
				}

				for range 2 {
					if err := enc.EncodeString("repeat"); err != nil {
						return err
					}
				}

				if err := enc.EncodeMapLen(1); err != nil {
					return err
				}

				if err := enc.EncodeString("repeat"); err != nil {
					return err
				}

				return enc.EncodeString("repeat")
			})
			for _, configured := range []bool{false, true} {
				enc := encoding.Encoder(ferretmsgpack.Default)
				posts := 0
				if configured {
					enc = ferretmsgpack.Default.EncodeWith().PostHook(func(context.Context, runtime.Value, error) error {
						posts++

						return nil
					}).Encoder()
				}

				data, err := encoding.EncodeBytes(t.Context(), enc, value)
				if err != nil || !bytes.Equal(data, want) || (configured && posts != 6) {
					t.Fatalf("configured=%t custom settings changed: data=%x want=%x posts=%d err=%v", configured, data, want, posts, err)
				}
			}

			data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, runtime.NewArrayWith(runtime.String("repeat"), runtime.String("repeat")))
			if err != nil || !bytes.Equal(data, mustMarshalNative(t, []string{"repeat", "repeat"})) {
				t.Fatalf("settings survived pool reuse: %x, %v", data, err)
			}
		})
	}
}

func TestMsgpackBufferWriterReportsIgnoredCancellationToHooks(t *testing.T) {
	customErr := errors.New("custom encoding failed")
	for _, returned := range []error{nil, customErr} {
		ctx, cancel := context.WithCancel(t.Context())
		var postErr error
		posts := 0
		enc := ferretmsgpack.Default.EncodeWith().PostHook(func(actual context.Context, _ runtime.Value, err error) error {
			if actual != ctx {
				t.Error("hook lost caller context")
			}

			posts++
			postErr = err

			return nil
		}).Encoder()
		var dst bytes.Buffer
		value := runtime.NewArrayWith(&cancelingWriterValue{cancel: cancel, failure: returned})
		err := enc.Encode(ctx, &dst, value)
		cancel()
		// The last post-hook belongs to the ancestor array; its header is the
		// only output committed before the custom child ignores cancellation.
		if !errors.Is(err, context.Canceled) || !errors.Is(postErr, context.Canceled) || (returned != nil && (!errors.Is(err, returned) || !errors.Is(postErr, returned))) || posts != 2 || !bytes.Equal(dst.Bytes(), []byte{0x91}) {
			t.Fatalf("ignored cancellation lost: result=%v post=%v posts=%d length=%d", err, postErr, posts, dst.Len())
		}
	}
}

func TestMsgpackMapUsesHostKeyString(t *testing.T) {
	key := &hostStringKey{Value: runtime.String("embedded value")}
	value := &hostKeyMap{
		Object: runtime.NewObjectWith(map[string]runtime.Value{"stored key": runtime.Int(7)}),
		key:    key,
	}
	data, err := encoding.EncodeBytes(t.Context(), ferretmsgpack.Default, value)
	want := mustMarshalNative(t, map[string]int{"host key": 7})
	if err != nil || key.calls != 1 || !bytes.Equal(data, want) {
		t.Fatalf("host key string changed: data=%x want=%x calls=%d err=%v", data, want, key.calls, err)
	}
}
