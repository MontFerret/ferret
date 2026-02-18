package runtime

// Hashable represents an interface of any type that can be hashed.
type Hashable interface {
	// Hash returns a hash value for the object.
	Hash() uint64
}
