package runtime

import (
	"context"
	"encoding/binary"
	"hash/fnv"
)

type (
	// Query represents a query literal used by the operator index.
	Query struct {
		Kind    String `json:"kind"`
		Payload String `json:"payload"`
		Params  Value  `json:"params"`
	}

	// Queryable allows values to handle operator index queries.
	Queryable interface {
		Query(ctx context.Context, q Query) (Value, error)
	}
)

func NewQuery(kind, payload String) Query {
	return Query{Kind: kind, Payload: payload, Params: None}
}

func (q Query) String() string {
	if q.Payload == EmptyString {
		return q.Kind.String()
	}

	return q.Kind.String() + ":" + q.Payload.String()
}

func (q Query) Hash() uint64 {
	h := fnv.New64a()
	h.Write([]byte("query:"))
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, q.Kind.Hash())
	h.Write(buf)
	binary.LittleEndian.PutUint64(buf, q.Payload.Hash())
	h.Write(buf)
	binary.LittleEndian.PutUint64(buf, q.Params.Hash())
	h.Write(buf)
	return h.Sum64()
}

func (q Query) Copy() Value {
	return q
}
