package strings

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns a FQL value described by the JSON-encoded input string.
 * @params text (String) - The string to parse as JSON.
 * @returns FQL value (Value)
 */
func JsonParse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	var val interface{}

	err = json.Unmarshal([]byte(args[0].String()), &val)

	if err != nil {
		return values.EmptyString, err
	}

	if val == nil {
		return values.None, nil
	}
	return values.Parse(val), nil
}

/*
 * Returns a JSON string representation of the input value.
 * @params value (Value) - The input value to serialize.
 * @returns json (String)
 */
func JsonStringify(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	out, err := json.Marshal(args[0])

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(string(out)), nil
}
