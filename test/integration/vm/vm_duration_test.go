package vm_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestNativeDurationValues(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`RETURN 250ms`, "250ms"),
		S(`RETURN 1.5s`, "1.5s"),
		S(`RETURN 1e3ms`, "1s"),
		S(`RETURN 1D`, "24h0m0s"),
		S(`RETURN [1MS, 1s, 1M, 1h, 1d]`, []any{"1ms", "1s", "1m0s", "1h0m0s", "24h0m0s"}),
		S(`RETURN { duration: 2.5E-1S }.duration`, "250ms"),
		S(`RETURN 1s + 500ms`, "1.5s"),
		S(`RETURN 2s - 500ms`, "1.5s"),
		S(`RETURN 500ms * 3`, "1.5s"),
		S(`RETURN 3 * 500ms`, "1.5s"),
		S(`RETURN 5s * TO_NUMBER("2")`, "10s"),
		S(`RETURN TO_NUMBER("2") * 5s`, "10s"),
		S(`RETURN 5s * TO_NUMBER("2.5")`, "12.5s"),
		S(`RETURN TO_NUMBER("2.5") * 5s`, "12.5s"),
		S(`RETURN 1s / 2`, "500ms"),
		S(`RETURN 1s / 250ms`, 4),
		S(`RETURN 1s / 3s`, 1.0/3.0),
		S(`RETURN -500ms`, "-500ms"),
		S(`RETURN +500ms`, "500ms"),
		S(`RETURN 5s == 5000ms`, true),
		S(`RETURN 5s > 4999ms`, true),
		S(`RETURN TYPENAME(5s)`, "Duration"),
		S(`RETURN IS_DURATION(5s)`, true),
		S(`RETURN TO_DURATION(5s)`, "5s"),
		S(`RETURN TO_DURATION("250ms")`, "250ms"),
		S(`RETURN TO_DURATION("1h30m")`, "1h30m0s"),
		S(`RETURN TO_DURATION(5000)`, "5s"),
		S(`RETURN TO_DURATION(1.5)`, "1.5ms"),
		S(`RETURN [TO_DURATION(NONE), TO_DURATION(false), TO_DURATION(true)]`, []any{"0s", "0s", "1ms"}),
		S(`RETURN TO_DURATION([[2]])`, "2ms"),
		S(`RETURN TYPENAME(TO_DURATION("500ms"))`, "Duration"),
		S(`RETURN TO_DURATION("500ms") + 500ms`, "1s"),
		S(`RETURN 1s + TO_DURATION(1)`, "1.001s"),
		S(`RETURN TO_DURATION(1) + 1s`, "1.001s"),
		S(`RETURN 1s - TO_DURATION(1)`, "999ms"),
		S(`RETURN 1s + TO_DURATION("1s")`, "2s"),
		S(`RETURN TO_DURATION("1s") + 1s`, "2s"),
		S(`RETURN 1s / TO_NUMBER("2")`, "500ms"),
		S(`RETURN 1s / TO_DURATION("250ms")`, 4),
		S(`RETURN 1s == "1s"`, false),
		S(`RETURN 1s != "1s"`, true),
		S(`RETURN 1s == TO_DURATION("1s")`, true),
		S(`RETURN 1s != "tomorrow"`, true),
		S(`RETURN 1s == "tomorrow"`, false),
		Error(`RETURN 1s > 999`),
		Error(`RETURN [1s, 2s] ALL > 999`),
		S(`RETURN [1s, 2s] ANY == "2s"`, false),
		S(`RETURN [1s, 2s] ANY == TO_DURATION("2s")`, true),
		S(`RETURN MATCH 5s {5000ms => true, _ => false}`, true),
		S(`RETURN MATCH 1s {"1s" => true, _ => false}`, false),
		S(`RETURN [1s] == ["1s"]`, false),
		S(`RETURN { value: 1s } == { value: "1s" }`, false),
		S(`RETURN "1s" IN [1s]`, false),
		S(`RETURN DISTINCT [5s, 5000ms]`, []any{"5s"}),
		S(`RETURN DISTINCT ["1s", 1000, 1s]`, []any{"1s", float64(1000), "1s"}),
		S(`RETURN UNION_DISTINCT(["1s", 1000], [1s, 1000ms])`, []any{"1s", float64(1000), "1s"}),
		S(`FOR value IN [1s, 1000ms] COLLECT key = value WITH COUNT INTO count RETURN { key, count }`, []any{
			map[string]any{"key": "1s", "count": float64(2)},
		}),
		Error(`RETURN SORTED([1s, "2s"])`),
		Error(`RETURN [1s] < ["1s"]`),
		S(`RETURN 0.000001ms * 0.5`, "0s"),
		S(`RETURN (-0.000001ms) * 0.5`, "0s"),
		S(`RETURN 0.000001ms / 2`, "0s"),
		S(`RETURN (-0.000001ms) / 2`, "0s"),
		Error(`RETURN TO_DURATION("1fortnight")`),
		Error(`RETURN TO_DURATION([1, 2])`),
		Error(`RETURN TO_DURATION({ value: 1 })`),
		spec.NewSpec(`RETURN "tomorrow" > 1s`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '>' cannot be applied to String and Duration", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1s <= "tomorrow"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '<=' cannot be applied to Duration and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN "tomorrow" < 1s`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '<' cannot be applied to String and Duration", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1s >= "tomorrow"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '>=' cannot be applied to Duration and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN [1s] ANY < "tomorrow"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '<' cannot be applied to Duration and String", ":1:8"},
		}),
		S("RETURN `${1s}`", "1s"),
		S("RETURN `${TO_STRING(1s)}`", "1s"),
		Error(`RETURN 1 / 1s`),
		spec.NewSpec(`RETURN 5s * "invalid"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '*' cannot be applied to Duration and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN "invalid" * 5s`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '*' cannot be applied to String and Duration", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1s + 1`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '+' cannot be applied to Duration and Int", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1 + 1s`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '+' cannot be applied to Int and Duration", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1s - 1`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to Duration and Int", ":1:8"},
		}),
		S(`RETURN 1s + "1s"`, "1s1s"),
		S(`RETURN "1s" + 1s`, "1s1s"),
		spec.NewSpec(`RETURN 1s / "2"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '/' cannot be applied to Duration and String", ":1:8"},
		}),
		spec.NewSpec(`RETURN 1s / "250ms"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '/' cannot be applied to Duration and String", ":1:8"},
		}),
		Error(`RETURN 1s * "2s"`),
		Error(`RETURN 1s * 1s`),
		Error(`RETURN 1s % 1`),
		Error(`RETURN 1s / 0`),
		Error(`RETURN 1s / 0s`),
		Error(`RETURN 9223372036.854775807s + 0.000001ms`),
		Error(`RETURN -((-9223372036.854775807s) - 0.000001ms)`),
	})
}

func TestDateTimeOperators(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z")`, "2026-08-01T12:00:00Z"),
		S(`RETURN TO_DATETIME(TO_DATETIME("2026-08-01T12:00:00Z"))`, "2026-08-01T12:00:00Z"),
		S(`RETURN TO_DATETIME(1690992000, "s")`, "2023-08-02T16:00:00Z"),
		S(`RETURN TO_DATETIME(1690992000000, "ms")`, "2023-08-02T16:00:00Z"),
		S(`RETURN TO_DATETIME(1690992000000000, "us")`, "2023-08-02T16:00:00Z"),
		S(`RETURN TO_DATETIME(1690992000000000000, "ns")`, "2023-08-02T16:00:00Z"),
		S(`RETURN TO_DATETIME(1, "SECONDS")`, "1970-01-01T00:00:01Z"),
		S(`RETURN TO_DATETIME(1.5, "s")`, "1970-01-01T00:00:01.5Z"),
		S(`RETURN TO_DATETIME(1.5, "ms")`, "1970-01-01T00:00:00.0015Z"),
		S(`RETURN TO_DATETIME(-1, "s")`, "1969-12-31T23:59:59Z"),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") + 30m`, "2026-08-01T12:30:00Z"),
		S(`RETURN 30m + TO_DATETIME("2026-08-01T12:00:00Z")`, "2026-08-01T12:30:00Z"),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") - 30m`, "2026-08-01T11:30:00Z"),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") - TO_DATETIME("2026-08-01T11:59:30Z")`, "30s"),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00+02:00") - TO_DATETIME("2026-08-01T10:00:00Z")`, "0s"),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00+02:00") == TO_DATETIME("2026-08-01T10:00:00Z")`, true),
		S(`RETURN TYPENAME(NOW() + TO_DURATION("5m"))`, "DateTime"),
		S(`RETURN @delay + TO_DURATION("500ms")`, "5.5s").Env(spec.WithParam("delay", 5*time.Second)),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") == "2026-08-01T12:00:00Z"`, false),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") == "2026-08-01T13:00:00Z"`, false),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") != "2026-08-01T12:00:00Z"`, true),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") != "2026-08-01T13:00:00Z"`, true),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") > "2026-08-01T12:00:00Z"`, true),
		S(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") > "not-a-date"`, true),
		spec.NewSpec(`RETURN TO_DATETIME(0)`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"numeric DateTime conversion requires an explicit epoch unit", `"s", "ms", "us", or "ns"`, ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(0.5)`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"numeric DateTime conversion requires an explicit epoch unit", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME("1", "s")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"String", "epoch units are only valid for Int or Float", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME("2026-08-01T12:00:00Z", "ms")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"String", "epoch units are only valid for Int or Float", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(TO_DATETIME("2026-08-01T12:00:00Z"), "s")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"DateTime", "epoch units are only valid for Int or Float", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(1, "minutes")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"unsupported epoch unit", `"s", "ms", "us", or "ns"`, ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(1, NONE)`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid type",
			Contains: []string{"DateTime epoch unit", "expected String", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(1, 1)`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid type",
			Contains: []string{"DateTime epoch unit", "expected String", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(@value, "s")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"must be finite", ":1:8"},
		}).Env(spec.WithParam("value", math.NaN())),
		spec.NewSpec(`RETURN TO_DATETIME(@value, "s")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid argument",
			Contains: []string{"must be finite", ":1:8"},
		}).Env(spec.WithParam("value", math.Inf(1))),
		spec.NewSpec(`RETURN TO_DATETIME(9223372036854775807, "s")`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Contains: []string{"out of range", "supported DateTime range", ":1:8"},
		}),
		spec.NewSpec(`RETURN TO_DATETIME()`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "invalid number of arguments",
			Contains: []string{
				"wrong number of arguments in call to TO_DATETIME",
				"Note: TO_DATETIME expects 1 or 2 arguments, but got 0",
				"Hint: Pass 1 or 2 arguments to TO_DATETIME",
				":1:8",
			},
		}),
		spec.NewSpec(`RETURN TO_DATETIME(1, "s", true)`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "invalid number of arguments",
			Contains: []string{
				"wrong number of arguments in call to TO_DATETIME",
				"Note: TO_DATETIME expects 1 or 2 arguments, but got 3",
				"Hint: Pass 1 or 2 arguments to TO_DATETIME",
				":1:8",
			},
		}),
		Error(`RETURN TO_DATETIME("invalid")`),
		spec.NewSpec(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") - "tomorrow"`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '-' cannot be applied to DateTime and String", ":1:8"},
		}),
		S(`RETURN TYPENAME(NOW() + "5m")`, "String"),
		spec.NewSpec(`RETURN NOW() + 5000`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator '+' cannot be applied to DateTime and Int", ":1:8"},
		}),
		Error(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") + TO_DATETIME("2026-08-01T12:00:00Z")`),
		Error(`RETURN 1s - TO_DATETIME("2026-08-01T12:00:00Z")`),
		Error(`RETURN TO_DATETIME("2026-08-01T12:00:00Z") * 2`),
	})
}

func TestTemporalComparisonDoesNotInspectOpaqueHostValues(t *testing.T) {
	lengthErr := errors.New("duration list length failed")
	atErr := runtime.Error(runtime.ErrRange, "duration list item failed")

	RunSpecs(t, []spec.Spec{
		S(`RETURN 1s == @value`, false).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		S(`RETURN 1s != @value`, true).
			Env(spec.WithParam("value", newFallibleDurationList(nil, atErr, runtime.NewInt(1)))),
		S(`RETURN [1s] ANY == @value`, false).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		S(`RETURN @left == @right ? 10 : 20`, 20).
			Env(
				spec.WithParam("left", time.Second),
				spec.WithParam("right", newFallibleDurationList(lengthErr, nil)),
			),
		S(`RETURN @left != @right ? 10 : 20`, 10).
			Env(
				spec.WithParam("left", time.Second),
				spec.WithParam("right", newFallibleDurationList(lengthErr, nil)),
			),
		S(`RETURN @value == 1s ? 10 : 20`, 20).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		S(`RETURN @value != 1s ? 10 : 20`, 10).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		S(`RETURN MATCH @value {1s => 10, _ => 20}`, 20).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		Error(`RETURN 1s < @value`).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		spec.NewSpec(`RETURN TO_DURATION(@value)`).
			Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Contains: []string{"duration list length failed", ":1:8"},
		}).
			Env(spec.WithParam("value", newFallibleDurationList(lengthErr, nil))),
		spec.NewSpec(`RETURN TO_DURATION(@value)`).
			Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Contains: []string{"duration list item failed", ":1:8"},
		}).
			Env(spec.WithParam("value", newFallibleDurationList(nil, atErr, runtime.NewInt(1)))),
	})
}

func TestSchedulingRequiresDuration(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`WAIT(0s) RETURN true`, true),
		S(`WAIT(0) RETURN true`, true),
		S(`WAIT("0s") RETURN true`, true),
		S(`WAIT(0.000001ms) RETURN true`, true),
		S(`LET delay = 0s WAIT(delay + 0s) RETURN true`, true),
		S(`RETURN WAITFOR FALSE TIMEOUT 0s EVERY 1ms`, false),
		S(`RETURN WAITFOR FALSE TIMEOUT 0 EVERY 0`, false),
		S(`RETURN WAITFOR FALSE TIMEOUT "0s" EVERY 0, "0ms"`, false),
		S(`LET timeout = 0s RETURN WAITFOR FALSE TIMEOUT timeout + 0s EVERY 1ms`, false),
		S(`LET timeout = "0s" RETURN WAITFOR FALSE TIMEOUT timeout`, false),
		Error(`WAIT(-1ms) RETURN true`),
		Error(`WAIT(-1) RETURN true`),
		S(`RETURN WAITFOR TRUE TIMEOUT 0`, true),
		spec.NewSpec(
			`RETURN WAITFOR TRUE TIMEOUT (-1ms)`,
			"negative_timeout.fql",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message: "invalid operation",
			Contains: []string{
				":1:29",
				"wait duration must not be negative",
			},
		}),
		Error(`RETURN WAITFOR TRUE TIMEOUT 0s EVERY (-1ms)`),
		Error(`LET timeout = -1ms RETURN WAITFOR TRUE TIMEOUT timeout`),
		Nil(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY 0 OR RETURN NONE`),
		Nil(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY "0s" OR RETURN NONE`),
		Nil(`LET delay = "0s" RETURN T::FAIL() ON ERROR RETRY 1 DELAY delay OR RETURN NONE`),
		Error(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY -1ms OR RETURN NONE`),
		Nil(`LET delay = 0ms RETURN T::FAIL() ON ERROR RETRY 1 DELAY delay * 2 OR RETURN NONE`),
		Nil(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY (0ms OR 1ms) OR RETURN NONE`),
	})
}
