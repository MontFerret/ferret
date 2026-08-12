package strings

import (
	"context"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// json_parse returns a value described by the JSON-encoded input string.
// @param str {String} The string to parse as JSON.
// @return {Any} Parsed value.
func JSONParse(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	out, err := encodingjson.Default.Decode([]byte(arg.String()))
	if err != nil {
		return runtime.EmptyString, err
	}

	return out, nil
}

// json_stringify returns a JSON string representation of the input value.
// @param str {Any} The input value to serialize.
// @return {String} JSON string.
func JSONStringify(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	out, err := encodingjson.Default.Encode(arg)

	if err != nil {
		return runtime.EmptyString, err
	}

	return runtime.NewString(string(out)), nil
}
