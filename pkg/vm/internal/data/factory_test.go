package data

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	_ runtime.Factory[runtime.List] = (*DataSet)(nil)
	_ runtime.Factory[runtime.Map]  = (*FastObject)(nil)
)

func TestDataSetNewPreservesIndependentDistinctPolicy(t *testing.T) {
	for _, distinct := range []bool{false, true} {
		source := NewDataSet(distinct)
		if err := source.Append(t.Context(), runtime.Int(1)); err != nil {
			t.Fatal(err)
		}

		result, err := source.New(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		destination, ok := result.(*DataSet)
		if !ok || destination == source || (destination.uniqueness != nil) != distinct {
			t.Fatal("DataSet policy/family changed")
		}

		if size, err := result.Length(t.Context()); size != 0 || err != nil {
			t.Fatal("new dataset is not empty")
		}

		for range 2 {
			if err := result.Append(t.Context(), runtime.Int(1)); err != nil {
				t.Fatal(err)
			}
		}

		want := runtime.Int(2)
		if distinct {
			want = 1
		}

		if size, err := result.Length(t.Context()); size != want || err != nil {
			t.Fatalf("distinct state shared or policy lost: %d, %v", size, err)
		}

		if size, err := source.Length(t.Context()); size != 1 || err != nil {
			t.Fatal("source changed")
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if result, err := source.New(ctx); result != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation: %v, %v", result, err)
		}
	}
}

func TestFastObjectNewPreservesConfiguration(t *testing.T) {
	for _, threshold := range []int{0, 1, 10} {
		source := NewFastObject(NewShapeCache(8), threshold)
		if err := source.Set(t.Context(), runtime.String("a"), runtime.Int(7)); err != nil {
			t.Fatal(err)
		}

		result, err := source.New(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		destination, ok := result.(*FastObject)
		if !ok || destination == source || destination.cache != source.cache || destination.dictThreshold != threshold {
			t.Fatal("configuration/family changed")
		}

		if size, err := result.Length(t.Context()); size != 0 || err != nil {
			t.Fatal("new map is not empty")
		}

		if err := result.Set(t.Context(), runtime.String("a"), runtime.Int(9)); err != nil {
			t.Fatal(err)
		}

		if value, err := source.Get(t.Context(), runtime.String("a")); value != runtime.Int(7) || err != nil {
			t.Fatal("source changed")
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		if result, err := source.New(ctx); result != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("cancellation: %v, %v", result, err)
		}
	}
}
