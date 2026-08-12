package strings

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// from_base64 returns the value of a base64 representation.
// @param str {String} The string to decode.
// @return {String} The decoded string.
func FromBase64(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	value := arg.String()

	out, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return runtime.EmptyString, err
	}

	return runtime.NewString(string(out)), nil
}

// decode_uri_component returns the decoded String of uri.
// @param uri {String} Uri to decode.
// @return {String} Decoded string.
func DecodeURIComponent(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	str, err := url.QueryUnescape(arg.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(str), nil
}
