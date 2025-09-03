package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestKV(t *testing.T) {
	Convey("KV", t, func() {
		Convey("NewKV", func() {
			Convey("Should create a new key-value pair", func() {
				key := vm.Operand(1)
				value := vm.Operand(2)
				
				kv := core.NewKV(key, value)
				
				So(kv, ShouldNotBeNil)
				So(kv.Key, ShouldEqual, key)
				So(kv.Value, ShouldEqual, value)
			})

			Convey("Should handle zero operands", func() {
				key := vm.Operand(0)
				value := vm.Operand(0)
				
				kv := core.NewKV(key, value)
				
				So(kv, ShouldNotBeNil)
				So(kv.Key, ShouldEqual, key)
				So(kv.Value, ShouldEqual, value)
			})

			Convey("Should handle large operand values", func() {
				key := vm.Operand(999999)
				value := vm.Operand(888888)
				
				kv := core.NewKV(key, value)
				
				So(kv, ShouldNotBeNil)
				So(kv.Key, ShouldEqual, key)
				So(kv.Value, ShouldEqual, value)
			})
		})
	})
}

func TestCatchStack(t *testing.T) {
	Convey("CatchStack", t, func() {
		Convey("NewCatchStack", func() {
			Convey("Should create a new empty catch stack", func() {
				cs := core.NewCatchStack()
				
				So(cs, ShouldNotBeNil)
				So(cs.Len(), ShouldEqual, 0)
				So(cs.All(), ShouldHaveLength, 0)
			})
		})

		Convey(".Push", func() {
			Convey("Should add catch entries", func() {
				cs := core.NewCatchStack()
				
				cs.Push(10, 20, 30)
				
				So(cs.Len(), ShouldEqual, 1)
				entries := cs.All()
				So(entries, ShouldHaveLength, 1)
				So(entries[0][0], ShouldEqual, 10) // start
				So(entries[0][1], ShouldEqual, 20) // end
				So(entries[0][2], ShouldEqual, 30) // jump
			})

			Convey("Should handle multiple entries", func() {
				cs := core.NewCatchStack()
				
				cs.Push(10, 20, 30)
				cs.Push(40, 50, 60)
				cs.Push(70, 80, 90)
				
				So(cs.Len(), ShouldEqual, 3)
				entries := cs.All()
				So(entries, ShouldHaveLength, 3)
			})

			Convey("Should handle zero values", func() {
				cs := core.NewCatchStack()
				
				cs.Push(0, 0, 0)
				
				So(cs.Len(), ShouldEqual, 1)
				entries := cs.All()
				So(entries[0][0], ShouldEqual, 0)
				So(entries[0][1], ShouldEqual, 0)
				So(entries[0][2], ShouldEqual, 0)
			})
		})

		Convey(".Pop", func() {
			Convey("Should remove last entry", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				cs.Push(40, 50, 60)
				
				cs.Pop()
				
				So(cs.Len(), ShouldEqual, 1)
				entries := cs.All()
				So(entries[0][0], ShouldEqual, 10)
				So(entries[0][1], ShouldEqual, 20)
				So(entries[0][2], ShouldEqual, 30)
			})

			Convey("Should handle empty stack", func() {
				cs := core.NewCatchStack()
				
				// Should not panic
				So(func() { cs.Pop() }, ShouldNotPanic)
				So(cs.Len(), ShouldEqual, 0)
			})

			Convey("Should handle popping until empty", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				cs.Push(40, 50, 60)
				
				cs.Pop()
				cs.Pop()
				
				So(cs.Len(), ShouldEqual, 0)
				
				// Should not panic
				So(func() { cs.Pop() }, ShouldNotPanic)
			})
		})

		Convey(".Find", func() {
			Convey("Should find catch entry for position", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				cs.Push(25, 35, 40)
				
				// Position within first range
				catch, found := cs.Find(15)
				So(found, ShouldBeTrue)
				So(catch[0], ShouldEqual, 10)
				So(catch[1], ShouldEqual, 20)
				So(catch[2], ShouldEqual, 30)
				
				// Position within second range
				catch2, found2 := cs.Find(30)
				So(found2, ShouldBeTrue)
				So(catch2[0], ShouldEqual, 25)
				So(catch2[1], ShouldEqual, 35)
				So(catch2[2], ShouldEqual, 40)
			})

			Convey("Should handle boundary conditions", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				
				// Start boundary
				catch, found := cs.Find(10)
				So(found, ShouldBeTrue)
				So(catch[0], ShouldEqual, 10)
				
				// End boundary
				catch2, found2 := cs.Find(20)
				So(found2, ShouldBeTrue)
				So(catch2[0], ShouldEqual, 10)
			})

			Convey("Should not find catch for position outside ranges", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				
				// Before range
				_, found := cs.Find(5)
				So(found, ShouldBeFalse)
				
				// After range
				_, found2 := cs.Find(25)
				So(found2, ShouldBeFalse)
			})

			Convey("Should handle empty stack", func() {
				cs := core.NewCatchStack()
				
				_, found := cs.Find(10)
				So(found, ShouldBeFalse)
			})
		})

		Convey(".Clear", func() {
			Convey("Should remove all entries", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				cs.Push(40, 50, 60)
				cs.Push(70, 80, 90)
				
				cs.Clear()
				
				So(cs.Len(), ShouldEqual, 0)
				So(cs.All(), ShouldHaveLength, 0)
			})

			Convey("Should handle empty stack", func() {
				cs := core.NewCatchStack()
				
				So(func() { cs.Clear() }, ShouldNotPanic)
				So(cs.Len(), ShouldEqual, 0)
			})
		})

		Convey(".Len", func() {
			Convey("Should return correct length", func() {
				cs := core.NewCatchStack()
				
				So(cs.Len(), ShouldEqual, 0)
				
				cs.Push(10, 20, 30)
				So(cs.Len(), ShouldEqual, 1)
				
				cs.Push(40, 50, 60)
				So(cs.Len(), ShouldEqual, 2)
				
				cs.Pop()
				So(cs.Len(), ShouldEqual, 1)
			})
		})

		Convey(".All", func() {
			Convey("Should return all entries", func() {
				cs := core.NewCatchStack()
				cs.Push(10, 20, 30)
				cs.Push(40, 50, 60)
				
				all := cs.All()
				
				So(all, ShouldHaveLength, 2)
				So(all[0][0], ShouldEqual, 10)
				So(all[0][1], ShouldEqual, 20)
				So(all[0][2], ShouldEqual, 30)
				So(all[1][0], ShouldEqual, 40)
				So(all[1][1], ShouldEqual, 50)
				So(all[1][2], ShouldEqual, 60)
			})

			Convey("Should return empty slice for empty stack", func() {
				cs := core.NewCatchStack()
				
				all := cs.All()
				
				So(all, ShouldHaveLength, 0)
			})
		})

		Convey("Integration", func() {
			Convey("Should handle complex operations", func() {
				cs := core.NewCatchStack()
				
				// Add multiple entries
				cs.Push(0, 10, 100)
				cs.Push(15, 25, 200)
				cs.Push(30, 40, 300)
				
				// Test finding
				catch, found := cs.Find(5)
				So(found, ShouldBeTrue)
				So(catch[2], ShouldEqual, 100)
				
				catch2, found2 := cs.Find(20)
				So(found2, ShouldBeTrue)
				So(catch2[2], ShouldEqual, 200)
				
				catch3, found3 := cs.Find(35)
				So(found3, ShouldBeTrue)
				So(catch3[2], ShouldEqual, 300)
				
				// Pop one
				cs.Pop()
				
				// Should not find in removed range
				_, found4 := cs.Find(35)
				So(found4, ShouldBeFalse)
				
				// Should still find in remaining ranges
				_, found5 := cs.Find(20)
				So(found5, ShouldBeTrue)
			})
		})
	})
}