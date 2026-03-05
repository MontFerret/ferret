package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func newSymbolTable() *core.SymbolTable {
	return core.NewSymbolTable(core.NewRegisterAllocator(), nil)
}

func TestSymbolTable_AddConstant(t *testing.T) {
	Convey("SymbolTable should deduplicate scalar constants and runtime.None", t, func() {
		st := newSymbolTable()

		tests := []struct {
			name  string
			value runtime.Value
		}{
			{name: "string", value: runtime.NewString("ferret")},
			{name: "int", value: runtime.NewInt(42)},
			{name: "boolean", value: runtime.True},
			{name: "none", value: runtime.None},
		}

		for _, tt := range tests {
			first := st.AddConstant(tt.value)
			second := st.AddConstant(tt.value)

			So(first, ShouldEqual, second)
			So(st.Constant(first), ShouldEqual, tt.value)
		}

		So(len(st.Constants()), ShouldEqual, len(tests))
	})
}

func TestSymbolTable_AddVariable(t *testing.T) {
	Convey("SymbolTable should reject duplicates in the same scope and allow shadowing", t, func() {
		st := newSymbolTable()

		first, ok := st.DeclareLocal("item", core.TypeInt)
		So(ok, ShouldBeTrue)
		So(first.IsRegister(), ShouldBeTrue)

		duplicate, ok := st.DeclareLocal("item", core.TypeInt)
		So(ok, ShouldBeFalse)
		So(duplicate, ShouldEqual, bytecode.NoopOperand)

		resolved, kind, found := st.Resolve("item")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, first)

		st.EnterScope()
		shadow, ok := st.DeclareLocal("item", core.TypeString)
		So(ok, ShouldBeTrue)

		resolved, kind, found = st.Resolve("item")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, shadow)

		st.ExitScope()
		resolved, kind, found = st.Resolve("item")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, first)
	})
}

func TestSymbolTable_GetConstant(t *testing.T) {
	Convey("SymbolTable should return constant values by operand", t, func() {
		st := newSymbolTable()

		addr := st.AddConstant(runtime.NewString("value"))
		So(addr.IsConstant(), ShouldBeTrue)
		So(st.Constant(addr), ShouldEqual, runtime.NewString("value"))
	})
}

func TestSymbolTable_GetVariable(t *testing.T) {
	Convey("SymbolTable should resolve nearest local, then global, and report missing names", t, func() {
		st := newSymbolTable()

		globalReg := bytecode.NewRegister(100)
		So(st.AssignGlobal("target", core.TypeInt, globalReg), ShouldBeTrue)

		resolved, kind, found := st.Resolve("target")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolGlobal)
		So(resolved, ShouldEqual, globalReg)

		localReg := bytecode.NewRegister(10)
		So(st.AssignLocal("target", core.TypeString, localReg), ShouldBeTrue)

		resolved, kind, found = st.Resolve("target")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, localReg)

		st.EnterScope()
		innerReg := bytecode.NewRegister(11)
		So(st.AssignLocal("target", core.TypeBool, innerReg), ShouldBeTrue)

		resolved, kind, found = st.Resolve("target")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, innerReg)

		st.ExitScope()
		resolved, kind, found = st.Resolve("target")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, localReg)

		resolved, kind, found = st.Resolve("missing")
		So(found, ShouldBeFalse)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, bytecode.NoopOperand)
	})
}

func TestSymbolTable_GetVariableOrConstant(t *testing.T) {
	Convey("SymbolTable should enforce duplicate checks for AssignLocal and AssignGlobal", t, func() {
		st := newSymbolTable()

		firstLocal := bytecode.NewRegister(1)
		secondLocal := bytecode.NewRegister(2)
		So(st.AssignLocal("local", core.TypeInt, firstLocal), ShouldBeTrue)
		So(st.AssignLocal("local", core.TypeString, secondLocal), ShouldBeFalse)

		resolved, kind, found := st.Resolve("local")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, firstLocal)

		st.EnterScope()
		So(st.AssignLocal("local", core.TypeString, secondLocal), ShouldBeTrue)
		resolved, kind, found = st.Resolve("local")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolLocal)
		So(resolved, ShouldEqual, secondLocal)
		st.ExitScope()

		firstGlobal := bytecode.NewRegister(20)
		secondGlobal := bytecode.NewRegister(21)
		So(st.AssignGlobal("global", core.TypeInt, firstGlobal), ShouldBeTrue)
		So(st.AssignGlobal("global", core.TypeBool, secondGlobal), ShouldBeFalse)

		resolved, kind, found = st.Resolve("global")
		So(found, ShouldBeTrue)
		So(kind, ShouldEqual, core.SymbolGlobal)
		So(resolved, ShouldEqual, firstGlobal)
	})
}

func TestSymbolTable_HasConstant(t *testing.T) {
	Convey("SymbolTable should keep maximum function arity and return a defensive copy", t, func() {
		st := newSymbolTable()

		st.BindFunction("FN", 3)
		st.BindFunction("FN", 1)
		st.BindFunction("FN", 5)

		functions := st.Functions()
		So(functions["FN"], ShouldEqual, 5)

		functions["FN"] = 0
		functions["NEW"] = 1

		updated := st.Functions()
		So(updated["FN"], ShouldEqual, 5)
		_, exists := updated["NEW"]
		So(exists, ShouldBeFalse)
	})
}
