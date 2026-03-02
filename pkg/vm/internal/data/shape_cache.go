package data

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type (
	FastObject struct {
		cache         *ShapeCache
		shape         *fastShape
		slots         []runtime.Value
		size          int
		dict          map[string]runtime.Value
		dictThreshold int
	}

	fastShape struct {
		id     uint64
		fields map[string]int
		names  []string
	}

	ShapeCache struct {
		limit       int
		nextID      uint64
		transitions map[shapeKey]*fastShape
		root        *fastShape
	}

	shapeKey struct {
		shapeID uint64
		key     string
	}
)

type Shape = fastShape

func NewShapeCache(limit int) *ShapeCache {
	if limit < 0 {
		limit = 0
	}

	cache := &ShapeCache{
		limit:       limit,
		transitions: make(map[shapeKey]*fastShape),
	}
	cache.root = cache.newShape(nil, nil)

	return cache
}

func (c *ShapeCache) Root() *Shape {
	if c == nil {
		return nil
	}

	return c.root
}

func (c *ShapeCache) Transition(shape *Shape, key string) *Shape {
	if c == nil || shape == nil {
		return nil
	}

	if c.limit > 0 {
		k := shapeKey{shapeID: shape.id, key: key}
		if next, ok := c.transitions[k]; ok {
			return next
		}

		if len(c.transitions) >= c.limit {
			// Hard cap: stop caching new transitions.
			return c.newShapeFrom(shape, key)
		}

		next := c.newShapeFrom(shape, key)
		c.transitions[k] = next
		return next
	}

	return c.newShapeFrom(shape, key)
}

func (c *ShapeCache) nextShapeID() uint64 {
	c.nextID++
	return c.nextID
}

func (c *ShapeCache) newShape(fields map[string]int, names []string) *fastShape {
	if fields == nil {
		fields = make(map[string]int)
	}

	if names == nil {
		names = make([]string, 0)
	}

	return &fastShape{
		id:     c.nextShapeID(),
		fields: fields,
		names:  names,
	}
}

func (c *ShapeCache) newShapeFrom(prev *fastShape, key string) *fastShape {
	fields := make(map[string]int, len(prev.fields)+1)
	for k, v := range prev.fields {
		fields[k] = v
	}

	slot := len(prev.names)
	fields[key] = slot

	names := make([]string, slot+1)
	copy(names, prev.names)
	names[slot] = key

	return c.newShape(fields, names)
}
