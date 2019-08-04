package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/textproto"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// HTTPHeaders HTTP header object
type HTTPHeaders map[string][]string

func NewHTTPHeaders(values map[string][]string) HTTPHeaders {
	return HTTPHeaders(values)
}

func (h HTTPHeaders) Type() core.Type {
	return HTTPHeaderType
}

func (h HTTPHeaders) String() string {
	var buf bytes.Buffer

	for k := range h {
		buf.WriteString(fmt.Sprintf("%s=%s;", k, h.Get(k)))
	}

	return buf.String()
}

func (h HTTPHeaders) Compare(other core.Value) int64 {
	if other.Type() != HTTPHeaderType {
		return Compare(HTTPHeaderType, other.Type())
	}

	oh := other.(HTTPHeaders)

	if len(h) > len(oh) {
		return 1
	} else if len(h) < len(oh) {
		return -1
	}

	for k := range h {
		c := strings.Compare(h.Get(k), oh.Get(k))

		if c != 0 {
			return int64(c)
		}
	}

	return 0
}

func (h HTTPHeaders) Unwrap() interface{} {
	return h
}

func (h HTTPHeaders) Hash() uint64 {
	hash := fnv.New64a()

	hash.Write([]byte(h.Type().String()))
	hash.Write([]byte(":"))
	hash.Write([]byte("{"))

	keys := make([]string, 0, len(h))

	for key := range h {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		hash.Write([]byte(key))
		hash.Write([]byte(":"))

		value := h.Get(key)

		hash.Write([]byte(value))

		if idx != endIndex {
			hash.Write([]byte(","))
		}
	}

	hash.Write([]byte("}"))

	return hash.Sum64()
}

func (h HTTPHeaders) Copy() core.Value {
	return *(&h)
}

func (h HTTPHeaders) MarshalJSON() ([]byte, error) {
	headers := map[string]string{}

	for key, val := range h {
		headers[key] = strings.Join(val, ", ")
	}

	out, err := json.Marshal(headers)

	if err != nil {
		return nil, err
	}

	return out, err
}

func (h HTTPHeaders) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func (h HTTPHeaders) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

func (h HTTPHeaders) GetIn(_ context.Context, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return values.None, nil
	}

	segment := path[0]

	err := core.ValidateType(segment, types.String)

	if err != nil {
		return values.None, err
	}

	return values.NewString(h.Get(segment.String())), nil
}
