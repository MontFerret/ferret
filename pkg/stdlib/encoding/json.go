package encoding

import (
	"context"
	"strings"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// json_parse returns a value described by the JSON-encoded input string.
// @param str {String} The string to parse as JSON.
// @return {Any} Parsed value.
func JSONParse(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	out, err := encodingjson.Default.Decode(ctx, strings.NewReader(string(text)))
	if err != nil {
		return runtime.EmptyString, err
	}

	return out, nil
}

// json_stringify returns a JSON string representation of the input value.
// @param str {Any} The input value to serialize.
// @return {String} JSON string.
func JSONStringify(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	out, err := ferretencoding.EncodeBytes(ctx, encodingjson.Default, arg)
	if err != nil {
		return runtime.EmptyString, err
	}

	return runtime.NewString(string(out)), nil
}
