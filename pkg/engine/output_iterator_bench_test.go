package engine

import (
	"context"
	"fmt"
	"testing"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkSessionCodecIterable(b *testing.B) {
	for _, contentType := range []string{encodingjson.ContentType, encodingmsgpack.ContentType} {
		b.Run(contentType, func(b *testing.B) {
			for _, shape := range []string{"scalars", "resources", "containers/depth_1", "containers/depth_8", "containers/depth_32"} {
				b.Run(shape, func(b *testing.B) {
					items := make([]runtime.Value, 1024)
					leaves := make([]*outputIteratorLeaf, 0, len(items))
					for i := range items {
						if shape == "scalars" {
							items[i] = runtime.NewInt(i)

							continue
						}

						leaf := &outputIteratorLeaf{text: "value"}
						leaves = append(leaves, leaf)
						items[i] = leaf
						var depth int
						if _, err := fmt.Sscanf(shape, "containers/depth_%d", &depth); err == nil {
							for range depth {
								items[i] = runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{"child": items[i]}))
							}
						}
					}

					iterator := &outputIterator{items: items}
					eng, err := New(WithFunctionsRegistrar(func(ns runtime.Namespace) {
						ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
							return &outputIterable{iter: iterator}, nil
						})
					}))
					if err != nil {
						b.Fatal(err)
					}

					defer eng.Close()
					ctx := context.Background()
					plan, err := eng.Compile(ctx, source.NewAnonymous("RETURN MAKE()"))
					if err != nil {
						b.Fatal(err)
					}

					defer plan.Close()
					session, err := plan.NewSession(ctx, WithOutputContentType(contentType))
					if err != nil {
						b.Fatal(err)
					}

					defer session.Close()
					b.ReportAllocs()
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						iterator.index = 0
						iterator.closes = 0
						for _, leaf := range leaves {
							leaf.closes = 0
						}

						handle, err := session.Run(ctx)
						if err != nil {
							b.Fatal(err)
						}

						content, err := handle.Collect(ctx)
						if err != nil {
							b.Fatal(err)
						}

						benchmarkOutputLength = len(content.Data)
						if err := handle.Close(); err != nil {
							b.Fatal(err)
						}
					}
				})
			}
		})
	}
}
