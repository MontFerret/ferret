package strings

import (
	"context"
	"strings"

	"github.com/gobwas/glob"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// like matches the whole string using glob syntax: * matches any sequence, ? matches one rune, and character classes and alternatives are supported. Percent and underscore are literal.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Boolean} Whether the whole string matches the pattern.
func Like(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 3)
	if err != nil {
		return runtime.False, err
	}

	if len(args) == 2 {
		return like2(ctx, args[0], args[1])
	}

	return like3(ctx, args[0], args[1], args[2])
}

// like matches the whole string using glob syntax: * matches any sequence, ? matches one rune, and character classes and alternatives are supported. Percent and underscore are literal.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @return {Boolean} Whether the whole string matches the pattern.
func like2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return like3(ctx, arg1, arg2, runtime.False)
}

// like matches the whole string using glob syntax: * matches any sequence, ? matches one rune, and character classes and alternatives are supported. Percent and underscore are literal.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Boolean} Whether the whole string matches the pattern.
func like3(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	text, search, insensitive, err := runtime.CastArgs3[runtime.String, runtime.String, runtime.Boolean](arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	str, pattern := string(text), string(search)
	if insensitive {
		str = strings.ToLower(str)
		pattern = strings.ToLower(pattern)
	}

	g, err := glob.Compile(pattern)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 1)
	}

	return runtime.NewBoolean(g.Match(str)), nil
}
