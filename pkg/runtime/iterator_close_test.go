package runtime_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestForEachPreservesTraversalAndCloseErrors(t *testing.T) {
	primary, cleanup := errors.New("traversal failed"), errors.New("close failed")
	for _, stage := range []string{"create", "next", "predicate", "close", "both next", "both predicate", "eof", "stop"} {
		t.Run(stage, func(t *testing.T) {
			iter := &traversalIterator{}
			src := &traversalSource{iter: iter}
			var predicateErr error
			wantPrimary, wantCleanup := false, false
			switch stage {
			case "create":
				src.err, wantPrimary = primary, true
			case "next":
				iter.nextErr, wantPrimary = primary, true
			case "predicate":
				predicateErr, wantPrimary = primary, true
			case "close":
				iter.closeErr, wantCleanup = cleanup, true
			case "both next":
				iter.nextErr, iter.closeErr, wantPrimary, wantCleanup = primary, cleanup, true, true
			case "both predicate":
				predicateErr, iter.closeErr, wantPrimary, wantCleanup = primary, cleanup, true, true
			case "eof":
				iter.nextErr = io.EOF
			}

			err := runtime.ForEach(context.Background(), src, func(context.Context, runtime.Value, runtime.Value) (runtime.Boolean, error) {
				return false, predicateErr
			})
			if errors.Is(err, primary) != wantPrimary || errors.Is(err, cleanup) != wantCleanup {
				t.Fatalf("lost error: %v", err)
			}

			if !wantPrimary && !wantCleanup && err != nil {
				t.Fatal(err)
			}

			wantCloses := 1
			if stage == "create" {
				wantCloses = 0
			}

			if iter.closes != wantCloses {
				t.Fatalf("closed %d times, want %d", iter.closes, wantCloses)
			}
		})
	}
}
