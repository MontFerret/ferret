package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

func TestHostParamTable_StableSlots(t *testing.T) {
	Convey("HostParamTable should assign stable 1-based slots in first-seen order", t, func() {
		tab := core.NewHostParamTable()

		foo := tab.Bind("foo")
		bar := tab.Bind("bar")
		fooAgain := tab.Bind("foo")

		So(foo, ShouldEqual, bytecode.Operand(1))
		So(bar, ShouldEqual, bytecode.Operand(2))
		So(fooAgain, ShouldEqual, foo)

		So(tab.Names(), ShouldResemble, []string{"foo", "bar"})

		names := tab.Names()
		names[0] = "changed"
		So(tab.Names(), ShouldResemble, []string{"foo", "bar"})
	})
}

func TestHostFunctionTable_KeepsMaxArityAndReturnsCopy(t *testing.T) {
	Convey("HostFunctionTable should keep maximum arity and return a defensive copy", t, func() {
		tab := core.NewHostFunctionTable()

		tab.Bind("FN", 3)
		tab.Bind("FN", 1)
		tab.Bind("FN", 5)

		fns := tab.All()
		So(fns["FN"], ShouldEqual, 5)

		fns["FN"] = 0
		fns["NEW"] = 1

		updated := tab.All()
		So(updated["FN"], ShouldEqual, 5)
		_, exists := updated["NEW"]
		So(exists, ShouldBeFalse)
	})
}
