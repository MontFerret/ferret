package strings

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// A limit larger than the host's int range is indistinguishable from its maximum:
// no in-memory string can contain more split results or literal replacements.
func nonNegativeLimit(value runtime.Value, pos int) (int, error) {
	limit, err := runtime.CastArg[runtime.Int](value, pos)
	if err != nil {
		return 0, err
	}

	if limit < 0 {
		return 0, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "limit must be non-negative"), pos)
	}

	return int(min(limit, runtime.Int(int(^uint(0)>>1)))), nil
}
