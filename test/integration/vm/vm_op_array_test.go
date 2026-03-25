package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestArrayComparisonOperator(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN [1,2,3] ALL IN [2,3,4]", false),
		S("RETURN [1,2,3] ALL IN [1,2,3]", true),
		S("RETURN [1,2,3] NONE IN [3]", false),
		S("RETURN [1,2,3] NONE IN [23,42]", true),
		S("RETURN [1,2,3] ANY IN [4,5,6]", false),
		S("RETURN [1,2,3] ANY IN [1,42]", true),
		S("RETURN [1,2,3] ANY == 2", true),
		S("RETURN [1,2,3] ANY == 4", false),
		S("RETURN [1,2,3] ANY > 0", true),
		S("RETURN [1,2,3] ANY <= 1", true),
		S("RETURN [1,2,3] NONE < 99", false),
		S("RETURN [1,2,3] NONE > 10", true),
		S("RETURN [1,2,3] ALL > 2", false),
		S("RETURN [1,2,3] ALL > 0", true),
		S("RETURN [1,2,3] ALL >= 3", false),
		S(`RETURN ["foo","bar"] ALL != "moo"`, true),
		S(`RETURN ["foo","bar"] NONE == "bar"`, false),
		S(`RETURN ["foo","bar"] ANY == "foo"`, true),

		// --- Empties (vacuous truth) ---
		S("RETURN [] ANY == 1", false),
		S("RETURN [] ALL == 1", true),
		S("RETURN [] NONE == 1", true),
		S("RETURN [] ANY IN [1,2]", false),
		S("RETURN [] ALL IN [1,2]", true),
		S("RETURN [] NONE IN [1,2]", true),
		S("RETURN [1,2] ANY IN []", false),
		S("RETURN [1,2] ALL IN []", false),
		S("RETURN [1,2] NONE IN []", true),

		// --- Duplicates ---
		S("RETURN [1,1,1] ALL == 1", true),
		S("RETURN [1,1,2] NONE == 1", false),
		S("RETURN [1,1,1] ALL IN [1]", true),
		S("RETURN [1,1,2] ANY IN [2]", true),

		// --- Mixed types (strict equality expected) ---
		S(`RETURN [1,"1",true,null] ANY == 1`, true),
		S(`RETURN ["1"] ANY == 1`, false),
		S(`RETURN [1,"1"] ALL == 1`, false),
		S("RETURN [true,false] ANY == true", true),
		S("RETURN [true,false] ALL IN [true,false]", true),
		S("RETURN [true,false] NONE == null", true),

		// --- Comparator coverage ---
		S("RETURN [1,2,3] ALL >= 1", true),
		S("RETURN [1,2,3] ANY < 0", false),
		S("RETURN [1,2,3] NONE >= 4", true),
		S("RETURN [1,2,3] ALL <= 3", true),

		// --- NULL handling ---
		S("RETURN [null,1] ANY == null", true),
		S("RETURN [null,1] NONE == null", false),
		S("RETURN [null] ALL == null", true),
		S("RETURN [null] ALL IN [null]", true),
		S("RETURN [null] ANY IN [1,2]", false),
		S("RETURN [1,2] NONE IN [null]", true),

		// --- Strings ---
		S(`RETURN ["foo","bar","baz"] ALL != "foo"`, false),
		S(`RETURN ["foo","bar","baz"] ANY != "foo"`, true),
		S(`RETURN ["foo","bar","baz"] NONE == "qux"`, true),

		// --- Nested arrays (deep equality for IN) ---
		S("RETURN [[1,2],[3]] ANY IN [[3],[4]]", true),
		S("RETURN [[1],[2]] ALL IN [[1],[2],[3]]", true),

		// --- Parentheses / precedence sanity ---
		S("RETURN ([1,2,3] ANY > 2)", true),

		// --- Edge cases ---
		S("RETURN [0] ANY == 0", true),
		S("RETURN [0] ALL == 0", true),
		S("RETURN [0] NONE == 0", false),
		S("RETURN [1,2,3] ANY < 1", false),
		S("RETURN [1,2,3] ALL < 4", true),
		S("RETURN [1,2,3] NONE > 3", true),

		// --- Nested arrays and objects ---
		S("RETURN [[1,2],[3,4]] ALL IN [[1,2],[3,4]]", true),
		S("RETURN [[1,2],[3,4]] NONE IN [[5,6]]", true),
		S("RETURN [[1,2],[3,4]] ANY IN [[3,4]]", true),

		// --- Mixed types and nulls ---
		S("RETURN [null,0,false] ANY == null", true),
		S("RETURN [null,0,false] ALL == null", false),
		S("RETURN [null,0,false] NONE == null", false),
		S("RETURN [0,false] ANY == 0", true),
		S("RETURN [0,false] ANY == false", true),

		// --- Strings and case sensitivity ---
		S(`RETURN ["Foo","foo"] ANY == "foo"`, true),
		S(`RETURN ["Foo","foo"] ALL == "foo"`, false),

		// --- Large arrays ---
		S("RETURN [1,2,3,4,5,6,7,8,9,10] ALL > 0", true),
		S("RETURN [1,2,3,4,5,6,7,8,9,10] NONE < 1", true),

		// --- Empty and single-element arrays ---
		S("RETURN [42] ALL == 42", true),
		S("RETURN [42] NONE == 42", false),
		S("RETURN [] ALL != 0", true),
		S("RETURN [] NONE == 0", true),

		// --- Boolean logic ---
		S("RETURN [true,true] ALL == true", true),
		S("RETURN [true,false] ALL == true", false),
		S("RETURN [false,false] NONE == true", true),

		S("RETURN [1,2,3] AT LEAST (2) IN [2,3,4]", true).Skip("TODO: Implement in v2.1"),
		S(`RETURN ["foo","bar"] AT LEAST (1+1) == "foo"`, false).Skip("// TODO: Implement in v2.1"),
	})
}
