package base

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/asm"
	"github.com/MontFerret/ferret/pkg/vm"
)

func PrintDebugInfo(name string, uc TestCase, prog *vm.Program) {
	fmt.Println("")
	fmt.Println("VM Test:", name)
	fmt.Println("Expression:", uc.Expression)
	fmt.Println("")
	fmt.Println("Bytecode:")

	out, e := asm.Disassemble(prog, asm.WithDebug())

	if e == nil {
		fmt.Println(out)
	}
}
