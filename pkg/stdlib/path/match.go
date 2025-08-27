package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MATCH reports whether name matches the pattern.
// @param {String} pattern - The pattern.
// @param {String} name - The name.
// @return {Boolean} - True if the name matches the pattern.
func Match(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 2)

	if err != nil {
		return runtime.False, err
	}

	err = runtime.ValidateType(args[0], runtime.TypeString)

	if err != nil {
		return runtime.False, err
	}

	err = runtime.ValidateType(args[1], runtime.TypeString)

	if err != nil {
		return runtime.False, err
	}

	pattern := args[0].String()
	name := args[1].String()

	matched, err := path.Match(pattern, name)

	if err != nil {
		return runtime.False, runtime.Error(err, "match")
	}

	return runtime.NewBoolean(matched), nil
}
