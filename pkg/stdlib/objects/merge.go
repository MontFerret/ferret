package objects

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Merge(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)
	if err != nil {
		return values.None, err
	}

	var mergedObj *values.Object

	var arr *values.Array
	objs := make([]*values.Object, len(args))
	isArr := false

	if len(args) == 1 {
		err = core.ValidateType(args[0], core.ArrayType)
		if err == nil {
			isArr = true
			arr = args[0].(*values.Array)
			err = validateArray(arr)
		}
	} else {
		for i, arg := range args {
			err = core.ValidateType(arg, core.ObjectType)
			if err != nil {
				break
			}
			objs[i] = arg.(*values.Object)
		}
	}

	if err != nil {
		return values.None, err
	}

	if isArr {
		mergedObj = mergeArray(arr)
	} else {
		mergedObj = mergeObjects(objs...)
	}

	return mergedObj, nil
}

func mergeArray(arr *values.Array) *values.Object {
	return values.NewObject()
}

func mergeObjects(objs ...*values.Object) *values.Object {
	return values.NewObject()
}

func validateArray(arr *values.Array) (err error) {
	var value core.Value

	for i := 0; i < int(arr.Length()); i++ {
		value = arr.Get(values.NewInt(i))
		err = core.ValidateType(value, core.ObjectType)
		if err != nil {
			break
		}
	}
	return
}
