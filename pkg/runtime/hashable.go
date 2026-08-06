package runtime

import (
	"encoding/binary"
	"hash/fnv"
	"sort"
)

// Hashable represents an interface of any type that can be hashed.
type Hashable interface {
	// Hash returns a hash value for the object.
	Hash() uint64
}

// Hash computes a hash value for the given typename and content using the FNV-1a hashing algorithm.
// It concatenates the typename, a colon, and the content bytes to generate a unique hash value.
// This function can be used to create consistent hash values for different types of content based on their type and content.
func Hash(typename string, content []byte) uint64 {
	h := fnv.New64a()

	h.Write([]byte(typename))
	h.Write([]byte(":"))
	h.Write(content)

	return h.Sum64()
}

// MapHash computes a hash value for a map of string keys to Value values.
// It uses the FNV-1a hashing algorithm to generate a consistent hash based on the keys and their corresponding value hashes.
// The keys are sorted to ensure that the order of key-value pairs does not affect the resulting hash.
func MapHash(input map[string]Value) uint64 {
	h := fnv.New64a()

	keys := make([]string, 0, len(input))

	for key := range input {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	h.Write([]byte("{"))

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := input[key]

		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		h.Write(bytes)

		if idx != endIndex {
			h.Write([]byte(","))
		}
	}

	h.Write([]byte("}"))

	return h.Sum64()
}
