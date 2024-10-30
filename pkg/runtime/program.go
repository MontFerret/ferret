package runtime

import (
	"bytes"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"io"
	"text/tabwriter"
)

type Program struct {
	Source     *core.Source
	Locations  []core.Location
	Bytecode   []Instruction
	Constants  []core.Value
	CatchTable [][2]int
	Registers  int
}

func (program *Program) Disassemble() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	for offset := 0; offset < len(program.Bytecode); {
		instruction := program.Bytecode[offset]
		program.disassembleInstruction(w, instruction, offset)
		offset++
		w.Write([]byte("\n"))
	}

	w.Flush()

	return buf.String()
}

func (program *Program) disassembleInstruction(out io.Writer, inst Instruction, offset int) {
	opcode := inst.Opcode
	out.Write([]byte(fmt.Sprintf("%d: [%d] ", offset, opcode)))
	dst, src1 := inst.Operands[0], inst.Operands[1]
	//src2 := program.disassembleOperand(inst.Operands[2])

	switch opcode {
	case OpMove:
		out.Write([]byte(fmt.Sprintf("MOVE %s %s", dst, src1)))
	case OpLoadConst:
		out.Write([]byte(fmt.Sprintf("LOADK %s %s", dst, src1)))
	case OpLoadGlobal:
		out.Write([]byte(fmt.Sprintf("LOADG %s %s", dst, src1)))
	case OpStoreGlobal:
		out.Write([]byte(fmt.Sprintf("STOREG %s %s", dst, src1)))
	case OpReturn:
		out.Write([]byte(fmt.Sprintf("RET")))
	default:
		return
	}
}
