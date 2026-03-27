package artifact

import "bytes"

// HasMagic reports whether data begins with the Ferret artifact magic bytes.
// It is intended for sniffing likely Ferret artifacts only. A true result does
// not guarantee that data is complete, valid, loadable, or compatible;
// callers that need actual artifact validation must still use Load or
// Unmarshal.
func HasMagic(data []byte) bool {
	return len(data) >= len(magic) && bytes.Equal(data[:len(magic)], magic[:])
}
