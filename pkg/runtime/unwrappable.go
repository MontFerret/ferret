package runtime

// Unwrappable represents an interface that can be unwrapped to get the underlying value.
type Unwrappable interface {
	Unwrap() any
}
