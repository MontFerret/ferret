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

func TestSessionOutputClosesSelfIteratorOnce(t *testing.T) {
	for _, codec := range []encoding.Codec{encodingjson.Default, encodingmsgpack.Default} {
		t.Run(codec.ContentType(), func(t *testing.T) {
			leaf := &outputIteratorLeaf{text: "leaf"}
			iterator := &outputIterator{
				items:        []runtime.Value{runtime.NewArrayWith(leaf)},
				duplicateErr: errors.New("iterator closed twice"),
			}
			eng := mustNewEngine(t,
				WithEncodingCodec(codec.ContentType(), codec),
				WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
						return iterator, nil
					})
				}),
			)
			t.Cleanup(func() { _ = eng.Close() })
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			t.Cleanup(func() { _ = plan.Close() })
			session := mustNewSession(t, plan, WithOutputContentType(codec.ContentType()))
			t.Cleanup(func() { _ = session.Close() })

			handle, err := session.Run(t.Context())
			if err != nil || handle == nil {
				t.Fatalf("Run failed: output=%v err=%v", handle, err)
			}

			if iterator.closes != 1 || leaf.closes != 1 {
				t.Fatalf("result cleanup: iterator=%d leaf=%d", iterator.closes, leaf.closes)
			}

			content, err := handle.Collect(t.Context())
			if err != nil || content == nil || !content.Metadata.LengthKnown {
				t.Fatalf("successful output: content=%v err=%v", content, err)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), codec, content.Data)
			if err != nil || decoded.String() != `[["leaf"]]` {
				t.Fatalf("output payload: value=%v err=%v", decoded, err)
			}

			for range 2 {
				for _, closeResource := range []func() error{handle.Close, session.Close, plan.Close, eng.Close} {
					if err := closeResource(); err != nil {
						t.Fatal(err)
					}
				}
			}

			if iterator.closes != 1 || leaf.closes != 1 {
				t.Fatalf("teardown repeated cleanup: iterator=%d leaf=%d", iterator.closes, leaf.closes)
			}
		})
	}
}

func TestSessionMsgpackOutputClosesYieldedResourcesWhenIteratorCloseFails(t *testing.T) {
	iteratorErr := errors.New("iterator cleanup failed")
	leafErr := errors.New("yielded resource cleanup failed")
	for _, cleanupErr := range []error{nil, leafErr} {
		name := "iterator failure"
		if cleanupErr != nil {
			name = "iterator and yielded resource failures"
		}

		t.Run(name, func(t *testing.T) {
			first := &outputIteratorLeaf{text: "first", closeErr: cleanupErr}
			second := &outputIteratorLeaf{text: "second"}
			iterator := &outputIterator{
				items: []runtime.Value{
					runtime.NewArrayWith(first),
					runtime.NewObjectWith(map[string]runtime.Value{"child": second}),
				},
				closeErr: iteratorErr,
			}
			eng := mustNewEngine(t,
				WithEncodingCodec(encodingmsgpack.ContentType, encodingmsgpack.Default),
				WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
						return &outputIterable{iter: iterator}, nil
					})
				}),
			)
			t.Cleanup(func() { _ = eng.Close() })
			plan := mustCompilePlan(t, eng, "RETURN MAKE()")
			t.Cleanup(func() { _ = plan.Close() })
			session := mustNewSession(t, plan, WithOutputContentType(encodingmsgpack.ContentType))
			t.Cleanup(func() { _ = session.Close() })

			handle, err := session.Run(t.Context())
			if err != nil || handle == nil {
				t.Fatalf("Run failed: output=%v err=%v", handle, err)
			}

			if iterator.closes != 1 || first.closes != 1 || second.closes != 1 {
				t.Fatalf("result cleanup: iterator=%d yielded=(%d,%d)", iterator.closes, first.closes, second.closes)
			}

			content, err := handle.Collect(t.Context())
			if !errors.Is(err, iteratorErr) || (cleanupErr != nil && !errors.Is(err, cleanupErr)) {
				t.Fatalf("lost error causes: %v", err)
			}

			if content == nil || len(content.Data) == 0 || content.Metadata.LengthKnown {
				t.Fatalf("failed encoding lost its bytes: %+v", content)
			}

			decoded, err := encoding.DecodeBytes(t.Context(), encodingmsgpack.Default, content.Data)
			if err != nil || decoded.String() != `[["first"],{"child":"second"}]` {
				t.Fatalf("output payload: value=%v err=%v", decoded, err)
			}

			for range 2 {
				if err := handle.Close(); !errors.Is(err, cleanupErr) || (cleanupErr == nil && err != nil) || errors.Is(err, iteratorErr) {
					t.Fatalf("output cleanup outcome: %v, want %v", err, cleanupErr)
				}
			}

			for _, closeResource := range []func() error{session.Close, plan.Close, eng.Close} {
				if err := closeResource(); err != nil {
					t.Fatal(err)
				}
			}

			if iterator.closes != 1 || first.closes != 1 || second.closes != 1 {
				t.Fatalf("teardown repeated cleanup: iterator=%d yielded=(%d,%d)", iterator.closes, first.closes, second.closes)
			}
		})
	}
}
