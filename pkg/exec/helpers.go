package exec

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func IsScalar(result Result) bool {
	_, ok := result.(*scalarResult)

	return ok
}

func ForEach(ctx context.Context, result Result, predicate func(val runtime.Value) error) error {
	hasNext, err := result.HasNext(ctx)

	if err != nil {
		return err
	}

	for hasNext {
		val, err := result.Next(ctx)

		if err != nil {
			return err
		}

		if err := predicate(val); err != nil {
			return err
		}

		hasNext, err = result.HasNext(ctx)

		if err != nil {
			return err
		}
	}

	return nil
}

func Collect(ctx context.Context, result Result) ([]runtime.Value, error) {
	var res []runtime.Value

	err := ForEach(ctx, result, func(val runtime.Value) error {
		res = append(res, val)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func One(ctx context.Context, result Result) (runtime.Value, error) {
	hasNext, err := result.HasNext(ctx)

	if err != nil {
		return nil, err
	}

	if !hasNext {
		return runtime.None, nil
	}

	return result.Next(ctx)
}

func JSONStream(ctx context.Context, input io.Writer, result Result) error {
	return ForEach(ctx, result, func(val runtime.Value) error {
		j, err := val.MarshalJSON()

		if err != nil {
			return err
		}

		_, err = input.Write(j)

		return err
	})
}
