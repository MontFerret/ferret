package json_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
	"time"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestJSONRootStringPaths(t *testing.T) {
	regexp, err := runtime.NewRegexp(`^<item>&[a-z]+>$`)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		value runtime.Value
		name  string
		want  string
	}{
		{runtime.String(""), "empty", `""`},
		{runtime.Duration(1500 * time.Millisecond), "duration", `"1.5s"`},
		{regexp, "regexp", `"^<item>&[a-z]+>$"`},
		{runtime.String("<>&\x7f"), "ascii", "\"<>&\x7f\""},
		{runtime.Duration(time.Microsecond), "unicode_duration", `"1µs"`},
		{runtime.Duration(999999), "unicode_duration_max", `"999.999µs"`},
		{runtime.Duration(-999999), "unicode_duration_min", `"-999.999µs"`},
		{runtime.Duration(math.MinInt64), "duration_min", `"-2562047h47m16.854775808s"`},
		{runtime.Duration(math.MaxInt64), "duration_max", `"2562047h47m16.854775807s"`},
		{runtime.String("λ\u2028\u2029"), "unicode", "\"λ\u2028\u2029\""},
		{runtime.String("\xffa\xc0"), "invalid_utf8", "\"\ufffda\ufffd\""},
		{runtime.String("\x00\n\"\\"), "escaping", `"\u0000\n\"\\"`},
	}
	for _, size := range []int{4088, 4089, 4090, 4096} {
		text := strings.Repeat("a", size)
		cases = append(cases, struct {
			value runtime.Value
			name  string
			want  string
		}{runtime.String(text), fmt.Sprintf("length_%d", size), `"` + text + `"`})
	}

	for _, tc := range cases {
		for _, mode := range []string{"buffer", "wrapped", "hooks", "adopter"} {
			for _, spare := range []int{0, 32, len(tc.want)} {
				t.Run(fmt.Sprintf("%s/%s/%d", tc.name, mode, spare), func(t *testing.T) {
					const consumed, prefix = "consumed:", "prefix:"
					data := make([]byte, len(consumed)+len(prefix), len(consumed)+len(prefix)+spare)
					copy(data, consumed+prefix)
					buffer := bytes.NewBuffer(data)
					buffer.Next(len(consumed))
					var dst io.Writer = buffer
					enc := ferretencoding.Encoder(json.Default)
					ctx := t.Context()
					pre, post, adopted := 0, 0, 0

					switch mode {
					case "wrapped":
						dst = struct{ io.Writer }{buffer}
					case "hooks":
						enc = json.Default.EncodeWith().PreHook(func(actual context.Context, value runtime.Value) error {
							if actual != ctx || value != tc.value {
								t.Fatal("pre-hook lost its context or value")
							}

							pre++

							return nil
						}).PostHook(func(actual context.Context, value runtime.Value, err error) error {
							if actual != ctx || value != tc.value || err != nil {
								t.Fatalf("post-hook lost its context, value, or error: %v", err)
							}

							post++

							return nil
						}).Encoder()
					case "adopter":
						ctx = encodingownership.WithValueAdopter(ctx, func(value runtime.Value, stored bool) {
							if value != tc.value || stored != (adopted == 0) {
								t.Fatal("root or visitor adoption changed")
							}

							adopted++
						})
					}

					if err := enc.Encode(ctx, dst, tc.value); err != nil || buffer.String() != prefix+tc.want {
						t.Fatalf("encode = %q, %v; want %q", buffer.String(), err, prefix+tc.want)
					}

					if mode == "hooks" && (pre != 1 || post != 1) {
						t.Fatalf("hooks = %d/%d; want 1/1", pre, post)
					}

					if mode == "adopter" && adopted != 2 {
						t.Fatalf("adoption calls = %d; want 2", adopted)
					}
				})
			}
		}
	}
}

func TestJSONStringChunkBoundariesAndFailures(t *testing.T) {
	prefix, suffix := strings.Repeat("a", 4088), strings.Repeat("b", 4096)
	for _, tc := range []struct {
		name      string
		text      string
		want      string
		firstSize int
	}{
		{"ascii", prefix + "aa" + suffix, `"` + prefix + "aa" + suffix + `"`, 4090},
		{"escape", prefix + "\x00" + suffix, `"` + prefix + `\u0000` + suffix + `"`, 4095},
		{"unicode", prefix + "λ" + suffix, `"` + prefix + "λ" + suffix + `"`, 4091},
		{"invalid_utf8", prefix + "\xff" + suffix, `"` + prefix + "\ufffd" + suffix + `"`, 4092},
	} {
		t.Run(tc.name, func(t *testing.T) {
			writer := &stringChunkWriter{}
			if err := json.Default.Encode(t.Context(), writer, runtime.String(tc.text)); err != nil || writer.String() != tc.want {
				t.Fatalf("encode = %q, %v; want %q", writer.String(), err, tc.want)
			}

			if len(writer.sizes) < 2 || writer.sizes[0] != tc.firstSize {
				t.Fatalf("writes = %v; want first size %d", writer.sizes, tc.firstSize)
			}

			for _, size := range writer.sizes {
				if size > 4096 {
					t.Fatalf("unbounded write: %d", size)
				}
			}

			failure := errors.New("write failed")
			for _, writeErr := range []error{failure, nil} {
				writer = &stringChunkWriter{failAt: 2, failure: writeErr}
				wantErr := writeErr
				if wantErr == nil {
					wantErr = io.ErrShortWrite
				}

				if err := json.Default.Encode(t.Context(), writer, runtime.String(tc.text)); !errors.Is(err, wantErr) || len(writer.sizes) != 2 || writer.String() != tc.want[:tc.firstSize] {
					t.Fatalf("failed encode: writes=%v bytes=%d err=%v", writer.sizes, writer.Len(), err)
				}
			}

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			writer = &stringChunkWriter{cancel: cancel}
			if err := json.Default.Encode(ctx, writer, runtime.String(tc.text)); !errors.Is(err, context.Canceled) || len(writer.sizes) != 1 || writer.String() != tc.want[:tc.firstSize] {
				t.Fatalf("canceled encode: writes=%v bytes=%d err=%v", writer.sizes, writer.Len(), err)
			}
		})
	}
}
