package drivers

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTTPCookies struct {
	values map[string]HTTPCookie
}

func NewHTTPCookies() *HTTPCookies {
	return NewHTTPCookiesWith(make(map[string]HTTPCookie))
}

func NewHTTPCookiesWith(values map[string]HTTPCookie) *HTTPCookies {
	return &HTTPCookies{values}
}

func (c *HTTPCookies) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(c.values, jettison.NoHTMLEscaping())
}

func (c *HTTPCookies) String() string {
	j, err := c.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(j)
}

func (c *HTTPCookies) Compare(other core.Value) int64 {
	oc, ok := other.(*HTTPCookies)

	if !ok {
		// TODO: Implement
		return 1
	}

	switch {
	case len(c.values) > len(oc.values):
		return 1
	case len(c.values) < len(oc.values):
		return -1
	}

	for name := range c.values {
		cEl, cExists := c.Get(values.NewString(name))

		if !cExists {
			return -1
		}

		ocEl, ocExists := oc.Get(values.NewString(name))

		if !ocExists {
			return 1
		}

		c := cEl.Compare(ocEl)

		if c != 0 {
			return c
		}
	}

	return 0
}

func (c *HTTPCookies) Unwrap() interface{} {
	return c.values
}

func (c *HTTPCookies) Hash() uint64 {
	hash := fnv.New64a()

	hash.Write([]byte("HTMLCookies"))
	hash.Write([]byte(":"))
	hash.Write([]byte("{"))

	keys := make([]string, 0, len(c.values))

	for key := range c.values {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		hash.Write([]byte(key))
		hash.Write([]byte(":"))

		el := c.values[key]

		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		hash.Write(bytes)

		if idx != endIndex {
			hash.Write([]byte(","))
		}
	}

	hash.Write([]byte("}"))

	return hash.Sum64()
}

func (c *HTTPCookies) Copy() core.Value {
	return NewHTTPCookiesWith(c.values)
}

func (c *HTTPCookies) Clone() core.Cloneable {
	clone := make(map[string]HTTPCookie)

	for _, cookie := range c.values {
		clone[cookie.Name] = cookie
	}

	return NewHTTPCookiesWith(clone)
}

func (c *HTTPCookies) Length() values.Int {
	return values.NewInt(len(c.values))
}

func (c *HTTPCookies) Keys() []values.String {
	result := make([]values.String, 0, len(c.values))

	for k := range c.values {
		result = append(result, values.NewString(k))
	}

	return result
}

func (c *HTTPCookies) Values() []HTTPCookie {
	result := make([]HTTPCookie, 0, len(c.values))

	for _, v := range c.values {
		result = append(result, v)
	}

	return result
}

func (c *HTTPCookies) Get(key values.String) (HTTPCookie, values.Boolean) {
	value, found := c.values[key.String()]

	if found {
		return value, values.True
	}

	return HTTPCookie{}, values.False
}

func (c *HTTPCookies) Set(cookie HTTPCookie) {
	c.values[cookie.Name] = cookie
}

func (c *HTTPCookies) GetByKey(ctx context.Context, key string) (core.Value, error) {
	// TODO: Implement
	return values.None, nil
}

func (c *HTTPCookies) ForEach(predicate func(value HTTPCookie, key values.String) bool) {
	for key, val := range c.values {
		if !predicate(val, values.NewString(key)) {
			break
		}
	}
}
