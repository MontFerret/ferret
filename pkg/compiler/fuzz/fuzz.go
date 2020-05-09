package compiler_test

import (
	"github.com/MontFerret/ferret/pkg/compiler"
)

func Fuzz(data []byte) int {
	c := compiler.New()
	p, err := c.Compile(string(data))
	if err!=nil || p==nil{
		return 0
	}
	return 1
}
