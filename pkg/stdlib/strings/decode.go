package strings

import (
	"context"
	"encoding/base64"
	"net/url"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// FROM_BASE64 returns the value of a base64 representation.
// @param {String} str - The string to decode.
// @return {String} - The decoded string.
func FromBase64(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	value := args[0].String()

	out, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(string(out)), nil
}

// DECODE_URI_COMPONENT returns the decoded String of uri.
// @param {String} uri - Uri to decode.
// @return {String} - Decoded string.
func DecodeURIComponent(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	str, err := url.QueryUnescape(args[0].String())

	if err != nil {
		return values.None, err
	}

	// hack for decoding unicode symbols.
	// eg. convert "\u0026" -> "&""
	str, err = strconv.Unquote("\"" + str + "\"")
	if err != nil {
		return values.None, err
	}

	return values.NewString(str), nil
}
