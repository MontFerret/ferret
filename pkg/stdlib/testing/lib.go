package testing

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// @namespace T
func RegisterLib(ns core.Namespace) error {
	t := ns.Namespace("T")

	if err := registerNOT(t); err != nil {
		return err
	}

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":    base.NewPositiveAssertion(Empty),
			"EQ":       base.NewPositiveAssertion(Equal),
			"FAIL":     base.NewPositiveAssertion(Fail),
			"FALSE":    base.NewPositiveAssertion(False),
			"GT":       base.NewPositiveAssertion(Gt),
			"GTE":      base.NewPositiveAssertion(Gte),
			"INCLUDE":  base.NewPositiveAssertion(Include),
			"LEN":      base.NewPositiveAssertion(Len),
			"MATCH":    base.NewPositiveAssertion(Match),
			"LT":       base.NewPositiveAssertion(Lt),
			"LTE":      base.NewPositiveAssertion(Lte),
			"NONE":     base.NewPositiveAssertion(None),
			"TRUE":     base.NewPositiveAssertion(True),
			"STRING":   base.NewPositiveAssertion(String),
			"INT":      base.NewPositiveAssertion(Int),
			"FLOAT":    base.NewPositiveAssertion(Float),
			"DATETIME": base.NewPositiveAssertion(DateTime),
			"ARRAY":    base.NewPositiveAssertion(Array),
			"OBJECT":   base.NewPositiveAssertion(Object),
			"BINARY":   base.NewPositiveAssertion(Binary),
		}),
	)
}

func registerNOT(ns core.Namespace) error {
	t := ns.Namespace("NOT")

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":    base.NewNegativeAssertion(Empty),
			"EQ":       base.NewNegativeAssertion(Equal),
			"FALSE":    base.NewNegativeAssertion(False),
			"GT":       base.NewNegativeAssertion(Gt),
			"GTE":      base.NewNegativeAssertion(Gte),
			"INCLUDE":  base.NewNegativeAssertion(Include),
			"LEN":      base.NewNegativeAssertion(Len),
			"MATCH":    base.NewNegativeAssertion(Match),
			"LT":       base.NewNegativeAssertion(Lt),
			"LTE":      base.NewNegativeAssertion(Lte),
			"NONE":     base.NewNegativeAssertion(None),
			"TRUE":     base.NewNegativeAssertion(True),
			"STRING":   base.NewNegativeAssertion(String),
			"INT":      base.NewNegativeAssertion(Int),
			"FLOAT":    base.NewNegativeAssertion(Float),
			"DATETIME": base.NewNegativeAssertion(DateTime),
			"ARRAY":    base.NewNegativeAssertion(Array),
			"OBJECT":   base.NewNegativeAssertion(Object),
			"BINARY":   base.NewNegativeAssertion(Binary),
		}),
	)
}
