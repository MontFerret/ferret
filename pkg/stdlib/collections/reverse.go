package collections

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Reverse reverses Unicode code points in a string or indexed elements in a list.
// Lists construct an independent outer list through New, preserving implementation
// family and backend configuration. Elements remain shallow borrowed references.
// Failed or canceled construction is handled by the factory; subsequent failures
// close the incomplete destination and preserve cleanup errors. Success transfers
// destination ownership to the caller. The source and its elements are not closed.
// @param value {String | List} The string or indexed list to reverse.
// @return {String | List} A reversed string or new list of the source's implementation family.
func Reverse(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	switch col := arg.(type) {
	case runtime.String:
		runes := []rune(string(col))
		size := len(runes)

		// Reverse
		for i := 0; i < size/2; i++ {
			runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
		}

		return runtime.NewString(string(runes)), nil
	case runtime.List:
		return reverseList(ctx, col)
	default:
		return runtime.None, runtime.TypeErrorOf(arg, runtime.TypeList, runtime.TypeString)
	}
}

func reverseList(ctx context.Context, col runtime.List) (_ runtime.Value, err error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	size, err := col.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if size < 0 {
		return runtime.None, runtime.Error(runtime.ErrInvalidOperation, "negative list length")
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	result, err := col.New(ctx)
	if err != nil {
		return runtime.None, err
	}

	defer func() {
		if err != nil {
			if closer, ok := result.(io.Closer); ok {
				if closeErr := closer.Close(); closeErr != nil {
					err = errors.Join(err, closeErr)
				}
			}
		}
	}()

	for i := size; i > 0; i-- {
		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		item, err := col.At(ctx, i-1)
		if err != nil {
			return runtime.None, err
		}

		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		if err := result.Append(ctx, item); err != nil {
			return runtime.None, err
		}
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return result, nil
}
