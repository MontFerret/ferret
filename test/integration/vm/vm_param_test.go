package vm_test

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestParam(t *testing.T) {
	RunUseCases(t,
		[]UseCase{
			CaseRuntimeErrorAs(`RETURN @foo`, runtime.Error(vm.ErrMissedParam, "@foo")),
			Case(`RETURN @str`, "bar", "Should return a value of a parameter"),
			Case(`RETURN @int + @int`, 2, "Should return a sum of two parameters"),
			Case(`RETURN @obj.str1 + @obj.str2`, "foobar", "Should return a concatenated string of two parameter properties"),
			CaseArray(`FOR i IN @values1 RETURN i`, []any{1, 2, 3, 4}, "Should iterate over an array parameter"),
			CaseArray(`FOR i IN @values2 SORT i RETURN i`, []any{"a", "b", "c", "d"}, "Should iterate over an object parameter"),
			CaseArray(`FOR i IN @start..@end RETURN i`, []any{1, 2, 3, 4, 5}, "Should iterate over a range parameter"),
			Case(`RETURN @obj.str1`, "foo", "Should be possible to use in member expression"),
		},
		vm.WithParam("str", "bar"),
		vm.WithParam("int", 1),
		vm.WithParam("bool", true),
		vm.WithParam("obj", map[string]interface{}{"str1": "foo", "str2": "bar"}),
		vm.WithParam("values1", []int{1, 2, 3, 4}),
		vm.WithParam("values2", map[string]interface{}{"a": "a", "b": "b", "c": "c", "d": "d"}),
		vm.WithParam("start", 1),
		vm.WithParam("end", 5),
	)
}
