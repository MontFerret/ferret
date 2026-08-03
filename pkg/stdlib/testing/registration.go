package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

func registerPositive(ns runtime.Namespace, name string, assertion base.Assertion) {
	registerAssertion(ns, name, assertion, base.NewPositiveAssertion(assertion))
}

func registerNegative(ns runtime.Namespace, name string, assertion base.Assertion) {
	registerAssertion(ns, name, assertion, base.NewNegativeAssertion(assertion))
}

func registerAssertion(ns runtime.Namespace, name string, assertion base.Assertion, fn runtime.Function) {
	for arity := assertion.Args.Min; arity <= assertion.Args.Max; arity++ {
		switch arity {
		case 0:
			ns.Function().A0().Add(name, func(ctx context.Context) (runtime.Value, error) {
				return fn(ctx)
			})
		case 1:
			ns.Function().A1().Add(name, func(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
				return fn(ctx, arg1)
			})
		case 2:
			ns.Function().A2().Add(name, func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
				return fn(ctx, arg1, arg2)
			})
		case 3:
			ns.Function().A3().Add(name, func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
				return fn(ctx, arg1, arg2, arg3)
			})
		case 4:
			ns.Function().A4().Add(name, func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
				return fn(ctx, arg1, arg2, arg3, arg4)
			})
		default:
			panic(fmt.Sprintf("unsupported assertion arity %d for %s", arity, name))
		}
	}
}
