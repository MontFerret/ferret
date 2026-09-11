package vm_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestMutableArrayContracts(t *testing.T) {
	specs := []spec.Spec{
		S("let a=[1,2] arrays::mut::push(a,3) arrays::mut::set(a,0,9) return a==[9,2,3]", true),
		S("LET values = [1, 2] LET copied = arrays::append(values, 3) LET changed = arrays::mut::push(values, 4) LET removed = arrays::mut::pop(values) RETURN {values, copied, changed, removed} == {values: [1,2], copied: [1,2,3], changed: [1,2], removed: 4}", true),
		S("LET original = [1, 2, 3] LET part = arrays::slice(original, 0, 2) LET changed = arrays::mut::set(part, 0, 9) RETURN {original, part} == {original: [1,2,3], part: [9,2]}", true),
		S("let a=[1,2] let alias=a let b=arrays::append(a,3) let result=arrays::mut::push(a,4) return a==[1,2,4] and alias==a and result==a and b==[1,2,3]", true),
		S("let a=[] let result=arrays::mut::unshift(a,2) let inserted=arrays::mut::insert(result,0,1) let appended=arrays::mut::insert(a,2,3) let changed=arrays::mut::set(a,1,9) return a==[1,9,3] and inserted==a and appended==a and changed==a", true),
		S("let a=[1,2,3] let last=arrays::mut::pop(a) let first=arrays::mut::shift(a) let middle=arrays::mut::remove_at(a,0) return last==3 and first==1 and middle==2 and a==[]", true),
		S("let a=[] return arrays::mut::pop(a)==none and arrays::mut::shift(a)==none and a==[]", true),
		S("let a=[1,2,1,3,3,1] let copy=arrays::remove(a,1) let result=arrays::mut::remove(a,1.0) return a==[2,3,3] and result==a and copy==a", true),
		S("let a=[3,1,2,1] let b=arrays::sorted(a) let result=arrays::mut::sort(a) let cleared=arrays::mut::clear(result) return a==[] and result==[] and cleared==[] and b==[1,1,2,3]", true),
		S("let a=[4,3,2,1] let b=arrays::slice(a,1,2) let c=arrays::slice(a,0,3) let changed=arrays::mut::set(b,0,99) let grown=arrays::mut::push(b,88) return a==[4,3,2,1] and b==[99,2,88] and c==[4,3,2]", true),
		S("let a=[1,2] let b=[1,2] let changed=arrays::mut::set(b,0,9) return a==[1,2] and b==[9,2]", true),
		S("let child=[1] let a=[child] let b=arrays::slice(a,0,1) let changed=arrays::mut::push(b[0],2) return a==[[1,2]] and b==[[1,2]]", true),
		S("let a=[1,2,3] let x=push(a,4) let y=pop(a) let z=shift(a) let w=unshift(a,0) return a==[1,2,3] and x==[1,2,3,4] and y==[1,2] and z==[2,3] and w==[0,1,2,3]", true),
		S("let a=[1] let changed=ArRaYs::MuT::PuSh(a,2) return a==[1,2]", true),
	}
	for _, call := range []string{
		"arrays::mut::set([],0,1)", "arrays::mut::set([1],1,2)", "arrays::mut::set([1],-1,2)",
		"arrays::mut::insert([1],-1,2)", "arrays::mut::insert([1],2,2)",
		"arrays::mut::remove_at([],0)", "arrays::mut::remove_at([1],1)", "arrays::mut::remove_at([1],-1)",
		"arrays::mut::set([1],9223372036854775807,2)", "arrays::mut::insert([1],9223372036854775807,2)",
		"arrays::mut::remove_at([1],9223372036854775807)", "arrays::mut::remove_at([1],0.0)",
		"arrays::mut::set([1],true,2)", "arrays::mut::insert([1],none,2)",
		"arrays::mut::push(1,2)", "arrays::mut::pop(none)", "arrays::mut::sort({})",
	} {
		specs = append(specs, S("return "+call+" on error return \"caught\"", "caught", call))
	}

	for _, call := range []string{
		"mut::push([],1)", "arrays::mut::append([],1)", "arrays::mut::reverse([])",
		"arrays::mut::push([])", "arrays::mut::push([],1,true)", "arrays::mut::unshift([],1,true)",
		"arrays::mut::remove([],1,1)", "arrays::mut::sort([],true)", "arrays::mut::pop([],1)",
		"arrays::mut::shift([],1)", "arrays::mut::clear([],1)", "arrays::mut::set([],0)",
		"arrays::mut::insert([],0)", "arrays::mut::remove_at([])",
	} {
		specs = append(specs, Error("return "+call, call))
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		RunSpecsWith(t, fmt.Sprintf("mutable arrays/%s", level), c, specs)

		target := &invalidArrayCopy{List: runtime.NewArrayWith(runtime.Int(2), runtime.Int(1))}
		RunSpecsWith(t, "mutable host/"+level.String(), c, []spec.Spec{
			S("let result=arrays::mut::push(@target,3) let first=arrays::mut::shift(result) let sorted=arrays::mut::sort(result) return first==2 and @target==[1,3] and sorted==@target", true),
		}, vm.WithParam("target", target))
	}
}
