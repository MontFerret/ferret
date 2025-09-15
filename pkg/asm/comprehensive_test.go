package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestDisassembleComprehensive(t *testing.T) {
	Convey("Should disassemble all arithmetic operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpAdd, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpSub, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpMulti, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpDiv, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
				vm.NewInstruction(vm.OpMod, vm.NewRegister(12), vm.NewRegister(13), vm.NewRegister(14)),
				vm.NewInstruction(vm.OpIncr, vm.NewRegister(15), vm.NewRegister(16), vm.NewRegister(17)),
				vm.NewInstruction(vm.OpDecr, vm.NewRegister(18), vm.NewRegister(19), vm.NewRegister(20)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "ADD R0 R1 R2")
		So(result, ShouldContainSubstring, "SUB R3 R4 R5")
		So(result, ShouldContainSubstring, "MUL R6 R7 R8")
		So(result, ShouldContainSubstring, "DIV R9 R10 R11")
		So(result, ShouldContainSubstring, "MOD R12 R13 R14")
		So(result, ShouldContainSubstring, "INCR R15 R16 R17")
		So(result, ShouldContainSubstring, "DECR R18 R19 R20")
	})

	Convey("Should disassemble all comparison operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpEq, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpNe, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpGt, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpLt, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
				vm.NewInstruction(vm.OpGte, vm.NewRegister(12), vm.NewRegister(13), vm.NewRegister(14)),
				vm.NewInstruction(vm.OpLte, vm.NewRegister(15), vm.NewRegister(16), vm.NewRegister(17)),
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
		So(result, ShouldContainSubstring, "LT R9 R10 R11")
		So(result, ShouldContainSubstring, "GTE R12 R13 R14")
		So(result, ShouldContainSubstring, "LTE R15 R16 R17")
	})

	Convey("Should disassemble type operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpCastBool, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpNegate, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpFlipPositive, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpFlipNegative, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
				vm.NewInstruction(vm.OpNot, vm.NewRegister(12), vm.NewRegister(13), vm.NewRegister(14)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "CASTB R0 R1 R2")
		So(result, ShouldContainSubstring, "NEG R3 R4 R5")
		So(result, ShouldContainSubstring, "FPP R6 R7 R8")
		So(result, ShouldContainSubstring, "FPN R9 R10 R11")
		So(result, ShouldContainSubstring, "NOT R12 R13 R14")
	})

	Convey("Should disassemble collection operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadArray, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLoadObject, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpLoadRange, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpLoadIndex, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
				vm.NewInstruction(vm.OpLoadKey, vm.NewRegister(12), vm.NewRegister(13), vm.NewRegister(14)),
				vm.NewInstruction(vm.OpLoadProperty, vm.NewRegister(15), vm.NewRegister(16), vm.NewRegister(17)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADARR R0 R1 R2")
		So(result, ShouldContainSubstring, "LOADOBJ R3 R4 R5")
		So(result, ShouldContainSubstring, "LOADRANGE R6 R7 R8")
		So(result, ShouldContainSubstring, "LOADI R9 R10 R11")
		So(result, ShouldContainSubstring, "LOADK R12 R13 R14")
		So(result, ShouldContainSubstring, "LOADPR R15 R16 R17")
	})

	Convey("Should disassemble optional collection operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadIndexOptional, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLoadKeyOptional, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpLoadPropertyOptional, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADIO R0 R1 R2")
		So(result, ShouldContainSubstring, "LOADKO R3 R4 R5")
		So(result, ShouldContainSubstring, "LOADPRO R6 R7 R8")
	})

	Convey("Should disassemble membership and pattern matching operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpIn, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLike, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpRegexp, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "IN R0 R1 R2")
		So(result, ShouldContainSubstring, "LIKE R3 R4 R5")
		So(result, ShouldContainSubstring, "REGEX R6 R7 R8")
	})

	Convey("Should disassemble array comparison operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpAnyEq, vm.NewRegister(0), vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpAnyNe, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpNoneEq, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpAllEq, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "ANYEQ R0 R1 R2")
		So(result, ShouldContainSubstring, "ANYNE R3 R4 R5")
		So(result, ShouldContainSubstring, "NONEQ R6 R7 R8")
		So(result, ShouldContainSubstring, "ALLEQ R9 R10 R11")
	})

	Convey("Should disassemble all function call variations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpCall0, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpCall1, vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpCall2, vm.NewRegister(3), vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpCall3, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpCall4, vm.NewRegister(9), vm.NewRegister(10), vm.NewRegister(11)),
				vm.NewInstruction(vm.OpProtectedCall0, vm.NewRegister(12)),
				vm.NewInstruction(vm.OpProtectedCall1, vm.NewRegister(13), vm.NewRegister(14)),
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
		So(result, ShouldContainSubstring, "CALL2 R3 R4 R5")
		So(result, ShouldContainSubstring, "CALL3 R6 R7 R8")
		So(result, ShouldContainSubstring, "CALL4 R9 R10 R11")
		So(result, ShouldContainSubstring, "PCALL0 R12")
		So(result, ShouldContainSubstring, "PCALL1 R13 R14")
	})

	Convey("Should disassemble iterator operations comprehensively", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpIter, vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpIterNext, vm.NewRegister(4), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpIterValue, vm.NewRegister(3), vm.NewRegister(0)),
				vm.NewInstruction(vm.OpIterKey, vm.NewRegister(4), vm.NewRegister(0)),
				vm.NewInstruction(vm.OpIterLimit, vm.NewRegister(4), vm.NewRegister(5), vm.NewRegister(6)),
				vm.NewInstruction(vm.OpIterSkip, vm.NewRegister(7), vm.NewRegister(8), vm.NewRegister(9)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "ITER R0 R1")
		So(result, ShouldContainSubstring, "ITNEXT") // Just check opcode is present
		So(result, ShouldContainSubstring, "ITVAL R3 R0")
		So(result, ShouldContainSubstring, "ITKEY R4 R0")
		So(result, ShouldContainSubstring, "ITLIMIT") // Just check opcode is present
		So(result, ShouldContainSubstring, "ITSKIP") // Just check opcode is present
	})

	Convey("Should disassemble dataset operations comprehensively", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpDataSet, vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpDataSetCollector, vm.NewRegister(2), vm.NewRegister(3)),
				vm.NewInstruction(vm.OpDataSetSorter, vm.NewRegister(4), vm.NewRegister(5)),
				vm.NewInstruction(vm.OpDataSetMultiSorter, vm.NewRegister(6), vm.NewRegister(7), vm.NewRegister(8)),
				vm.NewInstruction(vm.OpPush, vm.NewRegister(9), vm.NewRegister(10)),
				vm.NewInstruction(vm.OpPushKV, vm.NewRegister(11), vm.NewRegister(12), vm.NewRegister(13)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "DSET R0 1")
		So(result, ShouldContainSubstring, "DSETC R2 3")
		So(result, ShouldContainSubstring, "DSETS R4 5")
		So(result, ShouldContainSubstring, "DSETMS R6 R7 R8")
		So(result, ShouldContainSubstring, "PUSH R9 R10")
		So(result, ShouldContainSubstring, "PUSHKV R11 R12 R13")
	})

	Convey("Should disassemble utility operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLength, vm.NewRegister(0), vm.NewRegister(1)),
				vm.NewInstruction(vm.OpType, vm.NewRegister(2), vm.NewRegister(3)),
				vm.NewInstruction(vm.OpClose, vm.NewRegister(4)),
				vm.NewInstruction(vm.OpSleep, vm.NewRegister(5)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LEN R0 R1")
		So(result, ShouldContainSubstring, "TYPE R2 R3")
		So(result, ShouldContainSubstring, "CLOSE R4")
		So(result, ShouldContainSubstring, "SLEEP R5")
	})

	Convey("Should disassemble loading operations", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpLoadBool, vm.NewRegister(1), vm.NewRegister(2)),
				vm.NewInstruction(vm.OpLoadZero, vm.NewRegister(3)),
				vm.NewInstruction(vm.OpLoadParam, vm.NewRegister(4), vm.NewRegister(5), vm.NewRegister(6)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADN R0")
		So(result, ShouldContainSubstring, "LOADB R1 2")
		So(result, ShouldContainSubstring, "LOADZ R3")
		So(result, ShouldContainSubstring, "LOADP R4 R5 R6")
	})

	Convey("Should disassemble with different constant types", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(0), vm.NewConstant(0)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(1)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(2), vm.NewConstant(2)),
				vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(3), vm.NewConstant(3)),
			},
			Constants: []runtime.Value{
				runtime.NewString("text"),
				runtime.NewInt(123),
				runtime.True,
				runtime.False,
			},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{},
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "LOADC R0 C0 ; \"text\"")
		So(result, ShouldContainSubstring, "LOADC R1 C1 ; 123")
		So(result, ShouldContainSubstring, "LOADC R2 C2 ; \"true\"")
		So(result, ShouldContainSubstring, "LOADC R3 C3 ; \"false\"")
	})

	Convey("Should handle edge case with empty labels map", t, func() {
		program := &vm.Program{
			Bytecode: []vm.Instruction{
				vm.NewInstruction(vm.OpJump, vm.NewRegister(2)),
				vm.NewInstruction(vm.OpReturn, vm.NewRegister(0)),
				vm.NewInstruction(vm.OpLoadNone, vm.NewRegister(0)),
			},
			Constants: []runtime.Value{},
			Functions: map[string]int{},
			Params:    []string{},
			Labels:    map[int]string{}, // Empty labels map
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		// Should generate automatic labels
		So(result, ShouldContainSubstring, "JMP @L0")
		So(result, ShouldContainSubstring, "@L0:")
	})

	Convey("Should handle edge case with nil maps in program", t, func() {
		program := &vm.Program{
			Bytecode:  []vm.Instruction{vm.NewInstruction(vm.OpReturn, vm.NewRegister(0))},
			Constants: []runtime.Value{},
			Functions: nil, // nil map
			Params:    []string{},
			Labels:    nil, // nil map
		}

		result, err := asm.Disassemble(program)

		So(err, ShouldBeNil)
		So(result, ShouldContainSubstring, "RET R0")
	})
}