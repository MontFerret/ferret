package spec

import (
	"context"
	"testing"

	compilerpkg "github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func TestSequenceRunnerRunsMultipleStepsWithDistinctEnvironments(t *testing.T) {
	runner := NewSequenceRunner("Sequence/O0", compilerpkg.WithOptimizationLevel(compilerpkg.O0))

	runner.Run(t, []Sequence{
		{
			Base: NewBaseSpec("RETURN F()", "distinct environments"),
			Steps: []SequenceStep{
				{
					Name: "first env",
					EnvFactory: func() (*vm.Environment, error) {
						return vm.NewEnvironment([]vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(1), nil
							}),
						})
					},
					Result: NewExpectation(assert.ShouldEqual, 1),
				},
				{
					Name: "second env",
					EnvFactory: func() (*vm.Environment, error) {
						return vm.NewEnvironment([]vm.EnvironmentOption{
							vm.WithFunction("F", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								return runtime.NewInt(2), nil
							}),
						})
					},
					Result: NewExpectation(assert.ShouldEqual, 2),
				},
			},
		},
	})
}

func TestSequenceRunnerMatchesExpectedPanics(t *testing.T) {
	runner := NewSequenceRunner("Sequence/O0", compilerpkg.WithOptimizationLevel(compilerpkg.O0))

	runner.Run(t, []Sequence{
		{
			Base: NewBaseSpec("RETURN PANIC_FN()", "expected panic"),
			VM: []vm.Option{
				vm.WithPanicPolicy(vm.PanicPropagate),
			},
			Steps: []SequenceStep{
				{
					Name: "panic",
					EnvFactory: func() (*vm.Environment, error) {
						return vm.NewEnvironment([]vm.EnvironmentOption{
							vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
								panic("boom")
							}),
						})
					},
					Panic: NewExpectation(assert.ShouldEqual, "boom"),
				},
			},
		},
	})
}

func TestSequenceRunnerFailsUnexpectedPanics(t *testing.T) {
	passed := testing.RunTests(
		func(_, _ string) (bool, error) {
			return true, nil
		},
		[]testing.InternalTest{
			{
				Name: "unexpected panic",
				F: func(t *testing.T) {
					runner := NewSequenceRunner("Sequence/O0", compilerpkg.WithOptimizationLevel(compilerpkg.O0))

					runner.Run(t, []Sequence{
						{
							Base: NewBaseSpec("RETURN PANIC_FN()", "unexpected panic"),
							VM: []vm.Option{
								vm.WithPanicPolicy(vm.PanicPropagate),
							},
							Steps: []SequenceStep{
								{
									Name: "panic",
									EnvFactory: func() (*vm.Environment, error) {
										return vm.NewEnvironment([]vm.EnvironmentOption{
											vm.WithFunction("PANIC_FN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
												panic("boom")
											}),
										})
									},
									Result: NewExpectation(assert.ShouldEqual, "unreachable"),
								},
							},
						},
					})
				},
			},
		},
	)

	if passed {
		t.Fatal("expected unexpected panic subtest to fail")
	}
}
