package image

import (
	"context"
	"os"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"
)

var (
	binaryType = []core.Type{types.Binary}
	stringType = []core.Type{types.String}
)

func SaveFile(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)
	if err != nil {
		return values.None, err
	}

	pairs := []core.PairValueType{
		{Value: args[0], Types: binaryType},
		{Value: args[1], Types: stringType},
	}

	err = core.ValidateValueTypePairs(pairs...)
	if err != nil {
		return values.None, err
	}

	binary := args[0].(values.Binary).Unwrap().([]byte)
	filename := args[1].(values.String).String()

	_, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(err) {
		return values.None, errors.Errorf(`file "%s" already exists`, filename)
	}

	file, err := os.Create(filename)
	if err != nil {
		return values.None, err
	}

	_, err = file.Write(binary)
	if err != nil {
		return values.None, err
	}

	return values.None, nil
}
