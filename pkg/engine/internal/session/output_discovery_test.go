package session

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/internal/encodingownership"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestMaterializeCanceledRootRetainsDiscoveredResources(t *testing.T) {
	program, err := spec.Compile("RETURN @root")
	if err != nil {
		t.Fatal(err)
	}

	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			instance, err := vm.New(program)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = instance.Close() })
			root := runtime.NewArray(1)
			env := vm.NewDefaultEnvironment()
			env.Params["root"] = root

			result, err := instance.Run(t.Context(), env)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = result.Close() })
			if result.Root() != root {
				t.Fatal("execution did not preserve the stored root")
			}

			closeErr := errors.New("descendant cleanup failed")
			first := &abandonmentResource{Value: runtime.String("first")}
			failed := &abandonmentResource{Value: runtime.String("failed"), closeErr: closeErr}
			hostChild := &abandonmentResource{Value: runtime.String("host child")}
			lazyChild := &abandonmentResource{Value: runtime.String("lazy child")}
			host := &abandonmentMap{
				Object: runtime.NewObjectWith(map[string]runtime.Value{"hidden": hostChild}),
				value:  hostChild,
				key:    runtime.String("hidden"),
			}
			lazy := &abandonmentIterable{Value: runtime.String("lazy"), item: lazyChild}
			nested := runtime.NewObjectWith(map[string]runtime.Value{
				"resources": runtime.NewArrayWith(first, failed, runtime.NewArrayWith(first)),
				"host":      host,
				"lazy":      lazy,
			})

			// Add descendants after execution so only materialization can adopt them.
			if err := root.Append(t.Context(), nested); err != nil {
				t.Fatal(err)
			}

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			content, err := Materialize(ctx, encoding.NewRegistry(codec), codec.ContentType(), result)
			if !errors.Is(err, context.Canceled) || content != nil {
				t.Fatalf("canceled materialization: content=%v err=%v", content, err)
			}

			if first.closes != 0 || failed.closes != 0 || hostChild.closes != 0 || lazyChild.closes != 0 {
				t.Fatal("materialization prematurely closed descendants")
			}

			if host.forEachCalls != 0 || lazy.iterates != 0 || lazy.nexts != 0 {
				t.Fatalf("stored discovery evaluated collections: host=%d iterate=%d next=%d", host.forEachCalls, lazy.iterates, lazy.nexts)
			}

			if err := result.Close(); !errors.Is(err, closeErr) {
				t.Errorf("result cleanup lost descendant failure: %v", err)
			}

			if first.closes != 1 || failed.closes != 1 {
				t.Errorf("discovered cleanup: first=%d failed=%d", first.closes, failed.closes)
			}

			if err := result.Close(); err != nil {
				t.Errorf("repeated result cleanup: %v", err)
			}

			if err := instance.Close(); err != nil {
				t.Errorf("VM cleanup: %v", err)
			}

			if first.closes != 1 || failed.closes != 1 || hostChild.closes != 0 || lazyChild.closes != 0 {
				t.Errorf("owner cleanup: first=%d failed=%d host=%d lazy=%d", first.closes, failed.closes, hostChild.closes, lazyChild.closes)
			}
		})
	}
}

func TestCodecDiscoversStoredMutationBeforeKeyWriteFailure(t *testing.T) {
	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			first := &abandonmentResource{Value: runtime.String("first")}
			added := &abandonmentResource{Value: runtime.String("added")}
			array := runtime.NewArrayWith(first)
			key := &abandonmentKey{Value: runtime.String("key"), mutate: func() {
				if err := array.Append(t.Context(), added); err != nil {
					t.Fatal(err)
				}
			}}
			value := &abandonmentMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": array}), value: array, key: key}
			var result vm.Result
			adopter := outputValueAdopter{result: &result}
			failure := errors.New("key write failed")
			ctx := encodingownership.WithValueAdopter(t.Context(), adopter.Adopt)
			if err := codec.Encode(ctx, &abandonmentWriter{failure: failure}, value); !errors.Is(err, failure) {
				t.Fatalf("lost writer failure: %v", err)
			}

			if first.closes != 0 || added.closes != 0 {
				t.Fatal("codec closed children")
			}

			for range 2 {
				if err := result.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if first.closes != 1 || added.closes != 1 {
				t.Fatalf("mutation cleanup: first=%d added=%d", first.closes, added.closes)
			}
		})
	}
}
