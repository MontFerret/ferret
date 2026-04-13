package vm_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestWarmupClearsStaleHostCacheAcrossEnvironments(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.BaseSpec{
					Input: spec.NewExpressionInput("RETURN F()"),
				},
				Steps: []spec.SequenceStep{
					{
						Name: "with function",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(7), nil
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 7),
					},
					{
						Name: "missing function",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
				},
			},
		}
	})
}

func TestWarmupRebindsWhenEnvironmentFunctionNamesMatch(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F()"),
				Steps: []spec.SequenceStep{
					{
						Name: "env a",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(1), nil
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 1),
					},
					{
						Name: "env b",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(2), nil
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 2),
					},
				},
			},
		}
	})
}

func TestWarmupRebindThrashesAcrossSameNamedHostImplementations(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F()"),
				Steps: []spec.SequenceStep{
					{
						Name: "env a",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(1), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 1),
					},
					{
						Name: "env b",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(2), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 2),
					},
					{
						Name: "env a again",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(1), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 1),
					},
				},
			},
		}
	})
}

func TestWarmupRebindSwitchesBetweenFixedArityAndVarargImplementations(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F(1, 2)"),
				Steps: []spec.SequenceStep{
					{
						Name: "fixed arity",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
									return runtime.NewInt(12), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 12),
					},
					{
						Name: "vararg",
						Env: []vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(102), nil
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 102),
					},
					{
						Name: "fixed arity again",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A2().Add("F", func(context.Context, runtime.Value, runtime.Value) (runtime.Value, error) {
									return runtime.NewInt(22), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 22),
					},
				},
			},
		}
	})
}

func TestWarmupRepeatedMissingRunsAfterSuccessRecoverCleanly(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F()"),
				Steps: []spec.SequenceStep{
					{
						Name: "valid",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(7), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 7),
					},
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
						Name: "recovered",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(9), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 9),
					},
				},
			},
		}
	})
}

func TestWarmupArityMismatchProducesDescriptiveError(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F(1, 2)"),
				Steps: []spec.SequenceStep{
					{
						Name: "arity mismatch",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A1().Add("F", func(_ context.Context, arg runtime.Value) (runtime.Value, error) {
									return arg, nil
								})
							}),
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{
							Message:  "Invalid number of arguments",
							Contains: []string{"expected number of arguments 1, but got 2"},
						}),
					},
				},
			},
		}
	})
}

func TestWarmupArityMismatchRecovery(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN F(1, 2)"),
				Steps: []spec.SequenceStep{
					{
						Name: "correct arity",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A2().Add("F", func(_ context.Context, a, b runtime.Value) (runtime.Value, error) {
									return runtime.NewInt(12), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 12),
					},
					{
						Name: "wrong arity",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A1().Add("F", func(_ context.Context, arg runtime.Value) (runtime.Value, error) {
									return arg, nil
								})
							}),
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Invalid number of arguments"}),
					},
					{
						Name: "recovered",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A2().Add("F", func(_ context.Context, a, b runtime.Value) (runtime.Value, error) {
									return runtime.NewInt(22), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqual, 22),
					},
				},
			},
		}
	})
}

func TestWarmupTrulyMissingFunctionStillReportsUnresolved(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec("RETURN G(1)"),
				Steps: []spec.SequenceStep{
					{
						Name: "truly missing",
						EnvFactory: func() (*vm.Environment, error) {
							return vm.NewDefaultEnvironment(), nil
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
				},
			},
		}
	})
}

func TestWarmupMultiCallsiteRecoveryAfterPartialMissingRun(t *testing.T) {
	RunSequenceFactory(t, func() []spec.Sequence {
		return []spec.Sequence{
			{
				Base: spec.NewBaseSpec(`
LET a = F()
LET b = G()
LET c = F()
RETURN [a, b, c]
`,
				),
				Steps: []spec.SequenceStep{
					{
						Name: "all functions",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(1), nil
								})
								fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(2), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqualJSON, []any{1, 2, 1}),
					},
					{
						Name: "missing g",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(20), nil
								})
							}),
						},
						Error: spec.NewExpectation(ShouldBeRuntimeError, &ExpectedRuntimeError{Message: "Unresolved function"}),
					},
					{
						Name: "recovered",
						Env: []vm.EnvironmentOption{
							vm.WithFunctionsRegistrar(func(fns runtime.FunctionDefs) {
								fns.A0().Add("F", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(30), nil
								})
								fns.A0().Add("G", func(context.Context) (runtime.Value, error) {
									return runtime.NewInt(40), nil
								})
							}),
						},
						Result: spec.NewExpectation(assert.ShouldEqualJSON, []any{30, 40, 30}),
					},
				},
			},
		}
	})
}
