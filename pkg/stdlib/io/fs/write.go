package fs

import (
	"context"
	"os"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// write writes the given data into the file.
// @param path {String} File path to write into.
// @param data {Binary} Data to write.
// @param params {Map} Additional parameters. The mode field selects the write mode:
// * x - Exclusive: returns an error if the file exist. It can be combined with other modes
// * a - Append: will create a file if the specified file does not exist
// * w - Write (Default): will create a file if the specified file does not exist
// @return {None} None.
func Write(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return write2(ctx, args[0], args[1])
	}

	return write3(ctx, args[0], args[1], args[2])
}

// write writes the given data into the file.
// @param path {String} File path to write into.
// @param data {Binary} Data to write.
// @return {None} None.
func write2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return writeFile(ctx, arg1, arg2, defaultParams)
}

// write writes the given data into the file.
// @param path {String} File path to write into.
// @param data {Binary} Data to write.
// @param params {Map} Additional parameters. The mode field selects the write mode:
// * x - Exclusive: returns an error if the file exist. It can be combined with other modes
// * a - Append: will create a file if the specified file does not exist
// * w - Write (Default): will create a file if the specified file does not exist
// @return {None} None.
func write3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	params, err := parseParams(arg3, 2)
	if err != nil {
		return runtime.None, runtime.Error(
			err,
			"parse `params` argument",
		)
	}

	return writeFile(ctx, arg1, arg2, params)
}

func writeFile(ctx context.Context, arg1, arg2 runtime.Value, params parsedParams) (runtime.Value, error) {
	fpath, err := runtime.CastArg[runtime.String](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	data, err := runtime.CastArg[runtime.Binary](arg2, 1)
	if err != nil {
		return runtime.None, err
	}

	filesystem, err := fs.FileSystemFrom(ctx)

	if err != nil {
		return runtime.None, err
	}

	// 0666 - read & write
	file, err := filesystem.OpenFile(string(fpath), params.ModeFlag, 0666)

	if err != nil {
		return runtime.None, runtime.Error(err, "open file")
	}

	defer file.Close()

	_, err = file.Write(data)

	if err != nil {
		return runtime.None, runtime.Error(err, "write file")
	}

	return runtime.None, nil
}

// parsedParams contains parsed additional parameters.
type parsedParams struct {
	ModeFlag int
}

var defaultParams = parsedParams{
	// the same as `w`
	ModeFlag: os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
}

func parseParams(value runtime.Value, pos int) (parsedParams, error) {
	err := runtime.ValidateArgType(value, pos, runtime.TypeObject, runtime.TypeMap)

	if err != nil {
		return parsedParams{}, err
	}

	obj := value.(runtime.Map)

	params := defaultParams

	modestr, err := obj.Get(context.Background(), runtime.NewString("mode"))

	if err == nil {
		flag, err := parseWriteMode(modestr.String())

		if err != nil {
			return parsedParams{}, runtime.Error(
				runtime.ErrInvalidArgument,
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
		return -1, runtime.Errorf(
			runtime.ErrInvalidArgument,
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
			return -1, runtime.Errorf(
				runtime.ErrInvalidArgument,
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
		return -1, runtime.Errorf(
			runtime.ErrInvalidArgument,
			"invalid mode `%s`", s,
		)
	}

	return flag, nil
}
