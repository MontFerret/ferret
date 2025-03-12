package drivers

import (
	"bytes"
	"context"
	"fmt"
	"hash/fnv"
	"net/textproto"
	"sort"
	"strings"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// HTTPHeaders HTTP header object
type HTTPHeaders struct {
	values map[string][]string
}

func NewHTTPHeaders() *HTTPHeaders {
	return NewHTTPHeadersWith(make(map[string][]string))
}

func NewHTTPHeadersWith(values map[string][]string) *HTTPHeaders {
	return &HTTPHeaders{values}
}

func (h *HTTPHeaders) Length() core.Int {
	return core.NewInt(len(h.values))
}

func (h *HTTPHeaders) Type() core.Type {
	return HTTPHeaderType
}

func (h *HTTPHeaders) String() string {
	var buf bytes.Buffer

	for k := range h.values {
		buf.WriteString(fmt.Sprintf("%s=%s;", k, h.Get(k)))
	}

	return buf.String()
}

func (h *HTTPHeaders) Compare(other core.Value) int64 {
	oh, ok := other.(*HTTPHeaders)

	if !ok {
		return CompareTypes(h.Type(), core.Reflect(other))
	}

	if len(h.values) > len(oh.values) {
		return 1
	} else if len(h.values) < len(oh.values) {
		return -1
	}

	for k := range h.values {
		c := strings.Compare(h.Get(k), oh.Get(k))

		if c != 0 {
			return int64(c)
		}
	}

	return 0
}

func (h *HTTPHeaders) Unwrap() interface{} {
	return h.values
}

func (h *HTTPHeaders) Hash() uint64 {
	hash := fnv.New64a()

	hash.Write([]byte(h.Type().String()))
	hash.Write([]byte(":"))
	hash.Write([]byte("{"))

	keys := make([]string, 0, len(h.values))

	for key := range h.values {
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

func (h *HTTPHeaders) Copy() core.Value {
	return &HTTPHeaders{h.values}
}

func (h *HTTPHeaders) Clone() core.Cloneable {
	cp := make(map[string][]string)

	for k, v := range h.values {
		cp[k] = v
	}

	return &HTTPHeaders{cp}
}

func (h *HTTPHeaders) MarshalJSON() ([]byte, error) {
	headers := map[string]string{}

	for key, val := range h.values {
		headers[key] = strings.Join(val, ", ")
	}

	out, err := jettison.MarshalOpts(headers)

	if err != nil {
		return nil, err
	}

	return out, err
}

func (h *HTTPHeaders) Set(key, value string) {
	textproto.MIMEHeader(h.values).Set(key, value)
}

func (h *HTTPHeaders) SetArr(key string, value []string) {
	h.values[key] = value
}

func (h *HTTPHeaders) Get(key string) string {
	_, found := h.values[key]

	if !found {
		return ""
	}

	return textproto.MIMEHeader(h.values).Get(key)
}

func (h *HTTPHeaders) GetHeader(_ context.Context, key string) (core.Value, error) {
	if key == "" {
		return core.None, nil
	}

	return core.NewString(h.Get(key)), nil
}

func (h *HTTPHeaders) ForEach(predicate func(value []string, key string) bool) {
	for key, val := range h.values {
		if !predicate(val, key) {
			break
		}
	}
}
