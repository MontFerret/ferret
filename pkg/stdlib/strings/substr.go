package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SUBSTRING returns a substring of value.
// @param {String} str - The source string.
// @param {Int} offset - Start at offset, offsets start at position 0.
// @param {Int} [length] - At most length characters, omit to get the substring from offset to the end of the string.
// @return {String} - A substring of value.
func Substring(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[1], types.Int)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	offset := int(args[1].(values.Int))
	length := size

	if len(args) > 2 {
		if args[2].Type() == types.Int {
			length = int(args[2].(values.Int))
		}
	}

	var substr []rune

	if length == size {
		substr = runes[offset:]
	} else {
		end := offset + length

		if size > end {
			substr = runes[offset:end]
		} else {
			substr = runes[offset:]
		}
	}

	return values.NewStringFromRunes(substr), nil
}

// LEFT returns the leftmost characters of the string value by index.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The leftmost characters of the string value by index.
func Left(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)

	var pos int

	if args[1].Type() == types.Int {
		pos = int(args[1].(values.Int))
	}

	if len(text) < pos {
		return values.NewString(text), nil
	}

	return values.NewStringFromRunes(runes[0:pos]), nil
}

// RIGHT returns the rightmost characters of the string value.
// @param {String} str - The source string.
// @param {Int} length - The amount of characters to return.
// @return {String} - The rightmost characters of the string value.
func Right(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	pos := size

	if args[1].Type() == types.Int {
		pos = int(args[1].(values.Int))
	}

	if len(text) < pos {
		return values.NewString(text), nil
	}

	return values.NewStringFromRunes(runes[size-pos : size]), nil
}
