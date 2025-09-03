package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestConstantPool(t *testing.T) {
	Convey("ConstantPool", t, func() {
		Convey("NewConstantPool", func() {
			Convey("Should create an empty constant pool", func() {
				cp := core.NewConstantPool()
				So(cp, ShouldNotBeNil)
				So(cp.All(), ShouldHaveLength, 0)
			})
		})

		Convey(".Add", func() {
			Convey("Should add a constant and return an operand", func() {
				cp := core.NewConstantPool()
				val := runtime.NewString("test")
				
				op := cp.Add(val)
				
				So(op.IsConstant(), ShouldBeTrue)
				So(cp.All(), ShouldHaveLength, 1)
				So(cp.All()[0], ShouldEqual, val)
			})

			Convey("Should deduplicate scalar values", func() {
				cp := core.NewConstantPool()
				val1 := runtime.NewString("test")
				val2 := runtime.NewString("test")
				
				op1 := cp.Add(val1)
				op2 := cp.Add(val2)
				
				So(op1, ShouldEqual, op2)
				So(cp.All(), ShouldHaveLength, 1)
			})

			Convey("Should deduplicate None values", func() {
				cp := core.NewConstantPool()
				
				op1 := cp.Add(runtime.None)
				op2 := cp.Add(runtime.None)
				
				So(op1, ShouldEqual, op2)
				So(cp.All(), ShouldHaveLength, 1)
			})

			Convey("Should not deduplicate non-scalar values", func() {
				cp := core.NewConstantPool()
				val1 := runtime.NewArray(10)
				val2 := runtime.NewArray(10)
				
				op1 := cp.Add(val1)
				op2 := cp.Add(val2)
				
				So(op1, ShouldNotEqual, op2)
				So(cp.All(), ShouldHaveLength, 2)
			})

			Convey("Should handle different types of scalar values", func() {
				cp := core.NewConstantPool()
				
				strOp := cp.Add(runtime.NewString("hello"))
				intOp := cp.Add(runtime.NewInt(42))
				floatOp := cp.Add(runtime.NewFloat(3.14))
				boolOp := cp.Add(runtime.NewBoolean(true))
				
				So(strOp.IsConstant(), ShouldBeTrue)
				So(intOp.IsConstant(), ShouldBeTrue)
				So(floatOp.IsConstant(), ShouldBeTrue)
				So(boolOp.IsConstant(), ShouldBeTrue)
				
				So(cp.All(), ShouldHaveLength, 4)
			})
		})

		Convey(".Get", func() {
			Convey("Should retrieve a constant by operand", func() {
				cp := core.NewConstantPool()
				val := runtime.NewString("test")
				
				op := cp.Add(val)
				retrieved := cp.Get(op)
				
				So(retrieved, ShouldEqual, val)
			})

			Convey("Should panic with invalid operand type", func() {
				cp := core.NewConstantPool()
				invalidOp := vm.NewRegister(0)
				
				So(func() { cp.Get(invalidOp) }, ShouldPanic)
			})

			Convey("Should panic with invalid constant index", func() {
				cp := core.NewConstantPool()
				invalidOp := vm.NewConstant(999) // out of range
				
				So(func() { cp.Get(invalidOp) }, ShouldPanic)
			})

			Convey("Should panic with negative constant index", func() {
				cp := core.NewConstantPool()
				invalidOp := vm.NewConstant(-1)
				
				So(func() { cp.Get(invalidOp) }, ShouldPanic)
			})
		})

		Convey(".All", func() {
			Convey("Should return all constants in order", func() {
				cp := core.NewConstantPool()
				val1 := runtime.NewString("first")
				val2 := runtime.NewInt(42)
				val3 := runtime.NewFloat(3.14)
				
				cp.Add(val1)
				cp.Add(val2)
				cp.Add(val3)
				
				all := cp.All()
				So(all, ShouldHaveLength, 3)
				So(all[0], ShouldEqual, val1)
				So(all[1], ShouldEqual, val2)
				So(all[2], ShouldEqual, val3)
			})

			Convey("Should return empty slice for empty pool", func() {
				cp := core.NewConstantPool()
				
				all := cp.All()
				So(all, ShouldHaveLength, 0)
			})
		})

		Convey("Integration", func() {
			Convey("Should maintain consistency across operations", func() {
				cp := core.NewConstantPool()
				
				// Add multiple constants
				str1 := runtime.NewString("hello")
				str2 := runtime.NewString("hello") // duplicate
				int1 := runtime.NewInt(42)
				float1 := runtime.NewFloat(3.14)
				
				op1 := cp.Add(str1)
				op2 := cp.Add(str2) // should be same as op1
				op3 := cp.Add(int1)
				op4 := cp.Add(float1)
				
				// Check deduplication worked
				So(op1, ShouldEqual, op2)
				
				// Verify we can retrieve all values
				So(cp.Get(op1), ShouldEqual, str1)
				So(cp.Get(op3), ShouldEqual, int1)
				So(cp.Get(op4), ShouldEqual, float1)
				
				// Check total count (3 because str1 and str2 are deduplicated)
				So(cp.All(), ShouldHaveLength, 3)
			})
		})
	})
}