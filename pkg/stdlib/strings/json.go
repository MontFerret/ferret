package strings

import (
	"context"
	"encoding/json"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// JSON_PARSE returns a value described by the JSON-encoded input string.
// @param {String} str - The string to parse as JSON.
// @return {Any} - Parsed value.
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

// JSON_STRINGIFY returns a JSON string representation of the input value.
// @param {Any} str - The input value to serialize.
// @return {String} - JSON string.
func JSONStringify(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	out, err := jettison.MarshalOpts(args[0])

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(string(out)), nil
}
