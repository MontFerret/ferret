package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Collector struct {
	values map[uint64][]*KeyValuePair
}

func NewCollector() *Collector {
	return &Collector{values: make(map[uint64][]*KeyValuePair)}
}

func (c *Collector) Add(key core.Value, value core.Value) {
	hash := key.Hash()
	values, exists := c.values[hash]

	if !exists {
		values = make([]*KeyValuePair, 0, 5)
		c.values[hash] = values
	}

	c.values[hash] = append(values, &KeyValuePair{key, value})
}

func (c *Collector) MarshalJSON() ([]byte, error) {
	panic("not supported")
}

func (c *Collector) String() string {
	return "[Collector]"
}

func (c *Collector) Unwrap() interface{} {
	panic("not supported")
}

func (c *Collector) Hash() uint64 {
	panic("not supported")
}

func (c *Collector) Copy() core.Value {
	panic("not supported")
}
