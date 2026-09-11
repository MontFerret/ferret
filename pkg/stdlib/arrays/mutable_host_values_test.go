package arrays_test

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	mutableHostList struct {
		runtime.List
		expectedContext context.Context
		failure         error
		calls           map[string]int
		failAt          map[string]int
		length          *runtime.Int
		after           func(string)
		nilRemoval      bool
	}
	readableArray struct {
		runtime.Value
		runtime.IndexReadable
	}
	mutableBorrowedValue struct {
		runtime.Value
		copies int
		closes int
	}
	failingOrderingValue struct {
		runtime.Value
		failure error
	}
)

func (v failingOrderingValue) Compare(context.Context, runtime.Value) (runtime.Ordering, error) {
	return runtime.Equal, v.failure
}

func (v *mutableBorrowedValue) Copy() runtime.Value {
	v.copies++

	return v
}

func (v *mutableBorrowedValue) Close() error {
	v.closes++

	return nil
}

func (l *mutableHostList) Copy() runtime.Value {
	l.calls["Copy"]++

	return runtime.None
}

func (l *mutableHostList) Clone(context.Context) (runtime.Cloneable, error) {
	l.calls["Clone"]++

	return nil, l.failure
}

func (l *mutableHostList) New(context.Context) (runtime.List, error) {
	l.calls["New"]++

	return nil, l.failure
}

func (l *mutableHostList) Length(ctx context.Context) (runtime.Int, error) {
	if err := l.before(ctx, "Length"); err != nil {
		return 0, err
	}

	if l.length != nil {
		return *l.length, nil
	}

	return l.List.Length(ctx)
}

func (l *mutableHostList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := l.before(ctx, "At"); err != nil {
		return runtime.None, err
	}

	return l.List.At(ctx, index)
}

func (l *mutableHostList) Append(ctx context.Context, value runtime.Value) error {
	if err := l.before(ctx, "Append"); err != nil {
		return err
	}

	return l.List.Append(ctx, value)
}

func (l *mutableHostList) SetAt(ctx context.Context, index runtime.Int, value runtime.Value) error {
	if err := l.before(ctx, "SetAt"); err != nil {
		return err
	}

	return l.List.SetAt(ctx, index, value)
}

func (l *mutableHostList) Insert(ctx context.Context, index runtime.Int, value runtime.Value) error {
	if err := l.before(ctx, "Insert"); err != nil {
		return err
	}

	return l.List.Insert(ctx, index, value)
}

func (l *mutableHostList) RemoveAt(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := l.before(ctx, "RemoveAt"); err != nil {
		return runtime.None, err
	}

	if l.nilRemoval {
		return nil, nil
	}

	return l.List.RemoveAt(ctx, index)
}

func (l *mutableHostList) Clear(ctx context.Context) error {
	if err := l.before(ctx, "Clear"); err != nil {
		return err
	}

	return l.List.Clear(ctx)
}

func (l *mutableHostList) SortAsc(ctx context.Context) error {
	if err := l.before(ctx, "SortAsc"); err != nil {
		return err
	}

	return runtime.SortListAsc(ctx, l)
}

func (l *mutableHostList) Swap(ctx context.Context, a, b runtime.Int) error {
	if err := l.before(ctx, "Swap"); err != nil {
		return err
	}

	return l.List.Swap(ctx, a, b)
}

func newMutableHostList(list runtime.List) *mutableHostList {
	return &mutableHostList{
		List:    list,
		failure: errors.New("host operation failed"),
		calls:   make(map[string]int),
		failAt:  make(map[string]int),
	}
}

func (l *mutableHostList) before(ctx context.Context, name string) error {
	l.calls[name]++
	if l.expectedContext != nil && ctx != l.expectedContext {
		return errors.New("caller context replaced")
	}

	if l.after != nil {
		l.after(name)
	}

	if l.failAt[name] == l.calls[name] {
		return l.failure
	}

	return nil
}
