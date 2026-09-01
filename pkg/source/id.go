package source

import (
	"crypto/sha256"
	"encoding/binary"
)

type ID [32]byte

func generateID(name, content string) ID {
	h := sha256.New()
	binary.Write(h, binary.LittleEndian, uint64(len(name)))
	h.Write([]byte(name))

	binary.Write(h, binary.LittleEndian, uint64(len(content)))
	h.Write([]byte(content))

	var id ID

	copy(id[:], h.Sum(nil))

	return id
}
