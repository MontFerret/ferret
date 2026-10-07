package uapi

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkOutputExecution(b *testing.B) {
	for _, size := range []int{32, 64 << 10, 1 << 20} {
		for _, reusable := range []bool{false, true} {
			for _, adapter := range []string{"Native", "Universal"} {
				for _, mode := range []string{"Collect", "Discard", "Writer"} {
					b.Run(fmt.Sprintf("%d/%t/%s/%s", size, reusable, adapter, mode), func(b *testing.B) {
						native := newTestEngine(b)
						portable := Wrap(native, "benchmark")
						payload := strings.Repeat("x", size)
						query := source.NewAnonymous("RETURN @value")
						portableQuery := api.NewAnonymousSource("RETURN @value")
						run := func() (api.Output, error) {
							return native.Run(b.Context(), query, engine.WithSessionParam("value", payload))
						}
						if adapter == "Universal" {
							run = func() (api.Output, error) {
								return portable.Run(b.Context(), portableQuery, api.WithParam("value", payload))
							}
						}

						if reusable {
							plan, err := native.Compile(b.Context(), query)
							if err != nil {
								b.Fatal(err)
							}

							b.Cleanup(func() { _ = plan.Close() })
							nativeSession, err := plan.NewSession(b.Context(), engine.WithSessionParam("value", payload))
							if err != nil {
								b.Fatal(err)
							}

							b.Cleanup(func() { _ = nativeSession.Close() })
							run = func() (api.Output, error) { return nativeSession.Run(b.Context()) }
							if adapter == "Universal" {
								portableSession := &session{native: nativeSession}
								run = func() (api.Output, error) { return portableSession.Run(b.Context()) }
							}
						}

						var writer bytes.Buffer
						b.ReportAllocs()
						b.SetBytes(int64(size + 2))
						consumer := func(_ context.Context, data []byte) error {
							if len(data) != size+2 {
								b.Fatal("unexpected output")
							}

							if mode == "Writer" {
								_, err := writer.Write(data)

								return err
							}

							return nil
						}
						b.ResetTimer()
						for b.Loop() {
							out, err := run()
							if err != nil {
								b.Fatal(err)
							}

							if mode == "Collect" {
								content, err := out.Collect(b.Context())
								if err != nil || content == nil || len(content.Data) != size+2 {
									b.Fatalf("content=%v err=%v", content, err)
								}
							} else {
								writer.Reset()
								err := out.Consume(b.Context(), consumer)
								if err != nil {
									b.Fatal(err)
								}
							}

							if err := out.Close(); err != nil {
								b.Fatal(err)
							}
						}
					})
				}
			}
		}
	}
}
