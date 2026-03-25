package fs

import (
	"context"
	"os"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// READ reads from a given file.
// @param {String} path - Path to file to read from.
// @return {binary} - File content in binary format.
func Read(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	path, err := runtime.CastArg[runtime.String](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	data, err := os.ReadFile(path.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewBinary(data), nil
}
