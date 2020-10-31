package fs

import (
	"context"
	"os"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WRITE writes the given data into the file.
// @param {String} path - File path to write into.
// @param {Binary} data - Data to write.
// @param {Object} [params] - additional parameters:
// @param {String} [params.mode] - Write mode.
// * x - Exclusive: returns an error if the file exist. It can be combined with other modes
// * a - Append: will create a file if the specified file does not exist
// * w - Write (Default): will create a file if the specified file does not exist
func Write(_ context.Context, args ...core.Value) (core.Value, error) {
	err := validateRequiredWriteArgs(args)

	if err != nil {
		return values.None, err
	}

	fpath := args[0].String()
	arg1 := args[1]

	var data []byte

	if arg1.Type() == types.Binary {
		data = arg1.(values.Binary)
	} else {
		data = []byte(arg1.String())
	}

	params := defaultParams

	if len(args) == 3 {
		params, err = parseParams(args[2])

		if err != nil {
			return values.None, core.Error(
				err,
				"parse `params` argument",
			)
		}
	}

	// 0666 - read & write
	file, err := os.OpenFile(fpath, params.ModeFlag, 0666)

	if err != nil {
		return values.None, core.Error(err, "open file")
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return values.None, core.Error(err, "write file")
	}

	return values.None, nil
}

func validateRequiredWriteArgs(args []core.Value) error {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return err
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return err
	}

	return nil
}

// parsedParams contains parsed additional parameters.
type parsedParams struct {
	ModeFlag int
}

var defaultParams = parsedParams{
	// the same as `w`
	ModeFlag: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
}

func parseParams(value core.Value) (parsedParams, error) {
	err := core.ValidateType(value, types.Object)

	if err != nil {
		return parsedParams{}, err
	}

	obj := value.(*values.Object)

	params := defaultParams

	modestr, exists := obj.Get(values.NewString("mode"))

	if exists {
		flag, err := parseWriteMode(modestr.String())

		if err != nil {
			return parsedParams{}, core.Error(
				core.ErrInvalidArgument,
				"parse write mode",
			)
		}

		params.ModeFlag = flag
	}

	return params, nil
}

func parseWriteMode(s string) (int, error) {
	letters := []rune(s)
	count := len(letters)

	if count == 0 || count > 2 {
		return -1, core.Errorf(
			core.ErrInvalidArgument,
			"must be from 1 to 2 mode letters, got `%d`", count,
		)
	}

	// sort letters for more convenient work with it
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })

	// minimum flag for writing to file
	flag := os.O_WRONLY | os.O_CREATE

	if count == 2 {
		// since letter is sorted, `x` will always be the letters[1]
		if letters[1] != 'x' {
			return -1, core.Errorf(
				core.ErrInvalidArgument,
				"invalid mode `%s`", s,
			)
		}

		flag |= os.O_EXCL
	}

	switch letters[0] {
	case 'a':
		flag |= os.O_APPEND
	case 'w':
		flag |= os.O_TRUNC
	default:
		return -1, core.Errorf(
			core.ErrInvalidArgument,
			"invalid mode `%s`", s,
		)
	}

	return flag, nil
}
