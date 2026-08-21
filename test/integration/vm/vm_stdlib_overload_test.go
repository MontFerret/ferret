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
		spec.NewSpec(`RETURN T::EQ(1, 2, "custom assertion message")`, "assertion custom message").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "assertion error", Contains: []string{"custom assertion message"}},
		),
		Array(`RETURN REGEX_SPLIT("a,b,c", ",", 2, TRUE)`, []any{"a", "b,c"}, "preserve four-argument regex split behavior"),
		Array(`RETURN [FIND_FIRST("foobarbaz", "ba", 4, 9), FIND_LAST("foobarbaz", "ba", 4, 6)]`, []any{6, -1}, "four-argument string searches honor start and end"),
		Array(`RETURN [FIND_FIRST("éaéa", "a"), FIND_FIRST("éaéa", "a", 2), FIND_LAST("éaéa", "a"), FIND_LAST("éaéa", "a", 2)]`, []any{1, 3, 3, 3}, "string searches return Unicode character positions"),
		Array(`RETURN [FIND_FIRST("éaéa", "a", -10, 100), FIND_LAST("éaéa", "a", -10, 100), FIND_FIRST("éaéa", "a", 3, 1), FIND_LAST("éaéa", "a", 3, 1), FIND_FIRST("éa", "", 1, 100), FIND_LAST("éa", "", 1, 100)]`, []any{1, 3, -1, -1, 1, 2}, "string searches clamp character bounds"),
		Array(`RETURN RANGE(-3, -1)`, []any{-3, -2, -1}, "two-argument range supports negative endpoints"),
		Array(`RETURN RANGE(3, 1, -1)`, []any{3, 2, 1}, "three-argument range supports descending steps"),
	})
}

func TestTestingAssertionVocabulary(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Nil(`RETURN T::BOOL(FALSE)`, "boolean type assertion accepts false"),
		Nil(`RETURN T::NOT::BOOL("false")`, "negated boolean type assertion"),
		Nil(`RETURN T::NUMBER(42)`, "number assertion accepts integers"),
		Nil(`RETURN T::NUMBER(42.5)`, "number assertion accepts floats"),
		Nil(`RETURN T::DURATION(1s)`, "duration assertion accepts Duration values"),
		Nil(`RETURN T::NOT::DURATION("1s")`, "duration assertion rejects duration-like strings"),
		Nil(`RETURN T::APPROX(10, 10.5, 0.5, "unused")`, "approx supports inclusive mixed numeric boundary and message overload"),
		Nil(`RETURN T::NOT::APPROX(10, 11, 0.01)`, "negated approximate assertion"),
		Nil(`RETURN T::BETWEEN(200, 200, 299, "unused")`, "between includes its minimum and supports message overload"),
		Nil(`RETURN T::NOT::BETWEEN(500, 200, 299)`, "negated range assertion"),
		Nil(`RETURN T::CONTAINS(["ferret", "fql"], "ferret")`, "containment assertion"),
		Nil(`RETURN T::NOT::CONTAINS("ferret", "goose")`, "negated containment assertion"),
		Nil(`RETURN T::HAS({ id: 1, name: "Ferret", empty: NONE }, ["id", "name", "empty"])`, "has checks all properties including present None"),
		Nil(`RETURN T::HAS({}, [])`, "has accepts an empty key list"),
		spec.NewSpec(`RETURN T::APPROX(10, 11, 0.01, "not close")`, "approx custom message").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "assertion error", Contains: []string{"not close"}},
		),
		spec.NewSpec(`RETURN T::APPROX(10, 10, -1)`, "approx rejects negative tolerance").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid argument", Contains: []string{"tolerance must be non-negative"}},
		),
		spec.NewSpec(`RETURN T::BETWEEN(15, 20, 10)`, "between rejects reversed bounds").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{Message: "invalid argument", Contains: []string{"minimum boundary must not exceed maximum boundary"}},
		),
	})
}

func TestStdlibOverloadArityErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.NewSpec(`RETURN TRIM()`, "bounded overload rejects missing arguments").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Message: "invalid number of arguments",
				Contains: []string{
					"wrong number of arguments in call to trim",
					"Note: trim expects 1 or 2 arguments, but got 0",
					"Hint: Pass 1 or 2 arguments to trim",
				},
			},
		),
		spec.NewSpec(`RETURN SLICE([1, 2], 0, 1, TRUE)`, "slice rejects undocumented extra arguments").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Message: "invalid number of arguments",
				Contains: []string{
					"wrong number of arguments in call to slice",
					"Note: slice expects 2 or 3 arguments, but got 4",
					"Hint: Pass 2 or 3 arguments to slice",
				},
			},
		),
		spec.NewSpec(`RETURN RANGE(1, 3, 0)`, "range rejects a zero step").Expect().ExecError(
			ShouldBeRuntimeError,
			&ExpectedRuntimeError{
				Contains: []string{
					"argument 3 is invalid",
					"Note: argument 3: step must not be zero",
				},
			},
		),
	})
}
