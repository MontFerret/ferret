package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/frame"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type execState struct {
	program   *bytecode.Program
	env       *Environment
	registers *mem.RegisterFile
	scratch   *mem.Scratch
	frames    frame.CallStack
	catchByPC []int
	pc        int
}

func (s *execState) init(program *bytecode.Program, catchByPC []int) {
	s.program = program
	s.catchByPC = catchByPC
	s.registers = mem.NewRegisterFile(program.Registers)
	s.scratch = mem.NewScratch(len(program.Params))
	s.frames.Init(maxUDFRegisters(program.Functions.UserDefined))
}

func (s *execState) reset(env *Environment) {
	if s.registers.IsDirty() {
		s.registers.Reset()
	}

	s.registers.MarkDirty()
	s.env = env
	s.pc = 0
	s.frames.Reset()
}
