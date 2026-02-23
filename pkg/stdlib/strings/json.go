package strings

import (
	"context"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// JSON_PARSE returns a value described by the JSON-encoded input string.
// @param {String} str - The string to parse as JSON.
// @return {Any} - Parsed value.
func JSONParse(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	out, err := encodingjson.Default.Decode([]byte(args[0].String()))
	if err != nil {
		return runtime.EmptyString, err
	}

	return out, nil
}

// JSON_STRINGIFY returns a JSON string representation of the input value.
// @param {Any} str - The input value to serialize.
// @return {String} - JSON string.
func JSONStringify(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	out, err := encodingjson.Default.Encode(args[0])

	if err != nil {
		return runtime.EmptyString, err
	}

	return runtime.NewString(string(out)), nil
}
