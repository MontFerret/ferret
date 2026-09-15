package stdlib_test

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// These hosts intentionally ignore cancellation inside individual operations.
	// Returned errors remain independent of the caller context state.
	contextList struct {
		runtime.List
		expected context.Context
		failure  error
		cleanup  error
		cancel   context.CancelFunc
		calls    map[string]int
		created  *contextList
		stage    string
		values   []runtime.Value
		cancelAt int
		closes   int
	}

	contextIterable struct {
		runtime.Value
		list *contextList
	}

	contextIterator struct {
		list  *contextList
		index int
	}

	contextMap struct {
		runtime.Map
		list *contextList
	}
)

func (list *contextList) ForEach(ctx context.Context, visit runtime.IndexReadablePredicate) error {
	for _, value := range list.values {
		if err := list.observe(ctx, "visit"); err != nil {
			return err
		}

		// Repeated, nonzero indices must not affect counting or separators.
		more, err := visit(ctx, value, 7)
		if err != nil {
			return err
		}

		if !more {
			break
		}
	}

	return nil
}

func (list *contextList) Length(ctx context.Context) (runtime.Int, error) {
	if err := list.observe(ctx, "Length"); err != nil {
		return 0, err
	}

	return runtime.Int(len(list.values)), nil
}

func (list *contextList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := list.observe(ctx, "At"); err != nil {
		return runtime.None, err
	}

	return list.values[index], nil
}

func (list *contextList) Append(ctx context.Context, value runtime.Value) error {
	if err := list.observe(ctx, "Append"); err != nil {
		return err
	}

	list.values = append(list.values, value)

	return nil
}

func (list *contextList) SetAt(ctx context.Context, index runtime.Int, value runtime.Value) error {
	if err := list.observe(ctx, "SetAt"); err != nil {
		return err
	}

	list.values[index] = value

	return nil
}

func (list *contextList) RemoveAt(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if err := list.observe(ctx, "RemoveAt"); err != nil {
		return runtime.None, err
	}

	value := list.values[index]
	list.values = append(list.values[:index], list.values[index+1:]...)

	return value, nil
}

func (list *contextList) Swap(ctx context.Context, first, second runtime.Int) error {
	if err := list.observe(ctx, "Swap"); err != nil {
		return err
	}

	list.values[first], list.values[second] = list.values[second], list.values[first]

	return nil
}

func (list *contextList) New(ctx context.Context) (runtime.List, error) {
	if err := list.observe(ctx, "New"); err != nil {
		return nil, err
	}

	list.created = &contextList{
		expected: list.expected, cancel: list.cancel, stage: list.stage,
		cancelAt: list.cancelAt, failure: list.failure, cleanup: list.cleanup,
		calls: make(map[string]int),
	}

	return list.created, nil
}

func (list *contextList) Close() error {
	list.closes++

	return list.cleanup
}

func (list *contextList) observe(ctx context.Context, stage string) error {
	if ctx != list.expected {
		return errors.New("caller context changed")
	}

	list.calls[stage]++
	if stage == list.stage && list.calls[stage] == list.cancelAt {
		list.cancel()

		return list.failure
	}

	return nil
}

func (value *contextIterable) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if err := value.list.observe(ctx, "Iterate"); err != nil {
		return nil, err
	}

	return &contextIterator{list: value.list}, nil
}

func (iter *contextIterator) Next(ctx context.Context) (runtime.Value, runtime.Value, error) {
	if iter.index == len(iter.list.values) {
		return runtime.None, runtime.None, io.EOF
	}

	if err := iter.list.observe(ctx, "visit"); err != nil {
		return runtime.None, runtime.None, err
	}

	value := iter.list.values[iter.index]
	iter.index++

	return value, runtime.Int(7), nil
}

func (value *contextMap) ForEach(ctx context.Context, visit runtime.KeyReadablePredicate) error {
	return value.list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		return visit(c, value, runtime.String("key"))
	})
}

func (value *contextMap) Values(ctx context.Context) (runtime.List, error) {
	if err := value.list.observe(ctx, "Values"); err != nil {
		return nil, err
	}

	return value.list, nil
}

func (value *contextMap) ContainsKey(ctx context.Context, _ runtime.Value) (runtime.Boolean, error) {
	return runtime.False, value.list.observe(ctx, "ContainsKey")
}
