package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func HTTPHelpers(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"GET": httpGet,
		}),
	)
}

func Assertions(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"EXPECT": expect,
		}),
	)
}

func expect(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	if args[0].Compare(args[1]) == 0 {
		return values.EmptyString, nil
	}

	return values.NewString(fmt.Sprintf(`expected "%s", but got "%s"`, args[0], args[1])), nil
}

func httpGet(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	url := args[0].String()

	resp, err := http.Get(url)

	if err != nil {
		return values.None, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return values.None, err
	}

	return values.String(b), nil
}
