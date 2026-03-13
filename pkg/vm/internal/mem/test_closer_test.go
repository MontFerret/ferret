package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type testCloser struct {
	name   string
	closed int
}

func newTestCloser(name string) *testCloser {
	return &testCloser{name: name}
}

func (c *testCloser) Close() error {
	c.closed++
	return nil
}

func (c *testCloser) String() string {
	return c.name
}

func (c *testCloser) Hash() uint64 {
	return 0
}

func (c *testCloser) Copy() runtime.Value {
	return c
}
