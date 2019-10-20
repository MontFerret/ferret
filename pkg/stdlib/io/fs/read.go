package fs

import (
	"context"
	"io/ioutil"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Read reads from a given file.
// @params path (String) - path to file to read from.
// @returns data (Binary) - the read file in binary format.
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
