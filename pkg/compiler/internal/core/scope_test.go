package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestScopeProjection(t *testing.T) {
	Convey("ScopeProjection", t, func() {
		Convey("NewScopeProjection", func() {
			Convey("Should create a new scope projection", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				// Create some test variables
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(1)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(2)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				
				So(sp, ShouldNotBeNil)
			})
		})

		Convey(".EmitAsArray", func() {
			Convey("Should emit scope variables as array", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				// Create some test variables
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(1)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(2)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				
				// Should not panic
				So(func() { sp.EmitAsArray(dstReg) }, ShouldNotPanic)
				
				// Should have generated some instructions
				So(emitter.Size(), ShouldBeGreaterThan, 0)
			})

			Convey("Should handle empty scope", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				var scope []core.Variable
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				
				// Should not panic with empty scope
				So(func() { sp.EmitAsArray(dstReg) }, ShouldNotPanic)
				So(emitter.Size(), ShouldBeGreaterThan, 0)
			})
		})

		Convey(".EmitAsObject", func() {
			Convey("Should emit scope variables as object", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(1)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(2)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				
				// Should not panic
				So(func() { sp.EmitAsObject(dstReg) }, ShouldNotPanic)
				So(emitter.Size(), ShouldBeGreaterThan, 0)
			})

			Convey("Should handle empty scope efficiently", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				var scope []core.Variable
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				initialSize := emitter.Size()
				
				sp.EmitAsObject(dstReg)
				
				// Empty scope should generate minimal instructions
				So(emitter.Size(), ShouldEqual, initialSize+1) // Just one LoadObject instruction
			})

			Convey("Should emit key-value pairs for variables", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "testVar", Type: core.TypeString, Register: vm.Operand(1)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				
				sp.EmitAsObject(dstReg)
				
				// Should have generated instructions for key and value
				So(emitter.Size(), ShouldBeGreaterThan, 2)
				
				// Should have added the variable name as a constant
				constants := st.Constants()
				found := false
				for _, c := range constants {
					if str, ok := c.(runtime.String); ok && string(str) == "testVar" {
						found = true
						break
					}
				}
				So(found, ShouldBeTrue)
			})
		})

		Convey(".RestoreFromArray", func() {
			Convey("Should restore variables from array", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(0)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(0)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				
				// Should not panic
				So(func() { sp.RestoreFromArray(srcReg) }, ShouldNotPanic)
				
				// Should have generated instructions
				So(emitter.Size(), ShouldBeGreaterThan, 0)
				
				// Should have declared local variables
				locals := st.LocalVariables()
				So(len(locals), ShouldBeGreaterThan, 0)
			})

			Convey("Should handle empty scope", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				var scope []core.Variable
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				initialSize := emitter.Size()
				
				sp.RestoreFromArray(srcReg)
				
				// Empty scope should not generate many instructions
				So(emitter.Size(), ShouldEqual, initialSize)
			})

			Convey("Should create indices as constants", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(0)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(0)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				
				sp.RestoreFromArray(srcReg)
				
				// Should have added indices as constants
				constants := st.Constants()
				foundIndices := 0
				for _, c := range constants {
					if intVal, ok := c.(runtime.Int); ok && (intVal == 0 || intVal == 1) {
						foundIndices++
					}
				}
				So(foundIndices, ShouldBeGreaterThanOrEqualTo, 2)
			})
		})

		Convey(".RestoreFromObject", func() {
			Convey("Should restore variables from object", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(0)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(0)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				
				// Should not panic
				So(func() { sp.RestoreFromObject(srcReg) }, ShouldNotPanic)
				
				// Should have generated instructions
				So(emitter.Size(), ShouldBeGreaterThan, 0)
				
				// Should have declared local variables
				locals := st.LocalVariables()
				So(len(locals), ShouldBeGreaterThan, 0)
			})

			Convey("Should handle empty scope", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				var scope []core.Variable
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				initialSize := emitter.Size()
				
				sp.RestoreFromObject(srcReg)
				
				// Empty scope should not generate many instructions
				So(emitter.Size(), ShouldEqual, initialSize)
			})

			Convey("Should create keys as constants", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "testVar1", Type: core.TypeString, Register: vm.Operand(0)},
					{Name: "testVar2", Type: core.TypeInt, Register: vm.Operand(0)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				srcReg := ra.Allocate(core.Temp)
				
				sp.RestoreFromObject(srcReg)
				
				// Should have added variable names as constants
				constants := st.Constants()
				foundKeys := 0
				for _, c := range constants {
					if str, ok := c.(runtime.String); ok {
						if string(str) == "testVar1" || string(str) == "testVar2" {
							foundKeys++
						}
					}
				}
				So(foundKeys, ShouldEqual, 2)
			})
		})

		Convey("Integration", func() {
			Convey("Should handle complete projection cycle", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				// Original scope
				scope := []core.Variable{
					{Name: "var1", Type: core.TypeString, Register: vm.Operand(1)},
					{Name: "var2", Type: core.TypeInt, Register: vm.Operand(2)},
					{Name: "var3", Type: core.TypeFloat, Register: vm.Operand(3)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				
				// Test array projection and restoration
				arrayReg := ra.Allocate(core.Temp)
				sp.EmitAsArray(arrayReg)
				
				// Enter new scope
				st.EnterScope()
				sp.RestoreFromArray(arrayReg)
				
				// Verify variables were restored
				locals := st.LocalVariables()
				So(len(locals), ShouldEqual, 3)
				
				// Test object projection and restoration
				objectReg := ra.Allocate(core.Temp)
				sp.EmitAsObject(objectReg)
				
				// Enter another new scope
				st.EnterScope()
				sp.RestoreFromObject(objectReg)
				
				// Verify variables were restored again
				locals2 := st.LocalVariables()
				So(len(locals2), ShouldEqual, 3)
			})

			Convey("Should handle single variable scope", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "singleVar", Type: core.TypeBool, Register: vm.Operand(1)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				
				// Test all operations
				dstReg := ra.Allocate(core.Temp)
				
				So(func() { sp.EmitAsArray(dstReg) }, ShouldNotPanic)
				So(func() { sp.EmitAsObject(dstReg) }, ShouldNotPanic)
				So(func() { sp.RestoreFromArray(dstReg) }, ShouldNotPanic)
				So(func() { sp.RestoreFromObject(dstReg) }, ShouldNotPanic)
			})

			Convey("Should handle various variable types", func() {
				ra := core.NewRegisterAllocator()
				emitter := core.NewEmitter()
				st := core.NewSymbolTable(ra)
				
				scope := []core.Variable{
					{Name: "strVar", Type: core.TypeString, Register: vm.Operand(1)},
					{Name: "intVar", Type: core.TypeInt, Register: vm.Operand(2)},
					{Name: "floatVar", Type: core.TypeFloat, Register: vm.Operand(3)},
					{Name: "boolVar", Type: core.TypeBool, Register: vm.Operand(4)},
					{Name: "listVar", Type: core.TypeList, Register: vm.Operand(5)},
					{Name: "mapVar", Type: core.TypeMap, Register: vm.Operand(6)},
					{Name: "anyVar", Type: core.TypeAny, Register: vm.Operand(7)},
					{Name: "unknownVar", Type: core.TypeUnknown, Register: vm.Operand(8)},
				}
				
				sp := core.NewScopeProjection(ra, emitter, st, scope)
				dstReg := ra.Allocate(core.Temp)
				
				// All operations should work regardless of variable types
				So(func() { sp.EmitAsArray(dstReg) }, ShouldNotPanic)
				So(func() { sp.EmitAsObject(dstReg) }, ShouldNotPanic)
				So(func() { sp.RestoreFromArray(dstReg) }, ShouldNotPanic)
				So(func() { sp.RestoreFromObject(dstReg) }, ShouldNotPanic)
			})
		})
	})
}