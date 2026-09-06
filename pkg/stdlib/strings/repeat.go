package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const maxRepeatBytes = 64 << 20

// repeat concatenates count copies of text.
// @param text {String} The string to repeat.
// @param count {Int} Non-negative repetition count that fits in an Int on the host. The result must not exceed 64 MiB (67108864 UTF-8 bytes), including when count is one.
// @return {String} The repeated string. Zero count or empty text returns an empty string for a valid count. Oversized results return an argument error that ON ERROR can catch.
func Repeat(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	text, count, err := runtime.CastArgs2[runtime.String, runtime.Int](arg1, arg2)
	if err != nil {
		return runtime.None, err
	}

	maxInt := runtime.Int(int(^uint(0) >> 1))
	if count < 0 || count > maxInt {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "count must be non-negative and fit in an Int on the host"), 1)
	}

	// Divide before multiplying to reject oversized output without overflowing.
	if len(text) > 0 && count > maxRepeatBytes/runtime.Int(len(text)) {
		return runtime.None, runtime.ArgError(runtime.Errorf(runtime.ErrInvalidArgument, "repeated text must not exceed %d bytes (64 MiB)", maxRepeatBytes), 1)
	}

	return runtime.String(strings.Repeat(string(text), int(count))), nil
}
