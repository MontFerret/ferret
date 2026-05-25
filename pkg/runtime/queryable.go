package runtime

import (
	"context"
)

const (
	// QueryOneErrorMessage is the standard runtime failure for QUERY ONE cardinality mismatches.
	QueryOneErrorMessage = "QUERY ONE expected exactly one match"
)

type (
	// Query represents a query literal used by the operator index.
	Query struct {
		Options Value  `json:"options"`
		Kind    String `json:"kind"`
		Payload String `json:"payload"` // TODO: Rename to "expression" or "value"
	}

	// QueryFunc is the list-returning query implementation used by default modifier helpers.
	QueryFunc func(context.Context, Query) (List, error)

	// Queryable allows values to handle operator index queries and query modifiers.
	Queryable interface {
		Query(ctx context.Context, q Query) (List, error)
		QueryOne(ctx context.Context, q Query) (Value, error)
		QueryCount(ctx context.Context, q Query) (Int, error)
		QueryExists(ctx context.Context, q Query) (Boolean, error)
	}
)

// DefaultQueryOne implements QUERY ONE by materializing the list-returning query result.
func DefaultQueryOne(ctx context.Context, q Query, query QueryFunc) (Value, error) {
	out, err := query(ctx, q)
	if err != nil {
		return None, err
	}

	length, err := queryResultLength(ctx, out)
	if err != nil {
		return None, err
	}

	if length != 1 {
		return None, Error(ErrInvalidOperation, QueryOneErrorMessage)
	}

	return out.At(ctx, ZeroInt)
}

// DefaultQueryCount implements QUERY COUNT by materializing the list-returning query result.
func DefaultQueryCount(ctx context.Context, q Query, query QueryFunc) (Int, error) {
	out, err := query(ctx, q)
	if err != nil {
		return ZeroInt, err
	}

	return queryResultLength(ctx, out)
}

// DefaultQueryExists implements QUERY EXISTS by materializing the list-returning query result.
func DefaultQueryExists(ctx context.Context, q Query, query QueryFunc) (Boolean, error) {
	length, err := DefaultQueryCount(ctx, q, query)
	if err != nil {
		return False, err
	}

	return NewBoolean(length != 0), nil
}

func queryResultLength(ctx context.Context, out List) (Int, error) {
	if out == nil {
		return ZeroInt, nil
	}

	return out.Length(ctx)
}
