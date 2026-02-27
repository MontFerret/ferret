package strings

import (
	"context"
	"encoding/base64"
	"net/url"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FROM_BASE64 returns the value of a base64 representation.
// @param {String} str - The string to decode.
// @return {String} - The decoded string.
func FromBase64(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	value := arg.String()

	out, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return runtime.EmptyString, err
	}

	return runtime.NewString(string(out)), nil
}

// DECODE_URI_COMPONENT returns the decoded String of uri.
// @param {String} uri - Uri to decode.
// @return {String} - Decoded string.
func DecodeURIComponent(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	str, err := url.QueryUnescape(arg.String())

	if err != nil {
		return runtime.None, err
	}

	// hack for decoding unicode symbols.
	// eg. convert "\u0026" -> "&""
	str, err = strconv.Unquote("\"" + str + "\"")
	if err != nil {
		return runtime.None, err
	}

	return runtime.NewString(str), nil
}
