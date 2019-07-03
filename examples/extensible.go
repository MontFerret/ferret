package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func main() {
	strs, err := getStrings()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, str := range strs {
		fmt.Println(str)
	}
}

func getStrings() ([]string, error) {
	// function implements is a type of a function that ferret supports as a runtime function
	transform := func(ctx context.Context, args ...core.Value) (core.Value, error) {
		// it's just a helper function which helps to validate a number of passed args
		err := core.ValidateArgs(args, 1, 1)

		if err != nil {
			// it's recommended to return built-in None type, instead of nil
			return values.None, err
		}

		// this is another helper functions allowing to do type validation
		err = core.ValidateType(args[0], types.String)

		if err != nil {
			return values.None, err
		}

		// cast to built-in string type
		str := args[0].(values.String)

		return values.NewString(strings.ToUpper(str.String() + "_ferret")), nil
	}

	query := `
		FOR el IN ["foo", "bar", "qaz"]
			// conventionally all functions are registered in upper case
			RETURN TRANSFORM(el)
	`

	comp := compiler.New()

	if err := comp.RegisterFunction("transform", transform); err != nil {
		return nil, err
	}

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	out, err := program.Run(context.Background())

	if err != nil {
		return nil, err
	}

	res := make([]string, 0, 3)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
