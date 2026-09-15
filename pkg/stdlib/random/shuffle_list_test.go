package random

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	shuffleBackend struct {
		location string
		pageSize int
	}

	// The source supports only New and ForEach. Destinations additionally
	// support Append, Length, and Swap, with independent contents in one backend.
	shuffleList struct {
		*traversalList
		expected    context.Context
		failure     error
		cleanupErr  error
		after       func(string)
		backend     *shuffleBackend
		created     *shuffleList
		length      *runtime.Int
		operations  map[string]int
		failOp      string
		failAt      int
		rollbacks   int
		destination bool
	}
)

func newShuffleList(source *traversalList) *shuffleList {
	return &shuffleList{
		traversalList: source,
		backend:       &shuffleBackend{location: "storage/orders", pageSize: 128},
		operations:    make(map[string]int),
	}
}

func (list *shuffleList) New(ctx context.Context) (runtime.List, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	list.created = &shuffleList{
		traversalList: &traversalList{},
		expected:      list.expected,
		failure:       list.failure,
		cleanupErr:    list.cleanupErr,
		after:         list.after,
		backend:       list.backend,
		length:        list.length,
		operations:    make(map[string]int),
		failOp:        list.failOp,
		failAt:        list.failAt,
		destination:   true,
	}

	if err := list.operation(ctx, "New"); err != nil {
		list.rollbacks++

		return nil, errors.Join(err, list.created.Close())
	}

	return list.created, nil
}

func (list *shuffleList) Append(ctx context.Context, value runtime.Value) error {
	if err := list.operation(ctx, "Append"); err != nil {
		return err
	}

	list.values = append(list.values, value)

	return nil
}

func (list *shuffleList) Length(ctx context.Context) (runtime.Int, error) {
	if err := list.operation(ctx, "Length"); err != nil {
		return 0, err
	}

	if list.length != nil {
		return *list.length, nil
	}

	return runtime.Int(len(list.values)), nil
}

func (list *shuffleList) Swap(ctx context.Context, i, j runtime.Int) error {
	if err := list.operation(ctx, "Swap"); err != nil {
		return err
	}

	if i < 0 || j < 0 || i >= runtime.Int(len(list.values)) || j >= runtime.Int(len(list.values)) {
		return runtime.Error(runtime.ErrInvalidOperation, "out of bounds")
	}

	list.values[i], list.values[j] = list.values[j], list.values[i]

	return nil
}

func (list *shuffleList) Close() error {
	list.closes++

	return list.cleanupErr
}

func (list *shuffleList) operation(ctx context.Context, name string) error {
	list.operations[name]++

	if list.expected != nil && list.expected != ctx {
		return errors.New("host operation received the wrong context")
	}

	if !list.destination && name != "New" {
		return errors.New("source does not support " + name)
	}

	if list.after != nil {
		list.after(name)
	}

	if name == list.failOp && list.operations[name] >= list.failAt {
		return list.failure
	}

	return nil
}
