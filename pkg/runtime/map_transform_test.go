package runtime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	objectcases "github.com/MontFerret/ferret/v2/test/spec/objects"
)

func TestMapMergeSemantics(t *testing.T) {
	ctx := context.Background()
	for _, deep := range []bool{false, true} {
		for _, test := range objectcases.MergeCases() {
			t.Run(test.Name, func(t *testing.T) {
				dst := runtime.NewObject()
				want := test.Shallow
				merge := runtime.MergeMapsInto
				if deep {
					want, merge = test.Deep, runtime.MergeMapsDeepInto
				}

				if err := merge(ctx, dst, test.Sources...); err != nil {
					t.Fatal(err)
				}

				equal, err := runtime.EqualValues(ctx, dst, want)
				if err != nil || !equal {
					t.Fatalf("deep=%v result=%v want=%v: %v", deep, dst, want, err)
				}
			})
		}
	}
}

func TestMapFilterSemantics(t *testing.T) {
	ctx := context.Background()
	for _, keep := range []bool{false, true} {
		for _, test := range objectcases.FilterCases() {
			t.Run(test.Name, func(t *testing.T) {
				apply, want := runtime.OmitMapKeys, test.Omitted
				if keep {
					apply, want = runtime.KeepMapKeys, test.Kept
				}

				if err := apply(ctx, test.Source, test.Keys...); err != nil {
					t.Fatal(err)
				}

				equal, err := runtime.EqualValues(ctx, test.Source, want)
				if err != nil || !equal {
					t.Fatalf("keep=%v result=%v want=%v: %v", keep, test.Source, want, err)
				}
			})
		}
	}
}

func TestMapDeepIsolatedMerge(t *testing.T) {
	ctx := context.Background()
	for _, test := range objectcases.MergeCases() {
		t.Run(test.Name, func(t *testing.T) {
			dst := test.Sources[0]
			if err := runtime.MergeMapsDeepIsolatedInto(ctx, dst, test.Sources[1:]...); err != nil {
				t.Fatal(err)
			}

			equal, err := runtime.EqualValues(ctx, dst, test.Deep)
			if err != nil || !equal {
				t.Fatalf("result=%v want=%v err=%v", dst, test.Deep, err)
			}
		})
	}

	leaf := runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1)})
	branch := runtime.NewObjectWith(map[string]runtime.Value{"leaf": leaf})
	dst := runtime.NewObjectWith(map[string]runtime.Value{"change": branch, "retain": branch})
	source := runtime.NewObjectWith(map[string]runtime.Value{
		"change": runtime.NewObjectWith(map[string]runtime.Value{"leaf": runtime.NewObjectWith(map[string]runtime.Value{"right": runtime.Int(2)})}),
		"copy":   branch,
	})
	if err := runtime.MergeMapsDeepIsolatedInto(ctx, dst, source); err != nil {
		t.Fatal(err)
	}

	for _, value := range []runtime.Map{branch, leaf} {
		length, err := value.Length(ctx)
		if err != nil || length != 1 {
			t.Fatalf("external alias changed: %v, %v", value, err)
		}
	}

	retained, err := dst.Get(ctx, runtime.String("retain"))
	if err != nil || retained != branch {
		t.Fatalf("untouched branch identity changed: %v", err)
	}

	changed, err := dst.Get(ctx, runtime.String("change"))
	if err != nil || changed == branch {
		t.Fatalf("conflicting branch was not detached: %v", err)
	}

	want := runtime.NewObjectWith(map[string]runtime.Value{"leaf": runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1), "right": runtime.Int(2)})})
	equal, err := runtime.EqualValues(ctx, changed, want)
	if err != nil || !equal {
		t.Fatalf("changed=%v want=%v err=%v", changed, want, err)
	}

	original, err := source.Get(ctx, runtime.String("copy"))
	if err != nil || original != branch {
		t.Fatalf("source alias changed: %v", err)
	}
}
