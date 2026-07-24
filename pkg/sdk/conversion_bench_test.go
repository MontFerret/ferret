package sdk_test

import (
	"context"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type benchmarkAddress struct {
	City string `ferret:"city"`
	Zip  int    `ferret:"zip"`
}

type benchmarkUser struct {
	Name    string           `ferret:"name"`
	Address benchmarkAddress `ferret:"address"`
	Tags    []string         `ferret:"tags"`
}

var (
	benchmarkInput = benchmarkUser{
		Name: "Alice",
		Address: benchmarkAddress{
			City: "Paris",
			Zip:  75001,
		},
		Tags: []string{"admin", "author", "reviewer"},
	}
	benchmarkValue runtime.Value
	benchmarkBound = sdk.Bind2(func(_ context.Context, left runtime.String, right runtime.Int) (runtime.Value, error) {
		return runtime.NewString(left.String() + right.String()), nil
	})
)

func BenchmarkEncodeNested(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		var err error
		benchmarkValue, err = sdk.Encode(b.Context(), benchmarkInput)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNested(b *testing.B) {
	input, err := sdk.Encode(b.Context(), benchmarkInput)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var output benchmarkUser
		if err := sdk.Decode(b.Context(), input, &output); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeNestedStrict(b *testing.B) {
	input, err := sdk.Encode(b.Context(), benchmarkInput)
	if err != nil {
		b.Fatal(err)
	}

	options := []sdk.DecodeOption{
		sdk.RequireType(runtime.TypeMap),
		sdk.OnlyFields("name", "address", "tags"),
		sdk.DisallowUnknownFields(),
		sdk.DisallowNoneValues(),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		var output benchmarkUser
		if err := sdk.Decode(b.Context(), input, &output, options...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTypedBinder(b *testing.B) {
	left := runtime.NewString("value")
	right := runtime.NewInt(42)
	b.ReportAllocs()

	for b.Loop() {
		var err error
		benchmarkValue, err = benchmarkBound(b.Context(), left, right)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkManualCasts(b *testing.B) {
	left := runtime.NewString("value")
	right := runtime.NewInt(42)
	b.ReportAllocs()

	for b.Loop() {
		var err error
		benchmarkValue, err = benchmarkManualCasts(b.Context(), left, right)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSliceViewIteration(b *testing.B) {
	view := sdk.NewSliceView([]int{1, 2, 3, 4, 5, 6, 7, 8})
	b.ReportAllocs()

	for b.Loop() {
		iterator, err := view.Iterate(b.Context())
		if err != nil {
			b.Fatal(err)
		}

		for {
			benchmarkValue, _, err = iterator.Next(b.Context())
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchmarkManualCasts(_ context.Context, left, right runtime.Value) (runtime.Value, error) {
	leftValue, err := runtime.CastArg[runtime.String](left, 0)
	if err != nil {
		return runtime.None, err
	}

	rightValue, err := runtime.CastArg[runtime.Int](right, 1)
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(leftValue.String() + rightValue.String()), nil
}
