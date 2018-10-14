package strings

import (
	"context"
	"encoding/json"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// JSONParse returns a FQL value described by the JSON-encoded input string.
// @params text (String) - The string to parse as JSON.
// @returns FQL value (Read)
func JSONParse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	var val interface{}

	err = json.Unmarshal([]byte(args[0].String()), &val)

	if err != nil {
		return values.EmptyString, err
	}

	return values.Parse(val), nil
}

// JSONStringify returns a JSON string representation of the input value.
// @params value (Read) - The input value to serialize.
// @returns json (String)
func JSONStringify(_ context.Context, args ...core.Value) (core.Value, error) {
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
