package engine

import (
	"context"
	"testing"

	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkSessionMsgpackIterable(b *testing.B) {
	for _, shape := range []string{"scalars", "resources", "containers"} {
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
				if shape == "containers" {
					items[i] = runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{"child": leaf}))
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
			session, err := plan.NewSession(ctx, WithOutputContentType(encodingmsgpack.ContentType))
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
}
