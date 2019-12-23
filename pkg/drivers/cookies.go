package drivers

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"hash/fnv"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type HTTPCookies map[string]HTTPCookie

func NewHTTPCookies() HTTPCookies {
	return make(HTTPCookies)
}

func (c HTTPCookies) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]HTTPCookie(c))
}

func (c HTTPCookies) Type() core.Type {
	return HTTPCookiesType
}

func (c HTTPCookies) String() string {
	j, err := c.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(j)
}

func (c HTTPCookies) Compare(other core.Value) int64 {
	if other.Type() != HTTPCookiesType {
		return Compare(HTTPCookiesType, other.Type())
	}

	oc := other.(HTTPCookies)

	switch {
	case len(c) > len(oc):
		return 1
	case len(c) < len(oc):
		return -1
	}

	for name := range c {
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

func (c HTTPCookies) Unwrap() interface{} {
	return map[string]HTTPCookie(c)
}

func (c HTTPCookies) Hash() uint64 {
	hash := fnv.New64a()

	hash.Write([]byte(c.Type().String()))
	hash.Write([]byte(":"))
	hash.Write([]byte("{"))

	keys := make([]string, 0, len(c))

	for key := range c {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		hash.Write([]byte(key))
		hash.Write([]byte(":"))

		el := c[key]

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

func (c HTTPCookies) Copy() core.Value {
	copied := make(HTTPCookies)

	for k, v := range c {
		copied[k] = v
	}

	return copied
}

func (c HTTPCookies) Length() values.Int {
	return values.NewInt(len(c))
}

func (c HTTPCookies) Keys() []values.String {
	keys := make([]values.String, 0, len(c))

	for k := range c {
		keys = append(keys, values.NewString(k))
	}

	return keys
}

func (c HTTPCookies) Get(key values.String) (core.Value, values.Boolean) {
	value, found := c[key.String()]

	if found {
		return value, values.True
	}

	return values.None, values.False
}

func (c HTTPCookies) Set(key values.String, value core.Value) {
	if cookie, ok := value.(HTTPCookie); ok {
		c[key.String()] = cookie
	}
}

func (c HTTPCookies) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return values.None, nil
	}

	segment := path[0]

	err := core.ValidateType(segment, types.String)

	if err != nil {
		return values.None, err
	}

	cookie, found := c[segment.String()]

	if found {
		if len(path) == 1 {
			return cookie, nil
		}

		return values.GetIn(ctx, cookie, path[1:])
	}

	return values.None, nil
}
