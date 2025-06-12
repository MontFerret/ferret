package vm

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/MontFerret/ferret/pkg/runtime"
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

func (program *Program) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	var counter int

	for _, inst := range program.Bytecode {
		counter++
		program.writeInstruction(w, counter, inst)
		w.Write([]byte("\n"))
	}

	w.Flush()

	return buf.String()
}

func (program *Program) writeInstruction(w io.Writer, pos int, inst Instruction) {
	if inst.Opcode != OpReturn {
		w.Write([]byte(fmt.Sprintf("%d: %s", pos, inst)))
	} else {
		w.Write([]byte(fmt.Sprintf("%d: %s", pos, inst.Opcode)))
	}
}
