package runtime_v2

import (
	"bytes"
	"io"
	"text/tabwriter"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Program struct {
	Source    *core.Source
	Locations []core.Location
	Bytecode  []Opcode
	Arguments []int
	Constants []core.Value
}

func (program *Program) Disassemble() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	for offset := 0; offset < len(program.Bytecode); {
		instruction := program.Bytecode[offset]
		program.disassembleInstruction(w, instruction, offset)
		offset++
	}

	w.Flush()

	return buf.String()
}

func (program *Program) disassembleInstruction(out io.Writer, opcode Opcode, offset int) {
	switch opcode {
	case OpReturn:
		out.Write([]byte("return"))
	default:
		return
	}
}
