package events

import (
	"hash/fnv"
)

func New(name string) ID {
	h := fnv.New32a()

	h.Write([]byte(name))

	return ID(h.Sum32())
}
