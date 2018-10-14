package strings

import (
	"context"
	"encoding/base64"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

/*
 * Returns the value of a base64 representation.
 * @param base64String (String) - The string to decode.
 * @returns value (String) - The decoded string.
 */
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
