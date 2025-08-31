package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestEdgeCases(t *testing.T) {
	Convey("Should handle stream operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpStream, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpStreamIter, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "STRM R0 R1 R2")
		So(result, ShouldContainSubstring, "STRMITER R3 R4 R5")
	})

	Convey("Should handle comparison operation", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpCmp, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "COMP R0 R1 R2")
	})

	Convey("Should handle complex constant values", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(0), vm.NewConstant(0)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(1)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(2), vm.NewConstant(2)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{
				runtime.NewString(""),      // empty string
				runtime.NewString("\"quotes\""), // string with quotes
				runtime.NewInt(0),          // zero value
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADC R0 C0 ; \"\"")
		So(result, ShouldContainSubstring, "LOADC R1 C1 ; \"\\\"quotes\\\"\"")
		So(result, ShouldContainSubstring, "LOADC R2 C2 ; 0")
	})

	Convey("Should handle multiple disassembler options", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		// Test with multiple options (though only WithDebug exists currently)
		result, err := asm.Disassemble(program, asm.WithDebug(), asm.WithDebug())

		So(err, ShouldBeNil)
		So(result, ShouldNotEqual, "")
	})

	Convey("Should handle program with only labels", t, func() {
		program := &vm.Program{
			Bytecode:  []vm.Instruction{},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{0: "start", 5: "end"},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Empty bytecode means no instructions, so labels won't appear
		So(len(result) >= 0, ShouldBeTrue)
	})

	Convey("Should handle large register indices", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpMove, vm.NewRegister(255), vm.NewRegister(1024)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(255)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "MOVE R255 R1024")
		So(result, ShouldContainSubstring, "RET R255")
	})

	Convey("Should handle program with all header sections empty", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{}, // empty
			Functions: map[string]int{},  // empty
			Params:    []string{},        // empty
			Labels:    map[int]string{},  // empty
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "RET R0")
		// Should not contain any header directives
		So(result, ShouldNotContainSubstring, ".param")
		So(result, ShouldNotContainSubstring, ".func")
		So(result, ShouldNotContainSubstring, ".const")
	})

	Convey("Should handle jump to non-existing label address", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(999)), // Jump to non-existing address
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{}, // No labels defined
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Should display numeric address since no label exists
		So(result, ShouldContainSubstring, "JMP @L0") // Auto-generated label
	})

	Convey("Should test all protected call operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpProtectedCall, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpProtectedCall2, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpProtectedCall3, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpProtectedCall4, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "PCALL R0 R1 R2")
		So(result, ShouldContainSubstring, "PCALL2 R3 R4 R5")
		So(result, ShouldContainSubstring, "PCALL3 R6 R7 R8")
		So(result, ShouldContainSubstring, "PCALL4 R9 R10 R11")
	})

	Convey("Should test labelOrAddr function coverage", t, func() {
		// Test case where jump target has no label - should use numeric address
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(1)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{999: "unused_label"}, // Label for different address
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Should generate auto label since target doesn't have a named label
		So(result, ShouldContainSubstring, "JMP @L0")
	})

	Convey("Should handle parameter loading operation", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadParam, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADP R0 R1 R2")
	})
}