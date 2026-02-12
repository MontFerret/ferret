package runtime

import (
	"context"
	"encoding/binary"
	"hash/fnv"

	"github.com/wI2L/jettison"
)

// Query represents a query literal used by the operator index.
type Query struct {
	Kind    String
	Payload String
}

func NewQuery(kind, payload String) Query {
	return Query{Kind: kind, Payload: payload}
}

func (q Query) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(map[string]String{
		"kind":    q.Kind,
		"payload": q.Payload,
	}, jettison.NoHTMLEscaping())
}

func (q Query) String() string {
	if q.Payload == EmptyString {
		return q.Kind.String()
	}

	return q.Kind.String() + ":" + q.Payload.String()
}

func (q Query) Unwrap() interface{} {
	return map[string]interface{}{
		"kind":    q.Kind.Unwrap(),
		"payload": q.Payload.Unwrap(),
	}
}

func (q Query) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("query:"))
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, q.Kind.Hash())
	h.Write(buf)
	binary.LittleEndian.PutUint64(buf, q.Payload.Hash())
	h.Write(buf)
	return h.Sum64()
}

func (q Query) Copy() Value {
	return q
}

// Queryable allows values to handle operator index queries.
type Queryable interface {
	ApplyQuery(ctx context.Context, q Query) (Value, error)
}
