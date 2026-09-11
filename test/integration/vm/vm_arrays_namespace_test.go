package vm_test

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestArraysHostCopyErrors(t *testing.T) {
	source := &invalidArrayCopy{List: runtime.NewArrayWith(runtime.Int(1), runtime.Int(1))}
	expectCopyError := func(t *testing.T, actual any, _ ...any) {
		t.Helper()

		err, ok := actual.(error)
		if !ok || !errors.Is(err, runtime.ErrInvalidType) {
			t.Fatalf("expected copy type error, got %v", actual)
		}
	}

	var cases []spec.Spec
	for _, expression := range []string{
		`arrays::append(@source,1)`, `append(@source,1)`, `push(@source,1)`,
		`append(@source,1,true)`, `append(@source,1,false)`,
		`push(@source,1,true)`, `push(@source,1,false)`,
		`arrays::remove_at(@source,0)`, `arrays::remove_at(@source,-1)`,
		`remove_nth(@source,0)`, `remove_nth(@source,-1)`,
	} {
		cases = append(cases, spec.NewSpec("RETURN "+expression).Expect().ExecError(expectCopyError))
	}

	RunSpecs(t, cases, vm.WithParam("source", source))
	if source.String() != "[1,1]" {
		t.Fatalf("source changed: %v", source)
	}
}

func TestArraysNamespace(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`return [arrays::first([3,1]), arrays::last([3,1]), arrays::at([3,1],1)]`, []any{3, 1, 1}),
		Array(`return [arrays::first([]),arrays::last([]),arrays::at([],0),arrays::at([1],-1),arrays::at([1],9223372036854775807)]`, []any{nil, nil, nil, nil, nil}),
		Array(`return arrays::append([1,1], [2,3])`, []any{1, 1, []any{2, 3}}),
		Array(`return arrays::concat([3,1,3], [1,2])`, []any{3, 1, 3, 1, 2}),
		Array(`return arrays::union([3,1,3], [1,2])`, []any{3, 1, 2}),
		Array(`return arrays::intersection([3,1,3,2], [2,3,1], [1,3])`, []any{3, 1}),
		Array(`return arrays::difference([3,1,3,2], [2], [1])`, []any{3}),
		Array(`return arrays::symmetric_difference([7,3,7], [7,2,2], [7,2,5])`, []any{7, 3, 5}),
		Array(`return arrays::symmetric_difference([7,3], [7], [7], [7])`, []any{3}),
		Array(`return [arrays::contains([1,2,1],1),arrays::contains([],1),arrays::index_of([1,2,1],1),arrays::index_of([],1)]`, []any{true, false, 0, -1}),
		Array(`return arrays::unique([3,1,3,2])`, []any{3, 1, 2}),
		Array(`return arrays::sorted(arrays::unique([3,1,3,2]))`, []any{1, 2, 3}),
		Array(`return arrays::remove([3,1,3,2],3)`, []any{1, 2}),
		Array(`return arrays::remove_at([3,1,3,2],2)`, []any{3, 1, 2}),
		Array(`return arrays::remove_at([3,1],-1)`, []any{3, 1}),
		Array(`return arrays::remove_at([3,1],9223372036854775807)`, []any{3, 1}),
		Array(`return arrays::remove_any([3,1,3,2],[1,1,2])`, []any{3, 3}),
		Array(`return arrays::flatten([[1,[2]],[3]])`, []any{1, []any{2}, 3}),
		Array(`return arrays::flatten([[1,[2]],[3]],2)`, []any{1, 2, 3}),
		Array(`return arrays::slice([1,2,3],1)`, []any{2, 3}),
		Array(`return arrays::slice([1,2,3],1,1)`, []any{2}),
		Array(`return arrays::slice([1,2,3],1,9223372036854775807)`, []any{2, 3}),
		Array(`return arrays::slice([1,2,3],-1)`, []any{}),
		Array(`return ArRaYs::UnIoN([3,1], [3,2])`, []any{3, 1, 2}),
		Array(`let source=[4,3,2,1] let part=arrays::slice(source,1,2) part[0]=99 source[2]=88 return [source,part]`, []any{[]any{4, 3, 88, 1}, []any{99, 2}}),
		Array(`let source=[3,2,1] let part=pop(source) part[0]=99 return [source,part]`, []any{[]any{3, 2, 1}, []any{99, 2}}),
	})
}

func TestLegacyArraysCompatibility(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`return union([3,1,3],[1,2])`, []any{3, 1, 3, 1, 2}),
		Array(`return union_distinct([3,1,3],[1,2])`, []any{3, 1, 2}),
		Array(`return [first([3,1]),last([3,1]),nth([3,1],1)]`, []any{3, 1, 1}),
		Array(`return [position([2,1],1),position([2,1],1,false),position([2,1],1,true),position([],1,true)]`, []any{true, true, 1, -1}),
		Array(`return append([1,1],1,true)`, []any{1, 1}),
		Array(`return append([1,1],1,false)`, []any{1, 1, 1}),
		Array(`return append([1,1],2)`, []any{1, 1, 2}),
		Array(`return push([1,1],1,true)`, []any{1, 1}),
		Array(`return push([1,1],1,false)`, []any{1, 1, 1}),
		Array(`return push([1,1],2)`, []any{1, 1, 2}),
		Array(`return remove_value([1,2,1,1],1,1)`, []any{2, 1, 1}),
		Array(`return remove_value([1,2,1,1],1,0)`, []any{1, 2, 1, 1}),
		Array(`return remove_value([1,2,1,1],1,-2)`, []any{2}),
		Array(`return remove_value([1,2,1,1],1)`, []any{2}),
		Array(`return remove_values([1,2,1,3],[1])`, []any{2, 3}),
		Array(`return remove_nth([1,2,3],1)`, []any{1, 3}),
		Array(`return remove_nth([1],-1)`, []any{1}),
		Array(`return unshift([1,2,1,2],1,true)`, []any{1, 2, 2}),
		Array(`return unshift([1,2],1,false)`, []any{1, 1, 2}),
		Array(`return unshift([1,2],0)`, []any{0, 1, 2}),
		Array(`return sorted_unique([3,1,3,2])`, []any{1, 2, 3}),
		Array(`return sorted([3,1,3,2])`, []any{1, 2, 3, 3}),
		Array(`return unique([3,1,3,2])`, []any{3, 1, 2}),
		Array(`return pop([1])`, []any{}),
		Array(`return pop([])`, []any{}),
		Array(`return shift([1,2])`, []any{2}),
		Array(`return shift([])`, []any{}),
		Array(`return outersection([7,3,7],[7,2,2],[7,2,5])`, []any{3, 5}),
		Array(`return intersection([3,1,3],[1,3])`, []any{3, 1}),
		Array(`return minus([3,1,3],[1])`, []any{3}),
		Array(`return flatten([[1,[2]]])`, []any{1, []any{2}}),
		Array(`return flatten([[1,[2]]],2)`, []any{1, 2}),
		Array(`return slice([1,2,3],1)`, []any{2, 3}),
		Array(`return slice([1,2,3],1,1)`, []any{2}),
	})
}

func TestArraysRejectLegacyCanonicalNamesAndModes(t *testing.T) {
	var cases []spec.Spec
	for _, expression := range []string{
		`arrays::position([1],1)`, `arrays::union_distinct([1],[2])`, `arrays::minus([1],[2])`,
		`arrays::nth([1],0)`, `arrays::remove_nth([1],0)`, `arrays::outersection([1],[2])`,
		`arrays::sorted_unique([1])`, `arrays::pop([1])`, `arrays::push([1],2)`, `arrays::shift([1])`,
		`arrays::unshift([1],2)`, `arrays::mut::append([1],2)`, `arrays::remove_value([1],1)`,
		`arrays::remove_values([1],[1])`, `arrays::exclusive([1],[2])`,
	} {
		cases = append(cases, Error("return "+expression))
	}

	for _, expression := range []string{
		`arrays::append([1],1,true)`, `arrays::remove([1],1,1)`, `arrays::contains([1],1,true)`,
		`arrays::index_of([1],1,true)`, `arrays::union([1])`, `arrays::symmetric_difference([1])`,
		`arrays::remove_at([1],false)`, `arrays::at([1],false)`, `arrays::slice([1],0,1,true)`,
		`position([1])`, `remove_value([1])`, `append([1],1,0)`, `position([1],1,0)`,
	} {
		cases = append(cases, Error("return "+expression))
	}

	RunSpecs(t, cases)
}
