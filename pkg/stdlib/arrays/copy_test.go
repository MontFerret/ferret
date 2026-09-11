package arrays_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestArrayCopyFailures(t *testing.T) {
	for _, test := range []struct {
		name string
		args []runtime.Value
	}{
		{name: "arrays::append", args: []runtime.Value{runtime.Int(1)}},
		{name: "append", args: []runtime.Value{runtime.Int(1)}},
		{name: "push", args: []runtime.Value{runtime.Int(1)}},
		{name: "append", args: []runtime.Value{runtime.Int(1), runtime.True}},
		{name: "append", args: []runtime.Value{runtime.Int(1), runtime.False}},
		{name: "push", args: []runtime.Value{runtime.Int(1), runtime.True}},
		{name: "push", args: []runtime.Value{runtime.Int(1), runtime.False}},
		{name: "arrays::remove_at", args: []runtime.Value{runtime.Int(0)}},
		{name: "arrays::remove_at", args: []runtime.Value{runtime.Int(-1)}},
		{name: "arrays::remove_at", args: []runtime.Value{runtime.Int(9)}},
		{name: "remove_nth", args: []runtime.Value{runtime.Int(0)}},
		{name: "remove_nth", args: []runtime.Value{runtime.Int(-1)}},
	} {
		t.Run(test.name, func(t *testing.T) {
			probe := &copyMutationProbe{Value: runtime.True}
			source := &copyFailureList{List: runtime.NewArrayWith(runtime.Int(1), runtime.Int(1)), copied: probe}
			args := append([]runtime.Value{source}, test.args...)
			result, err := callLegacy(test.name, t.Context(), args...)
			if result != runtime.None || !errors.Is(err, runtime.ErrInvalidType) {
				t.Fatalf("result = %v, error = %v", result, err)
			}

			if source.copies != 1 || source.lookups != 0 || probe.mutations != 0 || source.String() != "[1,1]" {
				t.Fatalf("copies = %d, lookups = %d, mutations = %d, source = %v", source.copies, source.lookups, probe.mutations, source)
			}
		})
	}
}
