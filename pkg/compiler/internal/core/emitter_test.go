package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestEmitter(t *testing.T) {
	Convey("Emitter", t, func() {
		Convey("NewEmitter", func() {
			Convey("Should create a new emitter", func() {
				emitter := core.NewEmitter()
				
				So(emitter, ShouldNotBeNil)
				So(emitter.Size(), ShouldEqual, 0)
				So(emitter.Bytecode(), ShouldHaveLength, 0)
				So(emitter.Labels(), ShouldBeNil)
			})
		})

		Convey(".Size", func() {
			Convey("Should return the number of instructions", func() {
				emitter := core.NewEmitter()
				
				So(emitter.Size(), ShouldEqual, 0)
				
				emitter.Emit(vm.OpReturn)
				So(emitter.Size(), ShouldEqual, 1)
				
				emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
				So(emitter.Size(), ShouldEqual, 2)
			})
		})

		Convey(".Bytecode", func() {
			Convey("Should return all instructions", func() {
				emitter := core.NewEmitter()
				
				emitter.Emit(vm.OpReturn)
				emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
				
				bytecode := emitter.Bytecode()
				So(bytecode, ShouldHaveLength, 2)
				So(bytecode[0].Opcode, ShouldEqual, vm.OpReturn)
				So(bytecode[1].Opcode, ShouldEqual, vm.OpLoadNone)
				So(bytecode[1].Operands[0], ShouldEqual, vm.Operand(1))
			})
		})

		Convey(".Position", func() {
			Convey("Should return current position", func() {
				emitter := core.NewEmitter()
				
				// Position should be -1 when no instructions
				So(emitter.Position(), ShouldEqual, -1)
				
				emitter.Emit(vm.OpReturn)
				So(emitter.Position(), ShouldEqual, 0)
				
				emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
				So(emitter.Position(), ShouldEqual, 1)
			})
		})

		Convey("Label Management", func() {
			Convey(".NewLabel", func() {
				Convey("Should create new labels with unique IDs", func() {
					emitter := core.NewEmitter()
					
					label1 := emitter.NewLabel("test1")
					label2 := emitter.NewLabel("test2")
					
					So(label1, ShouldNotEqual, label2)
				})

				Convey("Should create unnamed labels", func() {
					emitter := core.NewEmitter()
					
					label := emitter.NewLabel()
					
					So(label, ShouldNotBeNil)
				})

				Convey("Should handle multiple name parts", func() {
					emitter := core.NewEmitter()
					
					label := emitter.NewLabel("part1", "part2", "part3")
					
					So(label, ShouldNotBeNil)
				})
			})

			Convey(".MarkLabel", func() {
				Convey("Should mark label at current position", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("test")
					
					emitter.Emit(vm.OpReturn)
					emitter.MarkLabel(label)
					
					pos, found := emitter.LabelPosition(label)
					So(found, ShouldBeTrue)
					So(pos, ShouldEqual, 1) // Position after the OpReturn
				})

				Convey("Should create labels map", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("test")
					
					emitter.MarkLabel(label)
					
					labels := emitter.Labels()
					So(labels, ShouldNotBeNil)
					So(len(labels), ShouldEqual, 1)
				})
			})

			Convey(".LabelPosition", func() {
				Convey("Should return false for unmarked labels", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("test")
					
					_, found := emitter.LabelPosition(label)
					So(found, ShouldBeFalse)
				})

				Convey("Should return correct position for marked labels", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("test")
					
					emitter.Emit(vm.OpReturn)
					emitter.Emit(vm.OpReturn)
					emitter.MarkLabel(label)
					
					pos, found := emitter.LabelPosition(label)
					So(found, ShouldBeTrue)
					So(pos, ShouldEqual, 2)
				})
			})
		})

		Convey("Instruction Emission", func() {
			Convey(".Emit", func() {
				Convey("Should emit instruction with no arguments", func() {
					emitter := core.NewEmitter()
					
					emitter.Emit(vm.OpReturn)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpReturn)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(0))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(0))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitA", func() {
				Convey("Should emit instruction with one operand", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitA(vm.OpLoadNone, vm.Operand(5))
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadNone)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(5))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(0))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitAB", func() {
				Convey("Should emit instruction with two operands", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAB(vm.OpMove, vm.Operand(1), vm.Operand(2))
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpMove)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitAb", func() {
				Convey("Should emit instruction with boolean argument true", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAb(vm.OpLoadBool, vm.Operand(1), true)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadBool)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})

				Convey("Should emit instruction with boolean argument false", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAb(vm.OpLoadBool, vm.Operand(1), false)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadBool)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(0))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitAx", func() {
				Convey("Should emit instruction with integer argument", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAx(vm.OpLoadConst, vm.Operand(1), 42)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadConst)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(42))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitAxy", func() {
				Convey("Should emit instruction with two integer arguments", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAxy(vm.OpCall, vm.Operand(1), 2, 3)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpCall)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(3))
				})
			})

			Convey(".EmitAs", func() {
				Convey("Should emit instruction with register sequence", func() {
					emitter := core.NewEmitter()
					seq := core.RegisterSequence{vm.Operand(5), vm.Operand(6), vm.Operand(7)}
					
					emitter.EmitAs(vm.OpCall, vm.Operand(1), seq)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpCall)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(5)) // First register
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(7)) // Last register
				})

				Convey("Should handle nil sequence", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitAs(vm.OpCall, vm.Operand(1), nil)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpCall)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(0))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(0))
				})
			})

			Convey(".EmitABx", func() {
				Convey("Should emit instruction with operand and integer", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitABx(vm.OpLoadIndex, vm.Operand(1), vm.Operand(2), 42)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadIndex)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(42))
				})
			})

			Convey(".EmitABC", func() {
				Convey("Should emit instruction with three operands", func() {
					emitter := core.NewEmitter()
					
					emitter.EmitABC(vm.OpAdd, vm.Operand(1), vm.Operand(2), vm.Operand(3))
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpAdd)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(3))
				})
			})
		})

		Convey("High-Level Emit Functions", func() {
			Convey("Should have convenience methods", func() {
				emitter := core.NewEmitter()
				ra := core.NewRegisterAllocator()
				
				// Test various high-level emit functions if they exist
				So(func() {
					// These functions should exist based on the scope.go usage
					dst := ra.Allocate(core.Temp)
					args := ra.AllocateSequence(2)
					
					// Test functions that should exist
					emitter.EmitArray(dst, args)
					emitter.EmitObject(dst, args)
					emitter.EmitMove(dst, args[0])
					
					st := core.NewSymbolTable(ra)
					constOp := st.AddConstant(runtime.NewString("test"))
					emitter.EmitLoadConst(dst, constOp)
				}, ShouldNotPanic)
			})
		})

		Convey("Instruction Modification", func() {
			Convey(".SwapAB", func() {
				Convey("Should modify instruction at label position", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("swap")
					
					// Emit a placeholder instruction
					emitter.MarkLabel(label)
					emitter.EmitAB(vm.OpMove, vm.Operand(1), vm.Operand(2))
					
					// Swap with different instruction
					emitter.SwapAB(label, vm.OpLoadBool, vm.Operand(3), vm.Operand(4))
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadBool)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(3))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(4))
				})
			})

			Convey(".SwapAx", func() {
				Convey("Should modify instruction with integer argument", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("swap")
					
					emitter.MarkLabel(label)
					emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
					
					emitter.SwapAx(label, vm.OpLoadConst, vm.Operand(2), 42)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadConst)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(42))
				})
			})

			Convey(".SwapAxy", func() {
				Convey("Should modify instruction with two integer arguments", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("swap")
					
					emitter.MarkLabel(label)
					emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
					
					emitter.SwapAxy(label, vm.OpCall, vm.Operand(1), 2, 3)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpCall)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(3))
				})
			})

			Convey(".SwapAs", func() {
				Convey("Should modify instruction with register sequence", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("swap")
					seq := core.RegisterSequence{vm.Operand(5), vm.Operand(6), vm.Operand(7)}
					
					emitter.MarkLabel(label)
					emitter.EmitA(vm.OpLoadNone, vm.Operand(1))
					
					emitter.SwapAs(label, vm.OpCall, vm.Operand(1), seq)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 1)
					So(bytecode[0].Opcode, ShouldEqual, vm.OpCall)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(1))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(5))
					So(bytecode[0].Operands[2], ShouldEqual, vm.Operand(7))
				})
			})

			Convey(".InsertAx", func() {
				Convey("Should insert instruction at label position", func() {
					emitter := core.NewEmitter()
					label := emitter.NewLabel("insert")
					
					emitter.MarkLabel(label)
					emitter.EmitA(vm.OpLoadBool, vm.Operand(1))
					
					// Insert should add instruction and shift others
					emitter.InsertAx(label, vm.OpLoadConst, vm.Operand(2), 42)
					
					bytecode := emitter.Bytecode()
					So(bytecode, ShouldHaveLength, 2)
					
					// Inserted instruction should be first
					So(bytecode[0].Opcode, ShouldEqual, vm.OpLoadConst)
					So(bytecode[0].Operands[0], ShouldEqual, vm.Operand(2))
					So(bytecode[0].Operands[1], ShouldEqual, vm.Operand(42))
					
					// Original instruction should be second
					So(bytecode[1].Opcode, ShouldEqual, vm.OpLoadBool)
					So(bytecode[1].Operands[0], ShouldEqual, vm.Operand(1))
				})
			})
		})

		Convey("Integration", func() {
			Convey("Should handle complex instruction sequences", func() {
				emitter := core.NewEmitter()
				
				// Create labels
				startLabel := emitter.NewLabel("start")
				endLabel := emitter.NewLabel("end")
				
				// Mark start
				emitter.MarkLabel(startLabel)
				
				// Emit various instructions
				emitter.Emit(vm.OpReturn)
				emitter.EmitA(vm.OpLoadBool, vm.Operand(1))
				emitter.EmitAB(vm.OpMove, vm.Operand(2), vm.Operand(1))
				emitter.EmitABC(vm.OpAdd, vm.Operand(3), vm.Operand(1), vm.Operand(2))
				
				// Mark end
				emitter.MarkLabel(endLabel)
				
				// Verify
				So(emitter.Size(), ShouldEqual, 4)
				
				startPos, found1 := emitter.LabelPosition(startLabel)
				endPos, found2 := emitter.LabelPosition(endLabel)
				
				So(found1, ShouldBeTrue)
				So(found2, ShouldBeTrue)
				So(startPos, ShouldEqual, 0)
				So(endPos, ShouldEqual, 4)
			})

			Convey("Should handle label patching", func() {
				emitter := core.NewEmitter()
				jumpLabel := emitter.NewLabel("jump_target")
				
				// Emit some instructions
				emitter.EmitA(vm.OpLoadBool, vm.Operand(1))
				
				// Later mark the label - this should trigger patching
				emitter.MarkLabel(jumpLabel)
				emitter.EmitA(vm.OpLoadNone, vm.Operand(2))
				
				pos, found := emitter.LabelPosition(jumpLabel)
				So(found, ShouldBeTrue)
				So(pos, ShouldEqual, 1) // After first instruction
			})
		})
	})
}