package runtime

import (
	"context"
)

type (
	// Ordering is the normalized result of a strict runtime comparison.
	Ordering int8

	// Equatable gives a host value an independent, fallible equality contract.
	//
	// Equal must be reflexive, symmetric, and transitive within the value's
	// comparison domain. It returns false without an error for an incompatible
	// value. If the value also implements Comparable, Equal must return true
	// exactly when Compare returns Equal. Semantically equal values must have equal
	// hashes; hashes are only candidate selectors and never prove equality.
	// Runtime dispatch forwards ctx without polling it; implementations that may
	// block are responsible for observing cancellation themselves.
	Equatable interface {
		Equal(ctx context.Context, other Value) (bool, error)
	}

	// Comparable gives a host value a fallible total ordering within its
	// comparison domain.
	//
	// Compare must be reflexive, antisymmetric, and transitive. It returns an
	// ErrInvalidOperation-compatible error for an incompatible value and returns
	// operational failures, including implementation-owned context cancellation,
	// unchanged. Runtime dispatch does not poll ctx. Any
	// negative or positive result is normalized by runtime dispatch to Less or
	// Greater. If the value also implements Equatable, Equal must agree with
	// Equal returning true.
	Comparable interface {
		Compare(ctx context.Context, other Value) (Ordering, error)
	}

	// Comparator compares two values for a fallible sort. Runtime sorting forwards
	// ctx to the comparator without polling it.
	Comparator = func(ctx context.Context, first, second Value) (Ordering, error)
)

const (
	Less    Ordering = -1
	Equal   Ordering = 0
	Greater Ordering = 1
)
