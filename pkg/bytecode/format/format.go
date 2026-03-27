package format

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

// Format encodes and decodes raw bytecode program payloads without an artifact
// header.
type Format interface {
	Name() string
	Marshal(*bytecode.Program) ([]byte, error)
	Unmarshal([]byte) (*bytecode.Program, error)
}
