package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestSymbolTable(t *testing.T) {
	Convey("SymbolTable", t, func() {
		Convey("NewSymbolTable", func() {
			Convey("Should create a new symbol table", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				So(st, ShouldNotBeNil)
				So(st.Scope(), ShouldEqual, 0)
				So(st.LocalVariables(), ShouldHaveLength, 0)
				So(st.Constants(), ShouldHaveLength, 0)
				So(st.Params(), ShouldHaveLength, 0)
			})
		})

		Convey("Scope Management", func() {
			Convey("Should handle scope enter and exit", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				So(st.Scope(), ShouldEqual, 0)
				
				st.EnterScope()
				So(st.Scope(), ShouldEqual, 1)
				
				st.EnterScope()
				So(st.Scope(), ShouldEqual, 2)
				
				st.ExitScope()
				So(st.Scope(), ShouldEqual, 1)
				
				st.ExitScope()
				So(st.Scope(), ShouldEqual, 0)
			})

			Convey("Should clean up local variables on exit scope", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				// Add variable in scope 0
				st.DeclareLocal("var0", core.TypeString)
				So(st.LocalVariables(), ShouldHaveLength, 1)
				
				// Enter scope 1 and add variable
				st.EnterScope()
				st.DeclareLocal("var1", core.TypeInt)
				So(st.LocalVariables(), ShouldHaveLength, 1) // Only current scope
				
				// Exit scope 1 - var1 should be cleaned up
				st.ExitScope()
				locals := st.LocalVariables()
				So(locals, ShouldHaveLength, 1)
				So(locals[0].Name, ShouldEqual, "var0")
			})
		})

		Convey("Local Variables", func() {
			Convey("Should declare local variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				op, ok := st.DeclareLocal("testVar", core.TypeString)
				
				So(ok, ShouldBeTrue)
				So(op, ShouldNotEqual, vm.NoopOperand)
				
				locals := st.LocalVariables()
				So(locals, ShouldHaveLength, 1)
				So(locals[0].Name, ShouldEqual, "testVar")
				So(locals[0].Type, ShouldEqual, core.TypeString)
				So(locals[0].Kind, ShouldEqual, core.SymbolVar)
			})

			Convey("Should not allow duplicate variables in same scope", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, ok1 := st.DeclareLocal("testVar", core.TypeString)
				_, ok2 := st.DeclareLocal("testVar", core.TypeInt)
				
				So(ok1, ShouldBeTrue)
				So(ok2, ShouldBeFalse)
				So(st.LocalVariables(), ShouldHaveLength, 1)
			})

			Convey("Should allow same name in different scopes", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, ok1 := st.DeclareLocal("testVar", core.TypeString)
				So(ok1, ShouldBeTrue)
				
				st.EnterScope()
				_, ok2 := st.DeclareLocal("testVar", core.TypeInt)
				So(ok2, ShouldBeTrue)
			})

			Convey("Should handle pseudo variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				// Should allow multiple ignore pseudo variables
				_, ok1 := st.DeclareLocal(core.IgnorePseudoVariable, core.TypeAny)
				_, ok2 := st.DeclareLocal(core.IgnorePseudoVariable, core.TypeAny)
				
				So(ok1, ShouldBeTrue)
				So(ok2, ShouldBeTrue)
				
				// Should allow multiple current pseudo variables
				_, ok3 := st.DeclareLocal(core.PseudoVariable, core.TypeAny)
				_, ok4 := st.DeclareLocal(core.PseudoVariable, core.TypeAny)
				
				So(ok3, ShouldBeTrue)
				So(ok4, ShouldBeTrue)
			})

			Convey("Should assign local variables to existing operands", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				reg := ra.Allocate(core.Var)
				ok := st.AssignLocal("testVar", core.TypeString, reg)
				
				So(ok, ShouldBeTrue)
				locals := st.LocalVariables()
				So(locals, ShouldHaveLength, 1)
				So(locals[0].Register, ShouldEqual, reg)
			})
		})

		Convey("Global Variables", func() {
			Convey("Should declare global variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				op, ok := st.DeclareGlobal("globalVar", core.TypeString)
				
				So(ok, ShouldBeTrue)
				So(op, ShouldNotEqual, vm.NoopOperand)
			})

			Convey("Should not allow duplicate global variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, ok1 := st.DeclareGlobal("globalVar", core.TypeString)
				_, ok2 := st.DeclareGlobal("globalVar", core.TypeInt)
				
				So(ok1, ShouldBeTrue)
				So(ok2, ShouldBeFalse)
			})

			Convey("Should assign global variables to existing operands", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				reg := ra.Allocate(core.Var)
				ok := st.AssignGlobal("globalVar", core.TypeString, reg)
				
				So(ok, ShouldBeTrue)
				
				// Try to assign again - should fail
				ok2 := st.AssignGlobal("globalVar", core.TypeInt, reg)
				So(ok2, ShouldBeFalse)
			})
		})

		Convey("Constants", func() {
			Convey("Should add constants", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				val := runtime.NewString("test")
				op := st.AddConstant(val)
				
				So(op.IsConstant(), ShouldBeTrue)
				So(st.Constants(), ShouldHaveLength, 1)
				So(st.Constants()[0], ShouldEqual, val)
			})

			Convey("Should retrieve constants by operand", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				val := runtime.NewString("test")
				op := st.AddConstant(val)
				
				retrieved := st.Constant(op)
				So(retrieved, ShouldEqual, val)
			})
		})

		Convey("Parameters", func() {
			Convey("Should bind parameters", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				op := st.BindParam("paramName")
				
				So(op.IsConstant(), ShouldBeTrue)
				params := st.Params()
				So(params, ShouldHaveLength, 1)
				So(params[0], ShouldEqual, "paramName")
			})

			Convey("Should handle multiple parameters", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.BindParam("param1")
				st.BindParam("param2")
				
				params := st.Params()
				So(params, ShouldHaveLength, 2)
				So(params, ShouldContain, "param1")
				So(params, ShouldContain, "param2")
			})
		})

		Convey("Functions", func() {
			Convey("Should bind functions", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.BindFunction("testFunc", 2)
				
				funcs := st.Functions()
				So(funcs, ShouldHaveLength, 1)
				So(funcs["testFunc"], ShouldEqual, 2)
			})

			Convey("Should not override function with fewer args", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.BindFunction("testFunc", 3)
				st.BindFunction("testFunc", 2) // Should not override
				
				funcs := st.Functions()
				So(funcs["testFunc"], ShouldEqual, 3)
			})

			Convey("Should override function with more args", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.BindFunction("testFunc", 2)
				st.BindFunction("testFunc", 3) // Should override
				
				funcs := st.Functions()
				So(funcs["testFunc"], ShouldEqual, 3)
			})
		})

		Convey("Symbol Resolution", func() {
			Convey("Should resolve local variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, _ = st.DeclareLocal("localVar", core.TypeString)
				
				resolved, kind, found := st.Resolve("localVar")
				So(found, ShouldBeTrue)
				So(kind, ShouldEqual, core.SymbolVar)
				So(resolved, ShouldNotEqual, vm.NoopOperand)
			})

			Convey("Should resolve global variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.DeclareGlobal("globalVar", core.TypeString)
				
				resolved, kind, found := st.Resolve("globalVar")
				So(found, ShouldBeTrue)
				So(kind, ShouldEqual, core.SymbolGlobal)
				So(resolved, ShouldNotEqual, vm.NoopOperand)
			})

			Convey("Should prioritize local over global", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.DeclareGlobal("var", core.TypeString)
				st.DeclareLocal("var", core.TypeInt)
				
				_, kind, found := st.Resolve("var")
				So(found, ShouldBeTrue)
				So(kind, ShouldEqual, core.SymbolVar) // Local should win
			})

			Convey("Should not find non-existent variable", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, _, found := st.Resolve("nonExistent")
				So(found, ShouldBeFalse)
			})

			Convey("Should find variables in correct scope", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.DeclareLocal("outerVar", core.TypeString)
				
				st.EnterScope()
				st.DeclareLocal("innerVar", core.TypeInt)
				
				// Both should be resolvable from inner scope
				_, _, found1 := st.Resolve("outerVar")
				_, _, found2 := st.Resolve("innerVar")
				So(found1, ShouldBeTrue)
				So(found2, ShouldBeTrue)
				
				st.ExitScope()
				
				// Only outer should be resolvable from outer scope
				_, _, found3 := st.Resolve("outerVar")
				_, _, found4 := st.Resolve("innerVar")
				So(found3, ShouldBeTrue)
				So(found4, ShouldBeFalse)
			})
		})

		Convey("Variable Lookup", func() {
			Convey("Should lookup local variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.DeclareLocal("testVar", core.TypeString)
				
				variable, found := st.Lookup("testVar")
				So(found, ShouldBeTrue)
				So(variable.Name, ShouldEqual, "testVar")
				So(variable.Type, ShouldEqual, core.TypeString)
			})

			Convey("Should not find non-existent variables", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				_, found := st.Lookup("nonExistent")
				So(found, ShouldBeFalse)
			})
		})

		Convey("Debug View", func() {
			Convey("Should provide debug information", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				st.DeclareLocal("localVar", core.TypeString)
				st.DeclareGlobal("globalVar", core.TypeInt)
				st.BindParam("param1")
				st.AddConstant(runtime.NewString("const"))
				
				debug := st.DebugView()
				So(debug, ShouldNotBeEmpty)
				So(len(debug), ShouldBeGreaterThan, 0)
			})
		})

		Convey("Integration", func() {
			Convey("Should handle complex symbol management", func() {
				ra := core.NewRegisterAllocator()
				st := core.NewSymbolTable(ra)
				
				// Add globals
				st.DeclareGlobal("global1", core.TypeString)
				st.DeclareGlobal("global2", core.TypeInt)
				
				// Add parameters
				st.BindParam("param1")
				st.BindParam("param2")
				
				// Add functions
				st.BindFunction("func1", 2)
				st.BindFunction("func2", 3)
				
				// Add constants
				st.AddConstant(runtime.NewString("const1"))
				st.AddConstant(runtime.NewInt(42))
				
				// Add locals in different scopes
				st.DeclareLocal("local1", core.TypeString)
				
				st.EnterScope()
				st.DeclareLocal("local2", core.TypeInt)
				
				// Verify everything is accessible
				_, _, found1 := st.Resolve("global1")
				_, _, found2 := st.Resolve("local1")
				_, _, found3 := st.Resolve("local2")
				
				So(found1, ShouldBeTrue)
				So(found2, ShouldBeTrue)
				So(found3, ShouldBeTrue)
				
				So(st.Params(), ShouldHaveLength, 2)
				So(st.Functions(), ShouldHaveLength, 2)
				So(st.Constants(), ShouldHaveLength, 4) // 2 explicit + 2 param names as constants
			})
		})
	})
}

// Legacy test stubs - keeping these for compatibility but implementing them
func TestSymbolTable_AddConstant(t *testing.T) {
	Convey("Should add constant", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		val := runtime.NewString("test")
		op := st.AddConstant(val)
		
		So(op.IsConstant(), ShouldBeTrue)
		So(st.Constant(op), ShouldEqual, val)
	})
}

func TestSymbolTable_AddVariable(t *testing.T) {
	Convey("Should add variable", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		op, ok := st.DeclareLocal("testVar", core.TypeString)
		
		So(ok, ShouldBeTrue)
		So(op, ShouldNotEqual, vm.NoopOperand)
	})
}

func TestSymbolTable_GetConstant(t *testing.T) {
	Convey("Should get constant", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		val := runtime.NewString("test")
		op := st.AddConstant(val)
		
		retrieved := st.Constant(op)
		So(retrieved, ShouldEqual, val)
	})
}

func TestSymbolTable_GetVariable(t *testing.T) {
	Convey("Should get variable", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		st.DeclareLocal("testVar", core.TypeString)
		
		variable, found := st.Lookup("testVar")
		So(found, ShouldBeTrue)
		So(variable.Name, ShouldEqual, "testVar")
	})
}

func TestSymbolTable_GetVariableOrConstant(t *testing.T) {
	Convey("Should resolve variables and constants", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		// Add a local variable
		st.DeclareLocal("testVar", core.TypeString)
		
		// Should resolve the variable
		_, _, found := st.Resolve("testVar")
		So(found, ShouldBeTrue)
	})
}

func TestSymbolTable_HasConstant(t *testing.T) {
	Convey("Should check for constants", t, func() {
		ra := core.NewRegisterAllocator()
		st := core.NewSymbolTable(ra)
		
		// Add a constant
		val := runtime.NewString("test")
		op := st.AddConstant(val)
		
		// Should be able to retrieve it
		retrieved := st.Constant(op)
		So(retrieved, ShouldEqual, val)
		
		// Constants should be in the constants list
		constants := st.Constants()
		So(constants, ShouldContain, val)
	})
}
