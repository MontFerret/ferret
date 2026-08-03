package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestStdlibArityOverloads(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`RETURN RAND() >= 0`, true, "zero-argument overload"),
		S(`RETURN TRIM(" value ")`, "value", "one-argument overload"),
		S(`RETURN TRIM("xxvaluexx", "x")`, "value", "two-argument overload"),
		S(`RETURN SUBSTITUTE("aba", "a", "x")`, "xbx", "three-argument overload"),
		S(`RETURN SUBSTITUTE("aaa", "a", "x", 2)`, "xxa", "four-argument overload"),
		S(`RETURN TO_DATETIME("1970-01-01T00:00:00Z") == TO_DATETIME(0, "s")`, true, "conversion overloads"),
		Nil(`RETURN T::EQ(1, 1)`, "positive assertion overload"),
		Nil(`RETURN T::EQ(1, 1, "unused")`, "positive assertion message overload"),
		Nil(`RETURN T::NOT::EQ(1, 2)`, "negative assertion overload"),
		Array(`RETURN REGEX_SPLIT("a,b,c", ",", 2, TRUE)`, []any{"a", "b,c"}, "preserve four-argument regex split behavior"),
	})
}

func TestStdlibOverloadArityErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(`RETURN TRIM()`, "bounded overload rejects missing arguments").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Message: "invalid number of arguments",
				Contains: []string{
					"wrong number of arguments in call to TRIM",
					"Note: TRIM expects 1 or 2 arguments, but got 0",
					"Hint: Pass 1 or 2 arguments to TRIM",
				},
			},
		),
		spec.NewSpec(`RETURN SLICE([1, 2], 0, 1, TRUE)`, "slice rejects undocumented extra arguments").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Message: "invalid number of arguments",
				Contains: []string{
					"wrong number of arguments in call to SLICE",
					"Note: SLICE expects 2 or 3 arguments, but got 4",
					"Hint: Pass 2 or 3 arguments to SLICE",
				},
			},
		),
	})
}
