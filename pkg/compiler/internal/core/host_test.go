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

func TestHostFunctionTable_AssignsStableSignatureIDsAndReturnsCopy(t *testing.T) {
	Convey("HostFunctionTable should assign first-seen IDs per signature", t, func() {
		tab := core.NewHostFunctionTable()

		fn3 := tab.Bind("FN", 3)
		fn1 := tab.Bind("FN", 1)
		fn3Again := tab.Bind("FN", 3)

		fns := tab.All()
		So(fn3, ShouldEqual, 0)
		So(fn1, ShouldEqual, 1)
		So(fn3Again, ShouldEqual, fn3)
		So(fns, ShouldResemble, []bytecode.HostFunction{
			{Name: "FN", ArgCount: 3},
			{Name: "FN", ArgCount: 1},
		})

		fns[0].Name = "changed"
		fns = append(fns, bytecode.HostFunction{Name: "NEW", ArgCount: 1})

		updated := tab.All()
		So(updated, ShouldResemble, []bytecode.HostFunction{
			{Name: "FN", ArgCount: 3},
			{Name: "FN", ArgCount: 1},
		})
	})
}
