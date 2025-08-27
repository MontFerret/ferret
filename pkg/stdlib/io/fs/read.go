package fs

import (
	"context"
	"os"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// READ reads from a given file.
// @param {String} path - Path to file to read from.
// @return {Binary} - File content in binary format.
func Read(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.None, runtime.Error(err, "validate arguments number")
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.None, runtime.Error(err, "validate [0] argument")
	}

	path := args[0].String()
	data, err := os.ReadFile(path)

	if err != nil {
		return runtime.None, runtime.Error(err, "read file")
	}

	return runtime.NewBinary(data), nil
}
