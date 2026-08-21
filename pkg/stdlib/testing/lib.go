package testing

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type assertionRegistration struct {
	name       string
	descriptor assertion
	negatable  bool
}

var assertionCatalog = []assertionRegistration{
	{name: "empty", descriptor: emptyAssertion, negatable: true},
	{name: "eq", descriptor: equalAssertion, negatable: true},
	{name: "fail", descriptor: failAssertion, negatable: false},
	{name: "false", descriptor: falseAssertion, negatable: true},
	{name: "gt", descriptor: gtAssertion, negatable: true},
	{name: "gte", descriptor: gteAssertion, negatable: true},
	{name: "include", descriptor: includeAssertion, negatable: true},
	{name: "len", descriptor: lenAssertion, negatable: true},
	{name: "match", descriptor: matchAssertion, negatable: true},
	{name: "lt", descriptor: ltAssertion, negatable: true},
	{name: "lte", descriptor: lteAssertion, negatable: true},
	{name: "none", descriptor: noneAssertion, negatable: true},
	{name: "true", descriptor: trueAssertion, negatable: true},
	{name: "string", descriptor: stringAssertion, negatable: true},
	{name: "int", descriptor: intAssertion, negatable: true},
	{name: "float", descriptor: floatAssertion, negatable: true},
	{name: "datetime", descriptor: dateTimeAssertion, negatable: true},
	{name: "array", descriptor: arrayAssertion, negatable: true},
	{name: "object", descriptor: objectAssertion, negatable: true},
	{name: "binary", descriptor: binaryAssertion, negatable: true},
}

// @namespace t
func RegisterLib(ns runtime.Namespace) {
	testingNamespace := ns.Namespace("t")
	negativeNamespace := testingNamespace.Namespace("not")

	for _, registration := range assertionCatalog {
		if registration.negatable {
			registerAssertion(negativeNamespace, registration.name, registration.descriptor, registration.descriptor.negative())
		}
	}

	for _, registration := range assertionCatalog {
		registerAssertion(testingNamespace, registration.name, registration.descriptor, registration.descriptor.positive())
	}
}

func registerAssertion(ns runtime.Namespace, name string, descriptor assertion, fn runtime.Function) {
	for arity := descriptor.args.min; arity <= descriptor.args.max; arity++ {
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
