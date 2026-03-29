package fs

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// READ reads from a given file.
// @param {String} path - Path to file to read from.
// @return {Binary} - File content in binary format.
func Read(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	path, err := runtime.CastArg[runtime.String](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	reader, err := fs.ReaderFrom(ctx)

	if err != nil {
		return runtime.None, err
	}

	data, err := reader.ReadFile(path.String())

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewBinary(data), nil
}
