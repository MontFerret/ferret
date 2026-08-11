package strings

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @param end {Int} Limit the search to a subset of the text, ending before end.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
func FindFirst(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 4)

	if err != nil {
		return runtime.NewInt(-1), err
	}

	switch len(args) {
	case 2:
		return findFirst2(ctx, args[0], args[1])
	case 3:
		return findFirst3(ctx, args[0], args[1], args[2])
	default:
		return findFirst4(ctx, args[0], args[1], args[2], args[3])
	}
}

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
func findFirst2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return findFirst(ctx, arg1, arg2, runtime.ZeroInt, runtime.ZeroInt, false)
}

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
func findFirst3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	return findFirst(ctx, arg1, arg2, start, runtime.ZeroInt, false)
}

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @param end {Int} Limit the search to a subset of the text, ending before end.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
func findFirst4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	end, hasEnd := arg4.(runtime.Int)
	return findFirst(ctx, arg1, arg2, start, end, hasEnd)
}

func findFirst(_ context.Context, arg1, arg2 runtime.Value, start, end runtime.Int, hasEnd bool) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)
	search := arg2.String()
	if !hasEnd {
		end = runtime.Int(len(runes))
	}

	startIndex, endIndex, ok := normalizeFindBounds(len(runes), start, end)
	if !ok {
		return runtime.NewInt(-1), nil
	}

	window := string(runes[startIndex:endIndex])
	found := strings.Index(window, search)

	if found > -1 {
		return runtime.NewInt(startIndex + utf8.RuneCountInString(window[:found])), nil
	}

	return runtime.NewInt(found), nil
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param src {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @param end {Int} Limit the search to a subset of the text, ending before end.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, end is returned.
func FindLast(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 4)

	if err != nil {
		return runtime.NewInt(-1), err
	}

	switch len(args) {
	case 2:
		return findLast2(ctx, args[0], args[1])
	case 3:
		return findLast3(ctx, args[0], args[1], args[2])
	default:
		return findLast4(ctx, args[0], args[1], args[2], args[3])
	}
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param src {String} The source string.
// @param search {String} The string to seek.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, end is returned.
func findLast2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return findLast(ctx, arg1, arg2, runtime.ZeroInt, runtime.ZeroInt, false)
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param src {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, end is returned.
func findLast3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	return findLast(ctx, arg1, arg2, start, runtime.ZeroInt, false)
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param src {String} The source string.
// @param search {String} The string to seek.
// @param start {Int} Limit the search to a subset of the text, beginning at start.
// @param end {Int} Limit the search to a subset of the text, ending before end.
// @return {Int} The character position of the match. If search is not contained in text, -1 is returned. If search is empty, end is returned.
func findLast4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	end, hasEnd := arg4.(runtime.Int)
	return findLast(ctx, arg1, arg2, start, end, hasEnd)
}

func findLast(_ context.Context, arg1, arg2 runtime.Value, start, end runtime.Int, hasEnd bool) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)
	search := arg2.String()
	if !hasEnd {
		end = runtime.Int(len(runes))
	}

	startIndex, endIndex, ok := normalizeFindBounds(len(runes), start, end)
	if !ok {
		return runtime.NewInt(-1), nil
	}

	window := string(runes[startIndex:endIndex])
	found := strings.LastIndex(window, search)

	if found > -1 {
		return runtime.NewInt(startIndex + utf8.RuneCountInString(window[:found])), nil
	}

	return runtime.NewInt(found), nil
}

func normalizeFindBounds(size int, start, end runtime.Int) (int, int, bool) {
	limit := runtime.Int(size)

	if start < 0 {
		start = 0
	} else if start > limit {
		start = limit
	}

	if end < 0 {
		end = 0
	} else if end > limit {
		end = limit
	}

	return int(start), int(end), start <= end
}
