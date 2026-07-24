package sdk

import "fmt"

type (
	// DecodeErrorKind classifies failures reported by Decode and related helpers.
	DecodeErrorKind string

	// DecodeError describes a failed Ferret-to-Go conversion.
	//
	// SafeToExpose distinguishes deterministic conversion diagnostics from errors
	// returned by runtime values or callbacks, which may contain sensitive details.
	DecodeError struct {
		cause      error
		path       string
		kind       DecodeErrorKind
		safe       bool
		renderPath bool
	}
)

const (
	// DecodeErrorKindType reports a source or target type mismatch.
	DecodeErrorKindType DecodeErrorKind = "type"
	// DecodeErrorKindUnknownField reports an object field rejected by strict decoding.
	DecodeErrorKindUnknownField DecodeErrorKind = "unknown_field"
	// DecodeErrorKindNone reports an explicit None rejected by DisallowNoneValues.
	DecodeErrorKindNone DecodeErrorKind = "none"
	// DecodeErrorKindRange reports a numeric or collection-size overflow.
	DecodeErrorKindRange DecodeErrorKind = "range"
	// DecodeErrorKindSource reports a failure returned by a runtime value, iterator, or caller context.
	DecodeErrorKindSource DecodeErrorKind = "source"
)

func newDecodeError(path string, kind DecodeErrorKind, cause error, safe bool) *DecodeError {
	return &DecodeError{
		path:       path,
		kind:       kind,
		cause:      cause,
		safe:       safe,
		renderPath: true,
	}
}

func newDecodeErrorWithoutPath(path string, kind DecodeErrorKind, cause error, safe bool) *DecodeError {
	err := newDecodeError(path, kind, cause, safe)
	err.renderPath = false

	return err
}

// Error renders the conversion path followed by the underlying error.
func (e *DecodeError) Error() string {
	if e == nil {
		return ""
	}

	if e.cause == nil {
		return e.path
	}

	if !e.renderPath || e.path == "" {
		return e.cause.Error()
	}

	return fmt.Sprintf("%s: %s", e.path, e.cause)
}

// Unwrap returns the underlying conversion or source failure.
func (e *DecodeError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

// Path returns the JSONPath-like location at which decoding failed.
func (e *DecodeError) Path() string {
	if e == nil {
		return ""
	}

	return e.path
}

// Kind returns the stable failure category.
func (e *DecodeError) Kind() DecodeErrorKind {
	if e == nil {
		return ""
	}

	return e.kind
}

// SafeToExpose reports whether Error contains only SDK-generated conversion detail.
func (e *DecodeError) SafeToExpose() bool {
	return e != nil && e.safe
}
