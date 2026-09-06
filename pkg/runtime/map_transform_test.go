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
