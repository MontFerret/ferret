package vm

import (
	"sync/atomic"
)

type failingCloser struct {
	err    error
	closes atomic.Int32
}

func newFailingCloser(err error) *failingCloser {
	return &failingCloser{err: err}
}

func (c *failingCloser) Close() error {
	c.closes.Add(1)
	return c.err
}

func (c *failingCloser) count() int {
	return int(c.closes.Load())
}
