package arrays

import (
	"context"
	"errors"
	"math"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type capacityList struct {
	runtime.List
	failure error
	length  runtime.Int
}

func (v capacityList) Length(context.Context) (runtime.Int, error) {
	return v.length, v.failure
}

func TestArrayCapacityHints(t *testing.T) {
	input := runtime.NewArrayWith(runtime.Int(7))
	host := capacityList{List: input, length: math.MaxInt64}
	functions := []struct {
		name string
		call func() (runtime.Value, error)
		want string
	}{
		{"concat", func() (runtime.Value, error) { return Concat(t.Context(), host, input) }, "[7,7]"},
		{"union", func() (runtime.Value, error) { return Union(t.Context(), host, input) }, "[7]"},
		{"flatten", func() (runtime.Value, error) { return Flatten(t.Context(), host) }, "[7]"},
		{"unshift", func() (runtime.Value, error) { return legacyUnshift2(t.Context(), host, runtime.Int(1)) }, "[1,7]"},
	}
	for _, tc := range functions {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.call()
			if err != nil || got.String() != tc.want {
				t.Fatalf("result = %v, %v; want %s", got, err, tc.want)
			}
		})
	}

	failure := errors.New("host length failure")
	host.failure = failure
	for _, tc := range functions {
		if _, err := tc.call(); !errors.Is(err, failure) {
			t.Errorf("%s lost length error: %v", tc.name, err)
		}
	}
}
