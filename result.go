package ferret

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	// ErrNoResult indicates that the result contains no values.
	// It is returned by methods that expect at least one value, such as First.
	ErrNoResult = errors.New("result contains no values")

	// ErrNotScalar indicates that the result represents a sequence of values
	// rather than a single top-level scalar value.
	ErrNotScalar = errors.New("result is not scalar")
)

// Result represents the output of a Ferret query execution.
//
// A result can be either:
//
//   - scalar: a single top-level value such as a number, string, object, or map
//   - sequence: a stream of values produced by an iterable result
//
// Result provides both low-level cursor-style access via Next and higher-level
// convenience methods such as First, Collect, and ForEach.
//
// Result is single-consumption for streaming operations. Methods such as Next,
// First, Collect, and ForEach consume values from the current position.
// In contrast, Value is a non-consuming accessor for scalar results only.
type Result struct {
	value    runtime.Value
	iterable runtime.Iterable

	iter   runtime.Iterator
	err    error
	done   bool
	scalar bool
}

// newResult constructs a Result from a runtime value.
//
// Non-iterable values are treated as scalar results.
// Iterable values are treated as streaming results, except runtime.Map,
// which is considered a single top-level scalar value for API purposes.
func newResult(val runtime.Value) *Result {
	iterable, ok := val.(runtime.Iterable)
	if !ok {
		return &Result{
			value:  val,
			scalar: true,
		}
	}

	if _, ok := val.(runtime.Map); ok {
		return &Result{
			value:  val,
			scalar: true,
		}
	}

	return &Result{
		value:    val,
		iterable: iterable,
		scalar:   false,
	}
}

// Close releases resources associated with the result.
//
// If the underlying iterator or root value implements io.Closer, Close calls it.
// After Close returns, the result becomes terminal and any subsequent call to
// Next will return io.EOF.
//
// Close is safe to call even if iteration has not started.
func (r *Result) Close() error {
	var err error

	if r.iter != nil {
		if c, ok := r.iter.(io.Closer); ok {
			err = errors.Join(err, c.Close())
		}
	}

	if c, ok := r.value.(io.Closer); ok {
		err = errors.Join(err, c.Close())
	}

	r.iter = nil
	r.iterable = nil
	r.value = runtime.None
	r.done = true

	return err
}

// IsScalar reports whether the result represents a single top-level value.
//
// Scalar results can be accessed with Value.
// Non-scalar results must be consumed through Next, First, Collect, or ForEach.
func (r *Result) IsScalar() bool {
	return r.scalar
}

// Next returns the next available value from the result.
//
// For scalar results, Next returns the scalar value once and then returns io.EOF
// on subsequent calls.
//
// For non-scalar results, Next lazily initializes the underlying iterator on the
// first call and then returns values one by one until the sequence is exhausted.
//
// Once a non-EOF error occurs, that error is remembered and returned again by
// subsequent calls.
//
// Next consumes the result from its current position.
func (r *Result) Next(ctx context.Context) (runtime.Value, error) {
	if r.err != nil {
		return runtime.None, r.err
	}

	if r.done {
		return runtime.None, io.EOF
	}

	if r.scalar {
		r.done = true
		return r.value, nil
	}

	if r.iter == nil {
		r.iter, r.err = r.iterable.Iterate(ctx)
		if r.err != nil {
			return runtime.None, r.err
		}
	}

	val, _, err := r.iter.Next(ctx)
	if errors.Is(err, io.EOF) {
		r.done = true
		return runtime.None, io.EOF
	}

	if err != nil {
		r.err = err
		return runtime.None, err
	}

	return val, nil
}

// Value returns the underlying scalar value without consuming the result.
//
// Value succeeds only when the result is scalar. If the result represents a
// sequence of values, Value returns ErrNotScalar.
//
// Unlike Next, Value does not advance or modify iteration state.
func (r *Result) Value() (runtime.Value, error) {
	if !r.scalar {
		return runtime.None, ErrNotScalar
	}

	return r.value, nil
}

// First returns the first available value from the result.
//
// For scalar results, it returns the scalar value.
// For non-scalar results, it returns the first item produced by the iterator.
//
// If the result contains no values, First returns ErrNoResult.
//
// First consumes one value from the result.
func (r *Result) First(ctx context.Context) (runtime.Value, error) {
	val, err := r.Next(ctx)
	if errors.Is(err, io.EOF) {
		return runtime.None, ErrNoResult
	}
	if err != nil {
		return runtime.None, err
	}

	return val, nil
}

// Collect consumes all remaining values from the result and returns them as a slice.
//
// For scalar results, Collect returns a slice containing the scalar value.
// For non-scalar results, it returns all remaining items from the iterator.
//
// If an iteration error occurs, Collect stops and returns that error.
func (r *Result) Collect(ctx context.Context) ([]runtime.Value, error) {
	var res []runtime.Value

	if err := r.ForEach(ctx, func(val runtime.Value) error {
		res = append(res, val)
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

// ForEach consumes all remaining values from the result and calls fn for each one.
//
// For scalar results, fn is called once with the scalar value.
// For non-scalar results, fn is called once per remaining item in the sequence.
//
// Iteration stops when:
//
//   - the result is exhausted
//   - fn returns an error
//   - the underlying iteration returns an error
//
// In the latter two cases, that error is returned.
func (r *Result) ForEach(ctx context.Context, fn func(runtime.Value) error) error {
	for {
		val, err := r.Next(ctx)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		if err := fn(val); err != nil {
			return err
		}
	}
}
