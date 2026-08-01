package vm_test

import (
	"testing"

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
		S(`RETURN 1s / 2`, "500ms"),
		S(`RETURN 1s / 250ms`, 4),
		S(`RETURN 1s / 3s`, 1.0/3.0),
		S(`RETURN -500ms`, "-500ms"),
		S(`RETURN +500ms`, "500ms"),
		S(`RETURN 5s == 5000ms`, true),
		S(`RETURN 5s > 4999ms`, true),
		S(`RETURN TYPENAME(5s)`, "Duration"),
		S(`RETURN IS_DURATION(5s)`, true),
		S(`RETURN MATCH 5s (5000ms => true, _ => false)`, true),
		S(`RETURN DISTINCT [5s, 5000ms]`, []any{"5s"}),
		spec.NewSpec(`RETURN 1s + 1`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
			Message:  "invalid operation",
			Contains: []string{"operator + is not supported for Duration and Int", ":1:8"},
		}),
		S(`RETURN 0.000001ms * 0.5`, "1ns"),
		S(`RETURN (-0.000001ms) * 0.5`, "-1ns"),
		S(`RETURN 0.000001ms / 2`, "1ns"),
		S(`RETURN (-0.000001ms) / 2`, "-1ns"),
		Error(`RETURN 1s - 1`),
		Error(`RETURN 1s + "1s"`),
		Error(`RETURN "1s" + 1s`),
		Error("RETURN `${1s}`"),
		S("RETURN `${TO_STRING(1s)}`", "1s"),
		Error(`RETURN 1 / 1s`),
		Error(`RETURN 1s * 1s`),
		Error(`RETURN 1s % 1`),
		Error(`RETURN 1s / 0`),
		Error(`RETURN 1s / 0s`),
		Error(`RETURN 9223372036.854775807s + 0.000001ms`),
		Error(`RETURN -((-9223372036.854775807s) - 0.000001ms)`),
	})
}

func TestSchedulingRequiresDuration(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`WAIT(0s) RETURN true`, true),
		S(`WAIT(0.000001ms) RETURN true`, true),
		S(`LET delay = 0s WAIT(delay + 0s) RETURN true`, true),
		S(`RETURN WAITFOR FALSE TIMEOUT 0s EVERY 1ms`, false),
		S(`LET timeout = 0s RETURN WAITFOR FALSE TIMEOUT timeout + 0s EVERY 1ms`, false),
		Error(`WAIT(0) RETURN true`),
		Error(`WAIT(-1ms) RETURN true`),
		Error(`RETURN WAITFOR TRUE TIMEOUT 0`),
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
		Error(`RETURN WAITFOR FALSE TIMEOUT 0s EVERY 0`),
		Error(`RETURN WAITFOR TRUE TIMEOUT 0s EVERY (-1ms)`),
		Error(`LET timeout = -1ms RETURN WAITFOR TRUE TIMEOUT timeout`),
		Error(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY 0 OR RETURN NONE`),
		Error(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY -1ms OR RETURN NONE`),
		Nil(`LET delay = 0ms RETURN T::FAIL() ON ERROR RETRY 1 DELAY delay * 2 OR RETURN NONE`),
		Nil(`RETURN T::FAIL() ON ERROR RETRY 1 DELAY (0ms OR 1ms) OR RETURN NONE`),
	})
}
