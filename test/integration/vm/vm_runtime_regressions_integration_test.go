package vm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestWarmupHostResolutionPrecedesMissingParamExecution(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec("RETURN MISSING_FN(@foo)").
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
		}
	})
}

func TestStrictWarmupAggregatesMissingParams(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
LET left = @foo
LET right = @bar
RETURN left + right
`).Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected warmup error")
				}

				formatted := diagnostics.Format(err)
				if got, want := strings.Count(formatted, "Missing parameter"), 2; got != want {
					return fmt.Errorf("unexpected missing parameter count: got %d, want %d\n%s", got, want, formatted)
				}

				return nil
			})),
		}
	})
}

func TestStrictWarmupReportsRepeatedMissingParamPerCallsite(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
LET left = @foo
LET right = @foo
RETURN left + right
`).Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected warmup error")
				}

				formatted := diagnostics.Format(err)
				if got, want := strings.Count(formatted, "Missing parameter"), 2; got != want {
					return fmt.Errorf("unexpected missing parameter count: got %d, want %d\n%s", got, want, formatted)
				}

				if got, want := strings.Count(formatted, "missed parameter: @foo"), 2; got != want {
					return fmt.Errorf("unexpected repeated missing parameter cause count: got %d, want %d\n%s", got, want, formatted)
				}

				return nil
			})),
		}
	})
}

func TestStrictWarmupReportsRepeatedUdfMissingParamPerLoadSite(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
FUNC read() => @foo
LET left = read()
LET right = read()
RETURN left + right
`).Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected warmup error")
				}

				formatted := diagnostics.Format(err)
				if got, want := strings.Count(formatted, "Missing parameter"), 1; got != want {
					return fmt.Errorf("unexpected missing parameter count: got %d, want %d\n%s", got, want, formatted)
				}

				for _, needle := range []string{"called from", "VM stack:"} {
					if strings.Contains(formatted, needle) {
						return fmt.Errorf("expected no warmup trace details, got:\n%s", formatted)
					}
				}

				return nil
			})),
		}
	})
}

func TestStrictWarmupFailsProtectedMissingParamUdfCall(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
FUNC risky() => @foo
RETURN risky()?
`).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
				Message:     "Missing parameter",
				NotContains: []string{"called from", "VM stack:"},
			}),
		}
	})
}

func TestStrictWarmupFailsProtectedMissingHostCallForDefaultAndBuiltEnvironment(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN MISSING_FN()?"),
				Steps: []spec.SequenceStep{
					{
						Name: "default",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
					{
						Name: "built",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewEnvironment(nil)
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
				},
			},
		}
	})
}

func TestStrictWarmupFailsOnDeadCodeUnresolvedHostCall(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN false ? MISSING_FN() : 1"),
				Steps: []spec.SequenceStep{
					{
						Name: "default",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
					{
						Name: "dummy",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("DUMMY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.None, nil
							}),
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
				},
			},
		}
	})
}

func TestStrictWarmupAggregatesMissingHostFunctions(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
LET a = MISSING_A()
LET b = MISSING_B()
RETURN a + b
`).Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected warmup error")
				}

				formatted := diagnostics.Format(err)
				if got, want := strings.Count(formatted, "Unresolved function"), 2; got != want {
					return fmt.Errorf("unexpected unresolved function count: got %d, want %d\n%s", got, want, formatted)
				}

				return nil
			})),
		}
	})
}

func TestStrictWarmupReportsRepeatedUnresolvedHostFunctionPerCallsite(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
LET a = MISSING()
LET b = MISSING()
RETURN a + b
`).Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected warmup error")
				}

				formatted := diagnostics.Format(err)
				if got, want := strings.Count(formatted, "Unresolved function"), 2; got != want {
					return fmt.Errorf("unexpected unresolved function count: got %d, want %d\n%s", got, want, formatted)
				}

				return nil
			})),
		}
	})
}

func TestStrictWarmupFailureIsRepeatableUntilEnvironmentFixed(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F()"),
				Steps: []spec.SequenceStep{
					{
						Name: "missing first",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
					{
						Name: "missing second",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
					{
						Name: "fixed env",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(7), nil
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 7),
					},
				},
			},
		}
	})
}

func TestResetDrainsLeakedFramesBetweenFailedRuns(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		counts := make([]int, 0, 2)

		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec(`
FUNC inner() (
	RETURN 1 / 0
)

FUNC outer() (
	RETURN inner()
)

RETURN outer()
`),
				Steps: []spec.SequenceStep{
					{
						Name: "first failure",
						Error: spec.NewExpectation(assert.NewUnaryAssertion(func(actual any) error {
							err, ok := actual.(error)
							if !ok || err == nil {
								return errors.New("expected runtime error")
							}

							var rtErr *vm.RuntimeError
							if !errors.As(err, &rtErr) {
								return fmt.Errorf("expected runtime error, got %T", err)
							}

							if rtErr.Message != "Division by zero" {
								return fmt.Errorf("unexpected runtime error message: got %q, want %q", rtErr.Message, "Division by zero")
							}

							counts = append(counts, strings.Count(rtErr.Format(), "called from"))
							return nil
						})),
					},
					{
						Name: "second failure",
						Error: spec.NewExpectation(assert.NewUnaryAssertion(func(actual any) error {
							err, ok := actual.(error)
							if !ok || err == nil {
								return errors.New("expected runtime error")
							}

							var rtErr *vm.RuntimeError
							if !errors.As(err, &rtErr) {
								return fmt.Errorf("expected runtime error, got %T", err)
							}

							if rtErr.Message != "Division by zero" {
								return fmt.Errorf("unexpected runtime error message: got %q, want %q", rtErr.Message, "Division by zero")
							}

							counts = append(counts, strings.Count(rtErr.Format(), "called from"))
							if len(counts) != 2 {
								return fmt.Errorf("expected two recorded counts, got %d", len(counts))
							}

							if counts[0] != counts[1] {
								return fmt.Errorf("expected stable stack depth across repeated failed runs: first=%d second=%d", counts[0], counts[1])
							}

							return nil
						})),
					},
				},
			},
		}
	})
}

func TestHostNilResultIsNormalizedToNone(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec("RETURN NIL_FN()").
				Env(vm.WithFunction("NIL_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
					return nil, nil
				})).
				Expect().Exec(assert.ShouldBeNil),
		}
	})
}

func TestModuloTypeErrorNotMisclassifiedAsModuloByZero(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`RETURN 5 % "x"`).
				Expect().ExecError(assert.NewUnaryAssertion(func(actual any) error {
				err, ok := actual.(error)
				if !ok || err == nil {
					return errors.New("expected runtime error")
				}

				var rtErr *vm.RuntimeError
				if !errors.As(err, &rtErr) {
					return fmt.Errorf("expected runtime error, got %T", err)
				}

				if rtErr.Message != "Invalid type" {
					return fmt.Errorf("unexpected runtime error message: got %q, want %q", rtErr.Message, "Invalid type")
				}

				if rtErr.Kind != diagnostics.TypeError {
					return fmt.Errorf("unexpected error kind: got %s, want %s", rtErr.Kind, diagnostics.TypeError)
				}

				if strings.Contains(strings.ToLower(rtErr.Format()), "modulo by zero") {
					return fmt.Errorf("expected non-modulo classification, got:\n%s", rtErr.Format())
				}

				return nil
			})),
		}
	})
}

func TestRuntimeErrorIncludesUDFCallStackContext(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
FUNC inner() (
	RETURN @x.foo
)
FUNC middle() (
	LET value = inner()
	RETURN value
)
FUNC outer() (
	LET value = middle()
	RETURN value
)
RETURN outer()
`).
				Env(vm.WithParam("x", runtime.None)).
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
				Message:  `Cannot read property "foo" of None`,
				Contains: []string{"called from inner (#1)", "VM stack: outer -> middle -> inner"},
			}),
		}
	})
}

func TestRuntimeErrorSingleUdfStackFormattingUsesSourceSpelling(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			spec.NewSpec(`
FUNC boo() (
	LET a = 1
	LET b = 0
	RETURN a / b
)
RETURN boo()
`).
				Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{
				Message:  "Division by zero",
				Contains: []string{"called from boo (#1)", "VM stack: boo"},
			}),
		}
	})
}
