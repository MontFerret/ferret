package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns a substring of value.
 * @params value (String) - The source string.
 * @param offset (Int) - Start at offset, offsets start at position 0.
 * @param length (Int, optional) - At most length characters, omit to get the substring from offset to the end of the string. Optional.
 * @returns substring (String) - A substring of value.
 */
func Substring(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[1], core.IntType)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	offset := int(args[1].(values.Int))
	length := size

	if len(args) > 2 {
		if args[2].Type() == core.IntType {
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

/*
 * Returns the leftmost characters of the string value by index.
 * @param src (String) - The source string.
 * @params length (Int) - The amount of characters to return.
 * @returns substr (String)
 */
func Left(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)

	var pos int

	if args[1].Type() == core.IntType {
		pos = int(args[1].(values.Int))
	}

	if len(text) < pos {
		return values.NewString(text), nil
	}

	return values.NewStringFromRunes(runes[0:pos]), nil
}

/*
 * Returns the rightmost characters of the string value.
 * @param src (String) - The source string.
 * @params length (Int) - The amount of characters to return.
 * @returns substr (String)
 */
func Right(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)
	pos := size

	if args[1].Type() == core.IntType {
		pos = int(args[1].(values.Int))
	}

	if len(text) < pos {
		return values.NewString(string(text)), nil
	}

	return values.NewStringFromRunes(runes[pos:size]), nil
}
