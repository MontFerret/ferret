package engine

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type trackingJSONCloser struct {
	closeErr  error
	encodeErr error
	name      string
	payload   string
	events    []string
	closed    int
}

func newTrackingJSONCloser(name, payload string) *trackingJSONCloser {
	return &trackingJSONCloser{name: name, payload: payload}
}

func (c *trackingJSONCloser) Close() error {
	c.events = append(c.events, "close")
	c.closed++

	return c.closeErr
}

func (c *trackingJSONCloser) MarshalJSON() ([]byte, error) {
	c.events = append(c.events, "encode")

	if c.closed != 0 {
		return nil, errors.New("encoding a closed resource")
	}

	return []byte(c.payload), c.encodeErr
}

func (c *trackingJSONCloser) String() string {
	return c.name
}

func (c *trackingJSONCloser) Hash() uint64 {
	return runtime.NewString(c.name).Hash()
}

func (c *trackingJSONCloser) Copy() runtime.Value {
	return c
}
