package vm_test

import (
	"testing"
)

func TestArrayComparisonOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN [1,2,3] ALL IN [2,3,4]", false),
		Case("RETURN [1,2,3] ALL IN [1,2,3]", true),
		Case("RETURN [1,2,3] NONE IN [3]", false),
		Case("RETURN [1,2,3] NONE IN [23,42]", true),
		Case("RETURN [1,2,3] ANY IN [4,5,6]", false),
		Case("RETURN [1,2,3] ANY IN [1,42]", true),
		Case("RETURN [1,2,3] ANY == 2", true),
		Case("RETURN [1,2,3] ANY == 4", false),
		Case("RETURN [1,2,3] ANY > 0", true),
		Case("RETURN [1,2,3] ANY <= 1", true),
		Case("RETURN [1,2,3] NONE < 99", false),
		Case("RETURN [1,2,3] NONE > 10", true),
		Case("RETURN [1,2,3] ALL > 2", false),
		Case("RETURN [1,2,3] ALL > 0", true),
		Case("RETURN [1,2,3] ALL >= 3", false),
		Case(`RETURN ["foo","bar"] ALL != "moo"`, true),
		Case(`RETURN ["foo","bar"] NONE == "bar"`, false),
		Case(`RETURN ["foo","bar"] ANY == "foo"`, true),

		// --- Empties (vacuous truth) ---
		Case("RETURN [] ANY == 1", false),
		Case("RETURN [] ALL == 1", true),
		Case("RETURN [] NONE == 1", true),
		Case("RETURN [] ANY IN [1,2]", false),
		Case("RETURN [] ALL IN [1,2]", true),
		Case("RETURN [] NONE IN [1,2]", true),
		Case("RETURN [1,2] ANY IN []", false),
		Case("RETURN [1,2] ALL IN []", false),
		Case("RETURN [1,2] NONE IN []", true),

		// --- Duplicates ---
		Case("RETURN [1,1,1] ALL == 1", true),
		Case("RETURN [1,1,2] NONE == 1", false),
		Case("RETURN [1,1,1] ALL IN [1]", true),
		Case("RETURN [1,1,2] ANY IN [2]", true),

		// --- Mixed types (strict equality expected) ---
		Case(`RETURN [1,"1",true,null] ANY == 1`, true),
		Case(`RETURN ["1"] ANY == 1`, false),
		Case(`RETURN [1,"1"] ALL == 1`, false),
		Case("RETURN [true,false] ANY == true", true),
		Case("RETURN [true,false] ALL IN [true,false]", true),
		Case("RETURN [true,false] NONE == null", true),

		// --- Comparator coverage ---
		Case("RETURN [1,2,3] ALL >= 1", true),
		Case("RETURN [1,2,3] ANY < 0", false),
		Case("RETURN [1,2,3] NONE >= 4", true),
		Case("RETURN [1,2,3] ALL <= 3", true),

		// --- NULL handling ---
		Case("RETURN [null,1] ANY == null", true),
		Case("RETURN [null,1] NONE == null", false),
		Case("RETURN [null] ALL == null", true),
		Case("RETURN [null] ALL IN [null]", true),
		Case("RETURN [null] ANY IN [1,2]", false),
		Case("RETURN [1,2] NONE IN [null]", true),

		// --- Strings ---
		Case(`RETURN ["foo","bar","baz"] ALL != "foo"`, false),
		Case(`RETURN ["foo","bar","baz"] ANY != "foo"`, true),
		Case(`RETURN ["foo","bar","baz"] NONE == "qux"`, true),

		// --- Nested arrays (deep equality for IN) ---
		Case("RETURN [[1,2],[3]] ANY IN [[3],[4]]", true),
		Case("RETURN [[1],[2]] ALL IN [[1],[2],[3]]", true),

		// --- Parentheses / precedence sanity ---
		Case("RETURN ([1,2,3] ANY > 2)", true),

		// --- Edge cases ---
		Case("RETURN [0] ANY == 0", true),
		Case("RETURN [0] ALL == 0", true),
		Case("RETURN [0] NONE == 0", false),
		Case("RETURN [1,2,3] ANY < 1", false),
		Case("RETURN [1,2,3] ALL < 4", true),
		Case("RETURN [1,2,3] NONE > 3", true),

		// --- Nested arrays and objects ---
		Case("RETURN [[1,2],[3,4]] ALL IN [[1,2],[3,4]]", true),
		Case("RETURN [[1,2],[3,4]] NONE IN [[5,6]]", true),
		Case("RETURN [[1,2],[3,4]] ANY IN [[3,4]]", true),

		// --- Mixed types and nulls ---
		Case("RETURN [null,0,false] ANY == null", true),
		Case("RETURN [null,0,false] ALL == null", false),
		Case("RETURN [null,0,false] NONE == null", false),
		Case("RETURN [0,false] ANY == 0", true),
		Case("RETURN [0,false] ANY == false", true),

		// --- Strings and case sensitivity ---
		Case(`RETURN ["Foo","foo"] ANY == "foo"`, true),
		Case(`RETURN ["Foo","foo"] ALL == "foo"`, false),

		// --- Large arrays ---
		Case("RETURN [1,2,3,4,5,6,7,8,9,10] ALL > 0", true),
		Case("RETURN [1,2,3,4,5,6,7,8,9,10] NONE < 1", true),

		// --- Empty and single-element arrays ---
		Case("RETURN [42] ALL == 42", true),
		Case("RETURN [42] NONE == 42", false),
		Case("RETURN [] ALL != 0", true),
		Case("RETURN [] NONE == 0", true),

		// --- Boolean logic ---
		Case("RETURN [true,true] ALL == true", true),
		Case("RETURN [true,false] ALL == true", false),
		Case("RETURN [false,false] NONE == true", true),

		SkipCase("RETURN [1,2,3] AT LEAST (2) IN [2,3,4]", true),        // TODO: Implement in v2.1
		SkipCase(`RETURN ["foo","bar"] AT LEAST (1+1) == "foo"`, false), // TODO: Implement in v2.1
	})
}
