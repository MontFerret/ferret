package engine

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSessionMsgpackOutputAdoptsYieldedResourcesBeforeTraversalFailure(t *testing.T) {
	traversalErr := errors.New("iterator next failed")
	iteratorErr := errors.New("iterator close failed")
	leafErr := errors.New("yielded resource close failed")
	for _, borrowed := range []bool{false, true} {
		for _, failure := range []error{context.Canceled, context.DeadlineExceeded, traversalErr} {
			for _, shape := range []string{"direct", "containers", "deep", "cycle", "resource aliases", "nil containers"} {
				t.Run(failure.Error()+"/"+shape+map[bool]string{false: "/owned", true: "/borrowed"}[borrowed], func(t *testing.T) {
					first := &outputIteratorLeaf{text: "first", closeErr: leafErr}
					second := &outputIteratorLeaf{text: "second"}
					items := []runtime.Value{first, second, first}
					switch shape {
					case "containers":
						items = []runtime.Value{runtime.NewArrayWith(first, runtime.NewObjectWith(map[string]runtime.Value{"child": second})), first}
					case "deep":
						var nested runtime.Value = runtime.NewArrayWith(first, second)
						for range 32_768 {
							nested = runtime.NewArrayWith(nested)
						}

						items = []runtime.Value{nested}
					case "cycle":
						array := runtime.NewArrayWith(first)
						object := runtime.NewObjectWith(map[string]runtime.Value{"child": second, "array": array})
						if err := array.Append(t.Context(), object); err != nil {
							t.Fatal(err)
						}

						items = []runtime.Value{array, object}
					case "resource aliases":
						items = []runtime.Value{&outputIteratorResource{outputIteratorLeaf: first, id: 1}, &outputIteratorResource{outputIteratorLeaf: first, id: 1}, &outputIteratorResource{outputIteratorLeaf: second, id: 2}}
					case "nil containers":
						items = []runtime.Value{(*runtime.Array)(nil), (*runtime.Object)(nil), first, second}
					}

					iterator := &outputIterator{items: items, nextErr: failure, closeErr: iteratorErr, duplicateErr: errors.New("iterator closed twice")}
					var value runtime.Value = &outputIterable{iter: iterator}
					if borrowed {
						value = iterator
					}

					preCalls, postCalls := 0, 0
					var postErr error
					codec := encodingmsgpack.Default.EncodeWith().PreHook(func(_ context.Context, item runtime.Value) error {
						preCalls++
						if item != value {
							t.Fatal("failed collection ran a child pre-hook")
						}

						return nil
					}).PostHook(func(_ context.Context, item runtime.Value, err error) error {
						postCalls++
						if item != value {
							t.Fatal("failed collection ran a child post-hook")
						}

						if first.closes != 0 || second.closes != 0 {
							t.Fatal("result resources closed before parent post-hooks")
						}

						postErr = err

						return nil
					}).Encoder().(encoding.Codec)
					eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
						ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
							return value, nil
						})
					}))
					t.Cleanup(func() { _ = eng.Close() })
					plan := mustCompilePlan(t, eng, "RETURN MAKE()")
					t.Cleanup(func() { _ = plan.Close() })
					session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
					t.Cleanup(func() { _ = session.Close() })

					handle, err := session.Run(t.Context())
					if err != nil || handle == nil {
						t.Fatalf("Run: output=%v err=%v", handle, err)
					}

					if iterator.closes != 1 || first.closes != 1 || second.closes != 1 || first.marshals != 0 || second.marshals != 0 {
						t.Fatalf("result cleanup: iterator=%d children=(%d,%d) marshals=(%d,%d)", iterator.closes, first.closes, second.closes, first.marshals, second.marshals)
					}

					if preCalls != 1 || postCalls != 1 || !errors.Is(postErr, failure) || errors.Is(postErr, leafErr) || errors.Is(postErr, iteratorErr) != !borrowed {
						t.Fatalf("root hooks: pre=%d post=%d err=%v", preCalls, postCalls, postErr)
					}

					content, err := handle.Collect(t.Context())
					if content != nil || !errors.Is(err, failure) || !errors.Is(err, iteratorErr) || !errors.Is(err, leafErr) {
						t.Fatalf("collection outcome: content=%v err=%v", content, err)
					}

					for range 2 {
						err := handle.Close()
						if !errors.Is(err, leafErr) || errors.Is(err, failure) || errors.Is(err, iteratorErr) != borrowed {
							t.Fatalf("cleanup outcome: %v", err)
						}

						for _, closeResource := range []func() error{session.Close, plan.Close, eng.Close} {
							if err := closeResource(); err != nil {
								t.Fatal(err)
							}
						}
					}

					if iterator.closes != 1 || first.closes != 1 || second.closes != 1 {
						t.Fatalf("teardown repeated cleanup: iterator=%d children=(%d,%d)", iterator.closes, first.closes, second.closes)
					}
				})
			}
		}
	}
}

func TestSessionMsgpackOutputAdoptsValueAtCancellationCheckpoint(t *testing.T) {
	for _, shape := range []string{"direct", "containers"} {
		t.Run(shape, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			leaf := &outputIteratorLeaf{text: "last yielded value"}
			items := make([]runtime.Value, 255)
			for i := range items {
				items[i] = runtime.Int(i)
			}

			items[len(items)-1] = leaf
			if shape == "containers" {
				array := runtime.NewArrayWith(leaf)
				object := runtime.NewObjectWith(map[string]runtime.Value{"array": array})
				if err := array.Append(t.Context(), object); err != nil {
					t.Fatal(err)
				}

				items[len(items)-1] = object
			}

			iterator := &outputIterator{items: items, onNext: func(_ context.Context, index int) {
				if index == 254 {
					cancel()
				}
			}}
			eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
				ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
					return &outputIterable{iter: iterator}, nil
				})
			}))
			t.Cleanup(func() { _ = eng.Close() })
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			t.Cleanup(func() { _ = plan.Close() })
			session := mustNewSession(t, plan, WithOutputContentType(encodingmsgpack.ContentType))
			t.Cleanup(func() { _ = session.Close() })

			handle, err := session.Run(ctx)
			if err != nil || handle == nil {
				t.Fatalf("Run: output=%v err=%v", handle, err)
			}

			content, err := handle.Collect(t.Context())
			if content != nil || !errors.Is(err, context.Canceled) || iterator.index != 255 || iterator.closes != 1 || leaf.closes != 1 || leaf.marshals != 0 {
				t.Fatalf("checkpoint cleanup: content=%v err=%v yielded=%d iterator=%d leaf=%d marshals=%d", content, err, iterator.index, iterator.closes, leaf.closes, leaf.marshals)
			}

			if err := handle.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSessionMsgpackEarlyAdoptionDoesNotReadHostCollections(t *testing.T) {
	leaf := &outputIteratorLeaf{text: "unvisited"}
	lazy := &outputIterator{items: []runtime.Value{leaf}}
	hostList := &outputIteratorHostList{Array: runtime.NewArrayWith(leaf)}
	iterator := &outputIterator{items: []runtime.Value{runtime.NewArrayWith(lazy, hostList)}, nextErr: context.Canceled}
	eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return &outputIterable{iter: iterator}, nil
		})
	}))
	t.Cleanup(func() { _ = eng.Close() })
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	t.Cleanup(func() { _ = plan.Close() })
	session := mustNewSession(t, plan, WithOutputContentType(encodingmsgpack.ContentType))
	t.Cleanup(func() { _ = session.Close() })

	handle, err := session.Run(t.Context())
	if err != nil || handle == nil {
		t.Fatalf("Run: output=%v err=%v", handle, err)
	}

	content, err := handle.Collect(t.Context())
	if content != nil || !errors.Is(err, context.Canceled) || lazy.index != 0 || hostList.reads != 0 || lazy.closes != 1 || hostList.closes != 1 || leaf.closes != 0 || leaf.marshals != 0 {
		t.Fatalf("host collection ownership: content=%v err=%v lazy=(reads %d, closes %d) list=(reads %d, closes %d) leaf=(closes %d, marshals %d)", content, err, lazy.index, lazy.closes, hostList.reads, hostList.closes, leaf.closes, leaf.marshals)
	}

	if err := handle.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionMsgpackEarlyAdoptionRescansMutatedYieldedContainer(t *testing.T) {
	first, second := &outputIteratorLeaf{text: "first"}, &outputIteratorLeaf{text: "second"}
	container := runtime.NewArrayWith(first)
	iterator := &outputIterator{items: []runtime.Value{container, container}, nextErr: context.Canceled, onNext: func(ctx context.Context, index int) {
		if index == 1 {
			if err := container.Append(ctx, second); err != nil {
				t.Fatal(err)
			}
		}
	}}
	eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return &outputIterable{iter: iterator}, nil
		})
	}))
	t.Cleanup(func() { _ = eng.Close() })
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	t.Cleanup(func() { _ = plan.Close() })
	session := mustNewSession(t, plan, WithOutputContentType(encodingmsgpack.ContentType))
	t.Cleanup(func() { _ = session.Close() })

	handle, err := session.Run(t.Context())
	if err != nil || handle == nil {
		t.Fatalf("Run: output=%v err=%v", handle, err)
	}

	content, err := handle.Collect(t.Context())
	if content != nil || !errors.Is(err, context.Canceled) || first.closes != 1 || second.closes != 1 {
		t.Fatalf("mutated container ownership: content=%v err=%v children=(%d,%d)", content, err, first.closes, second.closes)
	}

	if err := handle.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestSessionMsgpackCollectedResourcesSurviveChildHookFailure(t *testing.T) {
	first, second := &outputIteratorLeaf{text: "first"}, &outputIteratorLeaf{text: "second"}
	iterator := &outputIterator{items: []runtime.Value{first, runtime.NewArrayWith(second)}}
	hookErr := errors.New("child hook failed")
	codec := encodingmsgpack.Default.EncodeWith().PreHook(func(_ context.Context, item runtime.Value) error {
		if item == first {
			return hookErr
		}

		return nil
	}).Encoder().(encoding.Codec)
	eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return &outputIterable{iter: iterator}, nil
		})
	}))
	t.Cleanup(func() { _ = eng.Close() })
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	t.Cleanup(func() { _ = plan.Close() })
	session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
	t.Cleanup(func() { _ = session.Close() })

	handle, err := session.Run(t.Context())
	if err != nil || handle == nil {
		t.Fatalf("Run: output=%v err=%v", handle, err)
	}

	content, err := handle.Collect(t.Context())
	if !errors.Is(err, hookErr) || content == nil || content.Metadata.LengthKnown || first.closes != 1 || second.closes != 1 || first.marshals != 0 || second.marshals != 0 || iterator.closes != 1 {
		t.Fatalf("child hook settlement: content=%v err=%v children=(%d,%d) marshals=(%d,%d) iterator=%d", content, err, first.closes, second.closes, first.marshals, second.marshals, iterator.closes)
	}

	if err := handle.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMsgpackDirectEncodingLeavesYieldedResourcesWithCaller(t *testing.T) {
	for _, borrowed := range []bool{false, true} {
		for _, failure := range []error{nil, context.Canceled} {
			name := map[bool]string{false: "owned iterator", true: "borrowed iterator"}[borrowed]
			if failure != nil {
				name += "/failure"
			}

			t.Run(name, func(t *testing.T) {
				leaf := &outputIteratorLeaf{text: "yielded"}
				iterator := &outputIterator{items: []runtime.Value{runtime.NewArrayWith(leaf)}, nextErr: failure}
				var value runtime.Value = &outputIterable{iter: iterator}
				if borrowed {
					value = iterator
				}

				data, err := encoding.EncodeBytes(t.Context(), encodingmsgpack.Default, value)
				wantCloses, wantMarshals := 1, 1
				if borrowed {
					wantCloses = 0
				}

				if failure != nil {
					wantMarshals = 0
				}

				if (len(data) == 0) != (failure != nil) || !errors.Is(err, failure) || leaf.closes != 0 || leaf.marshals != wantMarshals || iterator.closes != wantCloses {
					t.Fatalf("direct encoding: data=%v err=%v iterator=%d leaf=%d marshals=%d", data, err, iterator.closes, leaf.closes, leaf.marshals)
				}
			})
		}
	}
}

func TestSessionMsgpackEarlyAdoptionPreservesHookAndCleanupOrder(t *testing.T) {
	var events []string
	first := &outputIteratorLeaf{text: "first", onClose: func() { events = append(events, "close first") }}
	second := &outputIteratorLeaf{text: "second", onClose: func() { events = append(events, "close second") }}
	iterator := &outputIterator{items: []runtime.Value{runtime.NewArrayWith(first, second)}, onClose: func() { events = append(events, "close iterator") }}
	value := &outputIterable{iter: iterator}
	label := func(item runtime.Value) string {
		if item == value {
			return "root"
		}

		if _, ok := item.(*runtime.Array); ok {
			return "array"
		}

		return item.String()
	}
	codec := encodingmsgpack.Default.EncodeWith().PreHook(func(_ context.Context, item runtime.Value) error {
		events = append(events, "pre "+label(item))

		return nil
	}).PostHook(func(_ context.Context, item runtime.Value, err error) error {
		if err != nil {
			t.Fatal(err)
		}

		events = append(events, "post "+label(item))

		return nil
	}).Encoder().(encoding.Codec)
	eng := mustNewEngine(t, WithEncodingCodec(codec.ContentType(), codec), WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return value, nil
		})
	}))
	t.Cleanup(func() { _ = eng.Close() })
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	t.Cleanup(func() { _ = plan.Close() })
	session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
	t.Cleanup(func() { _ = session.Close() })

	handle, err := session.Run(t.Context())
	if err != nil || handle == nil {
		t.Fatalf("Run: output=%v err=%v", handle, err)
	}

	content, err := handle.Collect(t.Context())
	if err != nil || content == nil || !content.Metadata.LengthKnown {
		t.Fatalf("successful output: content=%v err=%v", content, err)
	}

	want := []string{"pre root", "pre array", "pre first", "post first", "pre second", "post second", "post array", "close iterator", "post root", "close first", "close second"}
	if !slices.Equal(events, want) || first.closes != 1 || second.closes != 1 || first.marshals != 1 || second.marshals != 1 {
		t.Fatalf("hook and cleanup order: got %v, want %v", events, want)
	}

	if err := handle.Close(); err != nil {
		t.Fatal(err)
	}
}
