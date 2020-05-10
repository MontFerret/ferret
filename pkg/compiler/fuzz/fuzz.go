package compiler_test

import (
	"github.com/MontFerret/ferret/pkg/compiler"
)

// Fuzz is our fuzzer.
// The fuzzer is run on oss-fuzz.
// Link to the project on oss-fuzz
// will be added once the project is up.
func Fuzz(data []byte) int {
	c := compiler.New()
	p, err := c.Compile(string(data))
	if err != nil || p == nil {
		return 0
	}
	return 1
}
