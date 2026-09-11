package collections_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	listBackend struct {
		encoding string
		limit    int
	}

	configuredList struct {
		expected      context.Context
		closeErr      error
		newCleanupErr error
		*runtime.Array
		backend      *listBackend
		length       *runtime.Int
		created      *configuredList
		failures     map[string]error
		failAt       map[string]int
		after        func(string)
		calls        map[string]int
		newRollbacks int
		closes       int
	}
)

var (
	_ runtime.List                  = (*configuredList)(nil)
	_ runtime.Factory[runtime.List] = (*configuredList)(nil)
)

func (l *configuredList) Length(ctx context.Context) (runtime.Int, error) {
	if err := l.operation(ctx, "Length"); err != nil {
		return 0, err
	}

	if l.length != nil {
		return *l.length, nil
	}

	return l.Array.Length(ctx)
}

func (l *configuredList) New(ctx context.Context) (runtime.List, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := l.operation(ctx, "New"); err != nil {
		l.newRollbacks++

		return nil, errors.Join(err, l.newCleanupErr)
	}

	l.created = &configuredList{Array: runtime.NewArray(0), backend: l.backend, expected: l.expected, failures: l.failures, failAt: l.failAt, after: l.after, closeErr: l.closeErr, calls: make(map[string]int)}

	return l.created, nil
}

func (l *configuredList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := l.operation(ctx, "At"); err != nil {
		return runtime.None, err
	}

	return l.Array.At(ctx, index)
}

func (l *configuredList) Append(ctx context.Context, value runtime.Value) error {
	if err := l.operation(ctx, "Append"); err != nil {
		return err
	}

	return l.Array.Append(ctx, value)
}

func (l *configuredList) Close() error {
	l.closes++

	return l.closeErr
}

func (l *configuredList) operation(ctx context.Context, name string) error {
	l.calls[name]++
	if ctx != l.expected {
		return errors.New("list context changed")
	}

	if l.after != nil {
		l.after(name)
	}

	if l.calls[name] < l.failAt[name] {
		return nil
	}

	return l.failures[name]
}
