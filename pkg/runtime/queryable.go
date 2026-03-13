package runtime

import (
	"context"
)

type (
	// Query represents a query literal used by the operator index.
	Query struct {
		Options Value  `json:"options"`
		Kind    String `json:"kind"`
		Payload String `json:"payload"`
	}

	// Queryable allows values to handle operator index queries.
	Queryable interface {
		Query(ctx context.Context, q Query) (List, error)
	}
)
