package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSessionCustomCodecRetainsOwnershipHook(t *testing.T) {
	for _, base := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(base.ContentType(), func(t *testing.T) {
			leaf := &outputIteratorLeaf{text: "leaf"}
			root := runtime.NewArrayWith(leaf)
			var visits []runtime.Value
			configured := base.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
				if leaf.closes != 0 {
					t.Fatal("premature cleanup")
				}

				visits = append(visits, value)

				return nil
			}).Encoder().(encoding.Codec)
			codec := hookOwnershipCodec{Codec: configured}
			eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return root, nil })
			}))
			defer eng.Close()
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			defer plan.Close()
			session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
			defer session.Close()
			handle, err := session.Run(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			if _, err := handle.Collect(t.Context()); err != nil {
				t.Fatal(err)
			}

			if len(visits) != 2 || visits[0] != root || visits[1] != leaf || leaf.closes != 1 {
				t.Fatalf("hook order or result cleanup changed: visits=%v closes=%d", visits, leaf.closes)
			}

			for range 2 {
				if err := handle.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if leaf.closes != 1 {
				t.Fatalf("custom ownership hook lost: closes=%d", leaf.closes)
			}
		})
	}
}

func TestSessionBuiltinAdoptsStoredChildrenBeforeRootHook(t *testing.T) {
	for _, base := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		for _, pointer := range []bool{false, true} {
			t.Run(base.ContentType()+map[bool]string{false: "/value", true: "/pointer"}[pointer], func(t *testing.T) {
				first := &outputIteratorLeaf{text: "first"}
				added := &outputIteratorLeaf{text: "added"}
				root := runtime.NewArrayWith(first)
				failure := errors.New("root hook failed")
				codec := base.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
					if first.closes != 0 {
						t.Fatal("premature close")
					}

					if value == root {
						if err := root.Append(t.Context(), added); err != nil {
							t.Fatal(err)
						}

						return failure
					}

					return nil
				}).Encoder().(encoding.Codec)
				if pointer {
					switch c := codec.(type) {
					case encodingjson.Codec:
						codec = &c
					case encodingmsgpack.Codec:
						codec = &c
					}
				}

				eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return root, nil })
				}))
				plan := mustCompilePlan(t, eng, "RETURN MAKE()")
				session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
				handle, err := session.Run(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if _, err = handle.Collect(t.Context()); !errors.Is(err, failure) {
					t.Fatalf("lost hook failure: %v", err)
				}

				for range 2 {
					for _, close := range []func() error{handle.Close, session.Close, plan.Close, eng.Close} {
						if err := close(); err != nil {
							t.Fatal(err)
						}
					}
				}

				if first.closes != 1 || added.closes != 1 || first.marshals != 0 || added.marshals != 0 {
					t.Fatalf("abandoned resources: first=%+v added=%+v", first, added)
				}
			})
		}
	}
}

func TestSessionCodecAdoptsCanceledCyclicYield(t *testing.T) {
	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			first := &outputIteratorLeaf{text: "first"}
			second := &outputIteratorLeaf{text: "second"}
			cyclic := runtime.NewArrayWith(first)
			object := runtime.NewObjectWith(map[string]runtime.Value{"array": cyclic, "resource": second})
			if err := cyclic.Append(t.Context(), object); err != nil {
				t.Fatal(err)
			}

			items := make([]runtime.Value, 255)
			for i := range items {
				items[i] = runtime.Int(i)
			}

			items[254] = cyclic
			iter := &outputIterator{items: items, onNext: func(_ context.Context, index int) {
				if index == 254 {
					cancel()
				}
			}}
			eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return &outputIterable{iter: iter}, nil })
			}))
			defer eng.Close()
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			defer plan.Close()
			session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
			defer session.Close()
			handle, err := session.Run(ctx)
			if err != nil {
				t.Fatal(err)
			}

			_, err = handle.Collect(t.Context())
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("lost cancellation: %v", err)
			}

			for range 2 {
				if err := handle.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if first.closes != 1 || second.closes != 1 || iter.closes != 1 || first.marshals != 0 || second.marshals != 0 {
				t.Fatalf("cleanup: leaves=(%d,%d) iterator=%d marshals=(%d,%d)", first.closes, second.closes, iter.closes, first.marshals, second.marshals)
			}
		})
	}
}

func TestSessionCodecAdoptsHookMutationsAndDeepStoredResources(t *testing.T) {
	for _, base := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		for _, phase := range []string{"pre", "post"} {
			t.Run(base.ContentType()+"/"+phase, func(t *testing.T) {
				first := &outputIteratorLeaf{text: "first"}
				added := &outputIteratorLeaf{text: "added"}
				host := &outputIteratorHostList{Array: runtime.NewArrayWith(&outputIteratorLeaf{text: "unvisited"})}
				var nested runtime.Value = runtime.NewArrayWith(first, host)
				for range 32768 {
					nested = runtime.NewArrayWith(nested)
				}

				root := runtime.NewArrayWith(nested)
				failure := errors.New("hook failed")
				var codec encoding.Codec
				if phase == "pre" {
					codec = base.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
						if value == root {
							return failure
						}

						return nil
					}).Encoder().(encoding.Codec)
				} else {
					root = runtime.NewArrayWith(first)
					codec = base.EncodeWith().PostHook(func(_ context.Context, value runtime.Value, err error) error {
						if value == root {
							if appendErr := root.Append(t.Context(), added); appendErr != nil {
								t.Fatal(appendErr)
							}

							return failure
						}

						return nil
					}).Encoder().(encoding.Codec)
				}

				eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return root, nil })
				}))
				defer eng.Close()
				plan := mustCompilePlan(t, eng, "RETURN MAKE()")
				defer plan.Close()
				session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
				defer session.Close()
				handle, err := session.Run(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if _, err := handle.Collect(t.Context()); !errors.Is(err, failure) {
					t.Fatalf("lost hook error: %v", err)
				}

				for range 2 {
					if err := handle.Close(); err != nil {
						t.Fatal(err)
					}
				}

				if first.closes != 1 || host.reads != 0 || (phase == "pre" && host.closes != 1) || (phase == "post" && added.closes != 1) {
					t.Fatalf("cleanup: first=%d added=%d host=(reads=%d closes=%d)", first.closes, added.closes, host.reads, host.closes)
				}
			})
		}
	}
}

func TestSessionCodecAdoptsHostAcquisitionsBeforeChildHook(t *testing.T) {
	for _, base := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		for _, shape := range []string{"list", "map", "iterator"} {
			t.Run(base.ContentType()+"/"+shape, func(t *testing.T) {
				leaf := &outputIteratorLeaf{text: "leaf"}
				key := &outputIteratorLeaf{text: "key"}
				acquired := runtime.NewArrayWith(runtime.NewObjectWith(map[string]runtime.Value{"resource": leaf}))
				var value runtime.Value
				iter := &outputIterator{items: []runtime.Value{acquired}}
				switch shape {
				case "list":
					value = &outputAcquisitionList{Array: runtime.NewArrayWith(acquired)}
				case "map":
					value = &outputAcquisitionMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": acquired}), value: acquired, key: key}
				case "iterator":
					value = &outputIterable{iter: iter}
				}

				failure := errors.New("child hook failed")
				codec := base.EncodeWith().PreHook(func(ctx context.Context, item runtime.Value) error {
					if leaf.closes != 0 || key.closes != 0 {
						t.Fatal("premature child cleanup")
					}

					if item == acquired {
						if shape == "map" && ctx != context.Background() {
							t.Fatal("host predicate context replaced")
						}

						return failure
					}

					return nil
				}).Encoder().(encoding.Codec)
				eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return value, nil })
				}))
				defer eng.Close()
				plan := mustCompilePlan(t, eng, "RETURN MAKE()")
				defer plan.Close()
				session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
				defer session.Close()
				handle, err := session.Run(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if _, err := handle.Collect(t.Context()); !errors.Is(err, failure) {
					t.Fatalf("lost child error: %v", err)
				}

				for range 2 {
					if err := handle.Close(); err != nil {
						t.Fatal(err)
					}
				}

				wantKey, wantIterator := 0, 0
				if shape == "map" {
					wantKey = 1
				}

				if shape == "iterator" {
					wantIterator = 1
				}

				if leaf.closes != 1 || key.closes != wantKey || iter.closes != wantIterator || leaf.marshals != 0 {
					t.Fatalf("cleanup: leaf=%d key=%d iterator=%d marshals=%d", leaf.closes, key.closes, iter.closes, leaf.marshals)
				}
			})
		}
	}
}

func TestSessionMsgpackAdoptsMutatedBufferedContainerOnFailure(t *testing.T) {
	for _, phase := range []string{"hook", "next", "close"} {
		t.Run(phase, func(t *testing.T) {
			first := &outputIteratorLeaf{text: "first"}
			added := &outputIteratorLeaf{text: "added"}
			container := runtime.NewArrayWith(first)
			failure := errors.New("encoding stopped")
			mutate := func() {
				if err := container.Append(t.Context(), added); err != nil {
					t.Fatal(err)
				}
			}

			iter := &outputIterator{items: []runtime.Value{runtime.Int(0), container}}
			codec := encodingmsgpack.Default.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
				if phase == "hook" && value == runtime.Int(0) {
					mutate()

					return failure
				}

				return nil
			}).Encoder().(encoding.Codec)
			if phase == "next" {
				iter.nextErr = failure
				iter.onNext = func(_ context.Context, index int) {
					if index == 2 {
						mutate()
					}
				}
			}

			if phase == "close" {
				iter.closeErr = failure
				iter.onClose = mutate
			}

			eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return &outputIterable{iter: iter}, nil })
			}))
			defer eng.Close()
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			defer plan.Close()
			session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
			defer session.Close()
			handle, err := session.Run(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			if _, err := handle.Collect(t.Context()); !errors.Is(err, failure) {
				t.Fatalf("lost failure: %v", err)
			}

			for range 2 {
				if err := handle.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if first.closes != 1 || added.closes != 1 || iter.closes != 1 {
				t.Fatalf("buffer cleanup: first=%d added=%d iterator=%d", first.closes, added.closes, iter.closes)
			}
		})
	}
}

func TestSessionCodecAdoptsResourcesAppendedBySuccessfulChildHook(t *testing.T) {
	for _, base := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(base.ContentType(), func(t *testing.T) {
			first := &outputIteratorLeaf{text: "first"}
			added := &outputIteratorLeaf{text: "added"}
			root := runtime.NewArrayWith(first)
			codec := base.EncodeWith().PreHook(func(_ context.Context, value runtime.Value) error {
				if value == first {
					return root.Append(t.Context(), added)
				}

				return nil
			}).Encoder().(encoding.Codec)
			eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) { return root, nil })
			}))
			defer eng.Close()
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			defer plan.Close()
			session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
			defer session.Close()
			handle, err := session.Run(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			content, err := handle.Collect(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), base, content.Data)
			if err != nil || decoded.String() != `["first"]` {
				t.Fatalf("list snapshot changed: %v, %v", decoded, err)
			}

			for range 2 {
				if err := handle.Close(); err != nil {
					t.Fatal(err)
				}
			}

			if first.closes != 1 || added.closes != 1 || added.marshals != 0 {
				t.Fatalf("mutation cleanup: first=%d added=%d marshals=%d", first.closes, added.closes, added.marshals)
			}
		})
	}
}
