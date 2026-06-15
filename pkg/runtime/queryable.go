package runtime

import (
	"context"
)

type (
	// Query describes a query operation passed to a Queryable value.
	// Params contains query input supplied by WITH, while Options contains
	// execution policy supplied by OPTIONS.
	Query struct {
		Options    Value  `json:"options"`
		Params     Value  `json:"params"`
		Kind       String `json:"kind"`
		Expression String `json:"expression"`
	}

	// QueryFunc is the list-returning query implementation used by default modifier helpers.
	QueryFunc func(context.Context, Query) (List, error)

	// Queryable allows values to handle operator index queries and query modifiers.
	Queryable interface {
		// Query returns every matching value.
		Query(ctx context.Context, q Query) (List, error)
		// QueryOne returns the first matching value or None when there is no match.
		QueryOne(ctx context.Context, q Query) (Value, error)
		// QueryCount returns the number of matching values.
		QueryCount(ctx context.Context, q Query) (Int, error)
		// QueryExists reports whether at least one value matches.
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

	if length == 0 {
		return None, nil
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
