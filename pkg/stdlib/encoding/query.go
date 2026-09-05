package encoding

import (
	"context"
	"net/url"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// query_escape escapes query-form data: spaces become + and literal plus signs become %2B.
// @param text {String} The input string.
// @return {String} The converted string.
func QueryEscape(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	return runtime.String(url.QueryEscape(string(text))), nil
}

// query_unescape decodes query-form data, including + as a space. Malformed percent escapes return an error.
// @param text {String} The input string.
// @return {String} The converted string.
func QueryUnescape(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	text, err := runtime.CastArg[runtime.String](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	out, err := url.QueryUnescape(string(text))
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return runtime.String(out), nil
}
