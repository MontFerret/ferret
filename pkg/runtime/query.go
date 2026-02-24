package runtime

import (
	"context"
)

type (
	// Query represents a query literal used by the operator index.
	Query struct {
		Kind    String `json:"kind"`
		Payload String `json:"payload"`
		Options Value  `json:"options"`
	}

	// Queryable allows values to handle operator index queries.
	Queryable interface {
		Query(ctx context.Context, q Query) (Value, error)
	}
)
