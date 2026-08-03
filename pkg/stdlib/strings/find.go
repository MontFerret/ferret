package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FIND_FIRST returns the position of the first occurrence of the string search inside the string text. Positions start at 0.
// @param {String} str - The source string.
// @param {String} search - The string to seek.
// @param {Int} [start] - Limit the search to a subset of the text, beginning at start.
// @param {Int} [end] - Limit the search to a subset of the text, ending at end
// @return {Int} - The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
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

func findFirst2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return findFirst(ctx, arg1, arg2, runtime.ZeroInt, runtime.Int(len(arg1.String())))
}

func findFirst3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	return findFirst(ctx, arg1, arg2, start, runtime.Int(len(arg1.String())))
}

func findFirst4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	end := runtime.CastOr[runtime.Int](arg4, runtime.Int(len(arg1.String())))
	return findFirst(ctx, arg1, arg2, start, end)
}

func findFirst(_ context.Context, arg1, arg2 runtime.Value, start, end runtime.Int) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)
	search := arg2.String()

	found := strings.Index(string(runes[start:end]), search)

	if found > -1 {
		return runtime.NewInt(found + int(start)), nil
	}

	return runtime.NewInt(found), nil
}

// FIND_LAST returns the position of the last occurrence of the string search inside the string text. Positions start at 0.
// @param {String} src - The source string.
// @param {String} search - The string to seek.
// @param {Int} [start] - Limit the search to a subset of the text, beginning at start.
// @param {Int} [end] - Limit the search to a subset of the text, ending at end
// @return {Int} - The character position of the match. If search is not contained in text, -1 is returned. If search is empty, start is returned.
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

func findLast2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return findLast(ctx, arg1, arg2, runtime.ZeroInt, runtime.Int(len(arg1.String())))
}

func findLast3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	return findLast(ctx, arg1, arg2, start, runtime.Int(len(arg1.String())))
}

func findLast4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	start := runtime.CastOr[runtime.Int](arg3, runtime.ZeroInt)
	end := runtime.CastOr[runtime.Int](arg4, runtime.Int(len(arg1.String())))
	return findLast(ctx, arg1, arg2, start, end)
}

func findLast(_ context.Context, arg1, arg2 runtime.Value, start, end runtime.Int) (runtime.Value, error) {
	text := arg1.String()
	runes := []rune(text)
	search := arg2.String()

	found := strings.LastIndex(string(runes[start:end]), search)

	if found > -1 {
		return runtime.NewInt(found + int(start)), nil
	}

	return runtime.NewInt(found), nil
}
