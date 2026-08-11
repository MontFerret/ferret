package strings

import (
	"context"
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/gobwas/glob"
)

var (
	deprecatedLikeSyntax = regexp.MustCompile("[%_]")
)

// LIKE checks whether the pattern search is contained in the string text, using wildcard matching.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
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

// LIKE checks whether the pattern search is contained in the string text, using wildcard matching.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
func like2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return like3(ctx, arg1, arg2, runtime.False)
}

// LIKE checks whether the pattern search is contained in the string text, using wildcard matching.
// @param str {String} The string to search in.
// @param search {String} A search pattern that can contain the wildcard characters.
// @param caseInsensitive {Boolean} If set to true, the matching will be case-insensitive. The default is false.
// @return {Boolean} Returns true if the pattern is contained in text, and false otherwise.
func like3(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	str := arg1.String()
	pattern := arg2.String()

	if len(pattern) == 0 {
		return runtime.NewBoolean(len(str) == 0), nil
	}

	// TODO: Remove me in next releases
	replaced := deprecatedLikeSyntax.ReplaceAllFunc([]byte(pattern), func(b []byte) []byte {
		str := string(b)

		switch str {
		case "%":
			return []byte("*")
		case "_":
			return []byte("?")
		default:
			return b
		}
	})

	pattern = string(replaced)

	if runtime.ToBoolean(arg3) {
		str = strings.ToLower(str)
		pattern = strings.ToLower(pattern)
	}

	g, err := glob.Compile(pattern)

	if err != nil {
		return nil, err
	}

	return runtime.NewBoolean(g.Match(str)), nil
}
