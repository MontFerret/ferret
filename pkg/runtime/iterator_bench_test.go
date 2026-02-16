package runtime

import (
	"context"
	"testing"
)

type twoCallIter struct {
	n int
	i int
}

func (it *twoCallIter) HasNext(_ context.Context) (bool, error) {
	return it.i < it.n, nil
}

func (it *twoCallIter) Next(_ context.Context) (Value, Value, error) {
	it.i++
	return ZeroInt, ZeroInt, nil
}

type oneCallIter struct {
	n int
	i int
}

func (it *oneCallIter) Next(_ context.Context) (Value, Value, bool, error) {
	if it.i >= it.n {
		return None, None, false, nil
	}
	it.i++
	return ZeroInt, ZeroInt, true, nil
}

func BenchmarkIteratorTwoCall(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		it := &twoCallIter{n: 1024}
		for {
			ok, err := it.HasNext(ctx)
			if err != nil || !ok {
				break
			}
			if _, _, err := it.Next(ctx); err != nil {
				break
			}
		}
	}
}

func BenchmarkIteratorOneCall(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		it := &oneCallIter{n: 1024}
		for {
			_, _, ok, err := it.Next(ctx)
			if err != nil || !ok {
				break
			}
		}
	}
}
