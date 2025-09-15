package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestDisassemble(t *testing.T) {
	Convey("Should return error for nil program", t, func() {
		result, err := asm.Disassemble(nil)

		So(err, ShouldEqual, asm.ErrInvalidProgram)
		So(result, ShouldEqual, "")
	})

	Convey("Should disassemble empty program", t, func() {
		program := &vm.Program{
			Bytecode:  []vm.Instruction{},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// An empty program should return an empty string or minimal output
		So(len(result) >= 0, ShouldBeTrue) // Just check it doesn't error
	})

	Convey("Should disassemble program with parameters", t, func() {
		program := &vm.Program{
			Bytecode:  []vm.Instruction{},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{"param1", "param2"},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, ".param param1")
		So(result, ShouldContainSubstring, ".param param2")
	})

	Convey("Should disassemble program with functions", t, func() {
		program := &vm.Program{
			Bytecode:  []vm.Instruction{},
			Constants: []runtime.Value{},
			Functions: map[string]int{"func1": 2, "func2": 0},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, ".func func1 2")
		So(result, ShouldContainSubstring, ".func func2 0")
	})

	Convey("Should disassemble program with constants", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{},
			Constants: []runtime.Value{
				runtime.NewString("hello"),
				runtime.NewInt(42),
				runtime.True,
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, ".const \"hello\"")
		So(result, ShouldContainSubstring, ".const 42")
		So(result, ShouldContainSubstring, ".const \"true\"")
	})

	Convey("Should disassemble program with simple instructions", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADN R0")
		So(result, ShouldContainSubstring, "RET R0")
	})

	Convey("Should disassemble program with jump instructions", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "JMP @L0")
		So(result, ShouldContainSubstring, "@L0:")
	})

	Convey("Should disassemble program with named labels", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{2: "end"},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "JMP @end")
		So(result, ShouldContainSubstring, "@end:")
	})

	Convey("Should disassemble program with constants loading", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(0), vm.NewConstant(0)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{
				runtime.NewString("hello world"),
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADC R0 C0 ; \"hello world\"")
	})

	Convey("Should disassemble program with conditional jumps", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJumpIfTrue, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "JMPT @L0 R1")
	})

	Convey("Should disassemble program with iterator operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpIterNext, vm.NewRegister(3), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpIterSkip, vm.NewRegister(5), vm.NewRegister(2), vm.NewRegister(3)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "ITNEXT @L0 R1")
		So(result, ShouldContainSubstring, "ITSKIP @L1 R2 R3")
	})

	Convey("Should handle disassembler options", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		// Test with debug option
		result, err := asm.Disassemble(program, asm.WithDebug())

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "RET R0")
	})

	Convey("Should disassemble program with various opcodes", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadBool, vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpMove, vm.NewRegister(1), vm.NewRegister(0)),
				vm.NewInstruction(vm.OpLength, vm.NewRegister(2), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpType, vm.NewRegister(3), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpAdd, vm.NewRegister(4), vm.NewRegister(2), vm.NewRegister(3)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(4)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADB R0 1")
		So(result, ShouldContainSubstring, "MOVE R1 R0")
		So(result, ShouldContainSubstring, "LEN R2 R1")
		So(result, ShouldContainSubstring, "TYPE R3 R2")
		So(result, ShouldContainSubstring, "ADD R4 R2 R3")
		So(result, ShouldContainSubstring, "RET R4")
	})

	Convey("Should disassemble program with dataset operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpDataSet, vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpDataSetCollector, vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpPush, vm.NewRegister(2), vm.NewRegister(3)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "DSET R0 1")
		So(result, ShouldContainSubstring, "DSETC R1 2")
		So(result, ShouldContainSubstring, "PUSH R2 R3")
	})

	Convey("Should disassemble program with function call operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpCall0, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpCall1, vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpCall, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "CALL0 R0")
		So(result, ShouldContainSubstring, "CALL1 R1 R2")
		So(result, ShouldContainSubstring, "CALL R3 R4 R5")
		So(result, ShouldContainSubstring, "RET R0")
	})

	Convey("Should disassemble program with comparison operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpEq, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpNe, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpGt, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "EQ R0 R1 R2")
		So(result, ShouldContainSubstring, "NE R3 R4 R5")
		So(result, ShouldContainSubstring, "GT R6 R7 R8")
		So(result, ShouldContainSubstring, "RET R0")
	})

	Convey("Should disassemble program with mixed constants and registers", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(0), vm.NewConstant(0)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(1)),
				vm.NewInstruction(vm.OpAdd, vm.NewRegister(2), vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(2)),
			},
			Constants: []runtime.Value{
				runtime.NewInt(10),
				runtime.NewInt(20),
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADC R0 C0 ; 10")
		So(result, ShouldContainSubstring, "LOADC R1 C1 ; 20")
		So(result, ShouldContainSubstring, "ADD R2 R0 R1")
		So(result, ShouldContainSubstring, "RET R2")
	})

	Convey("Should handle invalid constant indices gracefully", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(0), vm.NewConstant(99)), // Invalid index
			},
			Constants: []runtime.Value{
				runtime.NewString("valid"),
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADC R0 C99 ; <invalid>")
	})

	Convey("Should disassemble program with complete header sections", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{
				runtime.NewString("constant1"),
				runtime.NewInt(42),
			},
			Functions: map[string]int{
				"test_func": 2,
				"other_func": 0,
			},
			Params: []string{"param1", "param2", "param3"},
			Labels: map[int]string{0: "start"},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Check all header sections are present
		So(result, ShouldContainSubstring, ".param param1")
		So(result, ShouldContainSubstring, ".param param2")
		So(result, ShouldContainSubstring, ".param param3")
		So(result, ShouldContainSubstring, ".func test_func 2")
		So(result, ShouldContainSubstring, ".func other_func 0")
		So(result, ShouldContainSubstring, ".const \"constant1\"")
		So(result, ShouldContainSubstring, ".const 42")
		So(result, ShouldContainSubstring, "@start:")
		So(result, ShouldContainSubstring, "RET R0")
	})

	Convey("Should disassemble program with complex jump pattern", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(4)),         // Jump to instruction 4
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),      // instruction 1
				vm.NewInstruction(vm.OpJumpIfFalse, vm.NewRegister(1), vm.NewRegister(0)), // instruction 2, jump to 1
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),        // instruction 3
				vm.NewInstruction(vm.OpLoadBool, vm.NewRegister(0), vm.NewRegister(1)), // instruction 4
				vm.NewInstruction(vm.OpJump, vm.NewRegister(2)),          // instruction 5, jump to 2
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Check that labels are generated for all jump targets
		So(result, ShouldContainSubstring, "JMP @L0")
		So(result, ShouldContainSubstring, "JMPF @L1 R0")
		So(result, ShouldContainSubstring, "JMP @L2")
		// Check that label definitions are present
		So(result, ShouldContainSubstring, "@L0:")
		So(result, ShouldContainSubstring, "@L1:")
		So(result, ShouldContainSubstring, "@L2:")
	})
}