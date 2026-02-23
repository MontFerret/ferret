package ferret

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func IsScalar(result Result) bool {
	_, ok := result.(*scalarResult)

	return ok
}

func ForEach(ctx context.Context, result Result, predicate func(val runtime.Value) error) error {
	for {
		val, err := result.Next(ctx)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}
		if err := predicate(val); err != nil {
			return err
		}
	}
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
	val, err := result.Next(ctx)
	if errors.Is(err, io.EOF) {
		return runtime.None, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func JSONStream(ctx context.Context, input io.Writer, result Result) error {
	encoder := encoding.Encoder(json.Default)
	if selected, err := encoding.EncoderFrom(ctx, encoding.ContentTypeJSON); err == nil {
		encoder = selected
	}

	return ForEach(ctx, result, func(val runtime.Value) error {
		j, err := encoder.Encode(val)

		if err != nil {
			return err
		}

		_, err = input.Write(j)

		return err
	})
}
