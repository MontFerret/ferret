package encoding_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/encoding"
)

func TestQueryEncoding(t *testing.T) {
	for _, tc := range []struct{ text, encoded string }{
		{text: "", encoded: ""}, {text: "a +/?&é", encoded: "a+%2B%2F%3F%26%C3%A9"},
		{text: "https://github.com/MontFerret/ferret", encoded: "https%3A%2F%2Fgithub.com%2FMontFerret%2Fferret"},
	} {
		got, err := encoding.QueryEscape(context.Background(), runtime.String(tc.text))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(tc.encoded), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(tc.encoded))
		}

		got, err = encoding.QueryUnescape(context.Background(), got)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(runtime.String(tc.text), got) {
			t.Fatalf("got %#v, want %#v", got, runtime.String(tc.text))
		}
	}
}

func TestBase64PreservesBytes(t *testing.T) {
	for _, tc := range []struct {
		encoded string
		bytes   runtime.Binary
	}{
		{bytes: runtime.Binary{}, encoded: ""}, {bytes: runtime.Binary("foobar"), encoded: "Zm9vYmFy"}, {bytes: runtime.Binary{0, 255, 128}, encoded: "AP+A"}, {bytes: runtime.Binary("é"), encoded: "w6k="},
	} {
		for _, input := range []runtime.Value{tc.bytes, runtime.String(tc.bytes)} {
			got, err := encoding.Base64Encode(context.Background(), input)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(runtime.String(tc.encoded), got) {
				t.Fatalf("got %#v, want %#v", got, runtime.String(tc.encoded))
			}

			got, err = encoding.Base64Decode(context.Background(), got)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.bytes, got) {
				t.Fatalf("got %#v, want %#v", got, tc.bytes)
			}
		}
	}

	for _, bad := range []string{"!", "Zg", "====", "w6k==="} {
		_, err := encoding.Base64Decode(context.Background(), runtime.String(bad))
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}
}

func TestEncodingStrictInputs(t *testing.T) {
	for _, fn := range []func(context.Context, runtime.Value) (runtime.Value, error){
		encoding.JSONParse, encoding.QueryEscape, encoding.QueryUnescape, encoding.HTMLEscape, encoding.HTMLUnescape, encoding.Base64Encode, encoding.Base64Decode,
	} {
		for _, bad := range []runtime.Value{runtime.Int(1), runtime.True, runtime.None, runtime.NewArray(0)} {
			_, err := fn(context.Background(), bad)
			if !errors.Is(err, runtime.ErrInvalidArgument) {
				t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
			}
		}
	}
}
