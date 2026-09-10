package objects_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func TestMutableDeepMergeBranchFailurePreservesAliases(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("branch failure")
	for _, failure := range []string{"clone", "clone type", "nested write", "nested lookup", "replacement"} {
		t.Run(failure, func(t *testing.T) {
			original := runtime.NewObjectWith(map[string]runtime.Value{"nested": runtime.NewObjectWith(map[string]runtime.Value{"old": runtime.Int(1)})})
			branch := &hostMap{Map: original}
			target := &hostMap{Map: runtime.NewObjectWith(map[string]runtime.Value{"branch": branch})}
			wantErr := sentinel
			switch failure {
			case "clone":
				branch.cloneErr = sentinel
			case "clone type":
				branch.cloned = runtime.NewArray(0)
				wantErr = runtime.ErrInvalidType
			case "nested write":
				branch.cloned = &hostMap{Map: runtime.NewObject(), setErr: sentinel}
			case "nested lookup":
				branch.cloned = &hostMap{Map: runtime.NewObject(), lookupErr: sentinel}
			case "replacement":
				target.setErr = sentinel
			}

			source := runtime.NewObjectWith(map[string]runtime.Value{"branch": runtime.NewObjectWith(map[string]runtime.Value{"nested": runtime.NewObjectWith(map[string]runtime.Value{"new": runtime.Int(2)})})})
			got, err := objects.MergeDeepMutable(ctx, target, source)
			if got != runtime.None || !errors.Is(err, wantErr) || !strings.Contains(err.Error(), "position 2") || !strings.Contains(err.Error(), "branch") {
				t.Fatalf("result=%v err=%v", got, err)
			}

			retained, err := target.Get(ctx, runtime.String("branch"))
			if err != nil || retained != branch {
				t.Fatalf("failed merge replaced original: %v", err)
			}

			assertObjectValue(t, original, runtime.NewObjectWith(map[string]runtime.Value{"nested": runtime.NewObjectWith(map[string]runtime.Value{"old": runtime.Int(1)})}))
		})
	}
}

func TestMutableMergesPreserveEarlierUpdatesOnFailure(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("incoming clone failed")
	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		clones := 0
		target := runtime.NewObject()
		first := runtime.NewObjectWith(map[string]runtime.Value{"first": runtime.True})
		second := runtime.NewObjectWith(map[string]runtime.Value{"second": &cloneValue{Value: runtime.None, clones: &clones, cloneErr: sentinel}})
		got, err := merge(ctx, target, first, second)
		if got != runtime.None || !errors.Is(err, sentinel) || !strings.Contains(err.Error(), "position 3") {
			t.Fatalf("result=%v err=%v", got, err)
		}

		assertObjectValue(t, target, first)
	}
}

func TestMutableMergesCopyIncomingValues(t *testing.T) {
	ctx := context.Background()
	for _, merge := range []runtime.Function{objects.MergeMutable, objects.MergeDeepMutable} {
		copies, clones := 0, 0
		copyLeaf := &copyOnlyValue{Value: runtime.None, copies: &copies, number: 1}
		cloneLeaf := &cloneValue{Value: runtime.None, clones: &clones, number: 1}
		branch := runtime.NewObjectWith(map[string]runtime.Value{"copy": copyLeaf, "clone": cloneLeaf})
		source := runtime.NewObjectWith(map[string]runtime.Value{"branch": branch})
		target := runtime.NewObject()
		got, err := merge(ctx, target, source)
		if err != nil || got != target || copies != 1 || clones != 1 {
			t.Fatalf("identity=%v copies=%d clones=%d err=%v", got == target, copies, clones, err)
		}

		inserted, err := target.Get(ctx, runtime.String("branch"))
		if err != nil || inserted == branch {
			t.Fatalf("source branch aliases target: %v", err)
		}

		if _, err := objects.MergeMutable(ctx, inserted, runtime.NewObjectWith(map[string]runtime.Value{"added": runtime.True})); err != nil {
			t.Fatal(err)
		}

		found, err := branch.ContainsKey(ctx, runtime.String("added"))
		if err != nil || found {
			t.Fatalf("source changed after later mutation: %v", err)
		}
	}
}

func TestMutableDeepMergeDiscardsPartiallyMergedBranch(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("second nested write")
	working := &fallibleWriteMap{hostMap: &hostMap{Map: runtime.NewObject()}, failure: sentinel}
	branch := &hostMap{Map: runtime.NewObject(), cloned: working}
	target := runtime.NewObjectWith(map[string]runtime.Value{"branch": branch})
	patch := &hostMap{Map: runtime.NewObject(), walk: func(ctx context.Context, fn runtime.KeyReadablePredicate) error {
		for _, key := range []runtime.String{"first", "second"} {
			if _, err := fn(ctx, runtime.True, key); err != nil {
				return err
			}
		}

		return nil
	}}
	got, err := objects.MergeDeepMutable(ctx, target,
		runtime.NewObjectWith(map[string]runtime.Value{"applied": runtime.True}),
		runtime.NewObjectWith(map[string]runtime.Value{"branch": patch}))
	if got != runtime.None || !errors.Is(err, sentinel) {
		t.Fatalf("nested failure: %v, %v", got, err)
	}

	assertObjectValue(t, working, runtime.NewObjectWith(map[string]runtime.Value{"first": runtime.True}))
	assertObjectValue(t, branch, runtime.NewObject())
	assertObjectValue(t, target, runtime.NewObjectWith(map[string]runtime.Value{"branch": runtime.NewObject(), "applied": runtime.True}))
	retained, err := target.Get(ctx, runtime.String("branch"))
	if err != nil || retained != branch {
		t.Fatalf("original branch replaced: %v", err)
	}
}
