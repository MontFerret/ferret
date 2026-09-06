package encoding

import (
	"context"
	"encoding/base64"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// base64_encode encodes raw bytes using standard padded Base64. Strings supply their UTF-8 bytes.
// @param value {String | Binary} The input data.
// @return {String} Standard padded Base64 text.
func Base64Encode(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString, runtime.TypeBinary); err != nil {
		return runtime.None, err
	}

	return runtime.String(base64.StdEncoding.EncodeToString(runtime.ToBinary(arg))), nil
}

// base64_decode decodes standard padded Base64, preserving arbitrary bytes. Invalid encoding returns an error.
// @param text {String} The encoded data.
// @return {Binary} The decoded bytes.
func Base64Decode(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	out, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return runtime.Binary(out), nil
}
