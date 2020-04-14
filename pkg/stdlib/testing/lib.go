package testing

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

var (
	ErrAssertion = errors.New("assertion error")
)

func RegisterLib(ns core.Namespace) error {
	t := ns.Namespace("T")

	if err := registerNOT(t); err != nil {
		return err
	}

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"ASSERT":  Assert,
			"EMPTY":   Empty,
			"EQUAL":   Equal,
			"FAIL":    Fail,
			"FALSE":   False,
			"INCLUDE": Include,
			"LEN":     Len,
			"MATCH":   Match,
			"NONE":    None,
			"TRUE":    True,
		}),
	)
}

func registerNOT(ns core.Namespace) error {
	t := ns.Namespace("NOT")

	return t.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EMPTY":   NotEmpty,
			"EQUAL":   NotEqual,
			"FALSE":   NotFalse,
			"INCLUDE": NotInclude,
			"MATCH":   NotMatch,
			"NONE":    NotNone,
			"TRUE":    NotTrue,
		}),
	)
}
