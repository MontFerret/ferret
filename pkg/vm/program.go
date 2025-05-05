package vm

import (
	"bytes"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime"
	"io"
	"text/tabwriter"
)

type (
	Catch [3]int

	Program struct {
		Source     *runtime.Source
		Locations  []runtime.Location
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		Params     []string
		Registers  int
	}
)

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
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	switch opcode {
	case OpMove:
		out.Write([]byte(fmt.Sprintf("MOVE %s %s", dst, src1)))
	case OpLoadNone:
		out.Write([]byte(fmt.Sprintf("LOADN %s", dst)))
	case OpLoadBool:
		out.Write([]byte(fmt.Sprintf("LOADB %s %d", dst, src1)))
	case OpLoadConst:
		out.Write([]byte(fmt.Sprintf("LOADC %s %s", dst, src1)))
	case OpLoadGlobal:
		out.Write([]byte(fmt.Sprintf("LOADG %s %s", dst, src1)))
	case OpStoreGlobal:
		out.Write([]byte(fmt.Sprintf("STOREG %s %s", dst, src1)))
	case OpCall:
		if src1 == 0 {
			out.Write([]byte(fmt.Sprintf("CALL %s", dst)))
		} else {
			out.Write([]byte(fmt.Sprintf("CALL %s %s %s", dst, src1, src2)))
		}
	case OpReturn:
		out.Write([]byte(fmt.Sprintf("RET")))
	default:
		return
	}
}
