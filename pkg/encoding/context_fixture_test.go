package encoding_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type contextList struct {
	runtime.List
	expected context.Context
	seen     bool
}

func (v *contextList) Length(ctx context.Context) (runtime.Int, error) {
	if ctx != v.expected {
		return 0, errors.New("Length lost caller context")
	}

	v.seen = true

	return v.List.Length(ctx)
}

func (v *contextList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if ctx != v.expected {
		return runtime.None, errors.New("At lost caller context")
	}

	return v.List.At(ctx, index)
}
