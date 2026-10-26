package json_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONCodecConfigurationSnapshots(t *testing.T) {
	var calls []string
	hook := func(name string) ferretencoding.PreEncoderHook {
		return func(context.Context, runtime.Value) error {
			calls = append(calls, name)

			return nil
		}
	}

	var zero json.Codec
	builder := zero.EncodeWith().PreHook(hook("base"))
	base := builder.Encoder().(json.Codec)
	copyOfBase := base
	builder.PreHook(hook("later"))
	later := builder.Encoder()
	derived := copyOfBase.EncodeWith().PreHook(hook("derived")).Encoder().(json.Codec)
	combined := derived.DecodeWith().PreHook(func(context.Context) error {
		calls = append(calls, "decode")

		return nil
	}).Decoder().(json.Codec)
	noOp := combined.EncodeWith().PreHook(nil).PostHook(nil).Encoder().(json.Codec)
	noOp = noOp.DecodeWith().PreHook(nil).PostHook(nil).Decoder().(json.Codec)

	for _, tc := range []struct {
		encoder ferretencoding.Encoder
		name    string
		want    []string
	}{
		{name: "zero", encoder: zero},
		{name: "published", encoder: base, want: []string{"base"}},
		{name: "copy_pointer", encoder: &copyOfBase, want: []string{"base"}},
		{name: "later", encoder: later, want: []string{"base", "later"}},
		{name: "derived", encoder: derived, want: []string{"base", "derived"}},
		{name: "combined", encoder: combined, want: []string{"base", "derived"}},
		{name: "no_op", encoder: noOp, want: []string{"base", "derived"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls = nil
			data, err := ferretencoding.EncodeBytes(t.Context(), tc.encoder, runtime.Int(1))
			if err != nil || string(data) != "1" || !reflect.DeepEqual(calls, tc.want) {
				t.Fatalf("encode = %s, %v, hooks %v; want hooks %v", data, err, calls, tc.want)
			}
		})
	}

	for _, codec := range []json.Codec{combined, noOp} {
		calls = nil
		value, err := codec.Decode(t.Context(), strings.NewReader("1"))
		if err != nil || value != runtime.Int(1) || !reflect.DeepEqual(calls, []string{"decode"}) {
			t.Fatalf("decode = %v, %v, hooks %v", value, err, calls)
		}
	}
}

func TestJSONCodecOperationRetainsConfiguration(t *testing.T) {
	var codec json.Codec
	pre, post := 0, 0
	codec = codec.EncodeWith().PreHook(func(context.Context, runtime.Value) error {
		codec = json.Codec{}
		pre++

		return nil
	}).PostHook(func(context.Context, runtime.Value, error) error {
		post++

		return nil
	}).Encoder().(json.Codec)

	data, err := ferretencoding.EncodeBytes(t.Context(), &codec, runtime.NewArrayWith(runtime.Int(1), runtime.Int(2)))
	if err != nil || string(data) != "[1,2]" || pre != 3 || post != 3 {
		t.Fatalf("encode = %s, %v, hooks %d/%d; want hooks 3/3", data, err, pre, post)
	}

	pre, post = 0, 0
	codec = codec.DecodeWith().PreHook(func(context.Context) error {
		codec = json.Codec{}
		pre++

		return nil
	}).PostHook(func(context.Context, runtime.Value, error) error {
		post++

		return nil
	}).Decoder().(json.Codec)

	value, err := codec.Decode(t.Context(), strings.NewReader("1"))
	if err != nil || value != runtime.Int(1) || pre != 1 || post != 1 {
		t.Fatalf("decode = %v, %v, hooks %d/%d; want hooks 1/1", value, err, pre, post)
	}
}
