package fs

import (
	"context"
	"io/ioutil"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// READ reads from a given file.
// @param {String} path - Path to file to read from.
// @return {Binary} - File content in binary format.
func Read(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, core.Error(err, "validate arguments number")
	}

	err = core.ValidateType(args[0], types.String)

	if err != nil {
		return values.None, core.Error(err, "validate [0] argument")
	}

	path := args[0].String()

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return values.None, core.Error(err, "read file")
	}

	return values.NewBinary(data), nil
}
