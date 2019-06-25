package strings

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// FromBase64 returns the value of a base64 representation.
// @param base64String (String) - The string to decode.
// @returns value (String) - The decoded string.
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

// DecodeURIComponent returns the decoded String of uri.
// @param (String) - Uri to decode.
// @returns String - Decoded string.
func DecodeURIComponent(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	str, err := url.QueryUnescape(args[0].String())

	if err != nil {
		return values.None, err
	}

	return values.NewString(str), nil
}
