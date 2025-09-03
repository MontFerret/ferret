package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestLoopTable(t *testing.T) {
	Convey("LoopTable", t, func() {
		Convey("NewLoopTable", func() {
			Convey("Should create a new loop table", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				So(lt, ShouldNotBeNil)
				So(lt.Depth(), ShouldEqual, 0)
				So(lt.Current(), ShouldBeNil)
			})
		})

		Convey(".NewForInLoop", func() {
			Convey("Should create ForIn loop with correct properties", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewForInLoop(core.NormalLoop, true)
				
				So(loop, ShouldNotBeNil)
				So(loop.Kind, ShouldEqual, core.ForInLoop)
				So(loop.Type, ShouldEqual, core.NormalLoop)
				So(loop.Distinct, ShouldBeTrue)
				So(loop.Allocate, ShouldBeTrue)
				So(loop.Dst, ShouldNotEqual, vm.NoopOperand)
			})

			Convey("Should handle different loop types", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				normalLoop := lt.NewForInLoop(core.NormalLoop, false)
				passThroughLoop := lt.NewForInLoop(core.PassThroughLoop, true)
				temporalLoop := lt.NewForInLoop(core.TemporalLoop, false)
				
				So(normalLoop.Type, ShouldEqual, core.NormalLoop)
				So(passThroughLoop.Type, ShouldEqual, core.PassThroughLoop)
				So(temporalLoop.Type, ShouldEqual, core.TemporalLoop)
			})
		})

		Convey(".NewForWhileLoop", func() {
			Convey("Should create ForWhile loop with correct properties", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewForWhileLoop(core.NormalLoop, false)
				
				So(loop, ShouldNotBeNil)
				So(loop.Kind, ShouldEqual, core.ForWhileLoop)
				So(loop.Type, ShouldEqual, core.NormalLoop)
				So(loop.Distinct, ShouldBeFalse)
			})
		})

		Convey(".NewLoop", func() {
			Convey("Should create loop with specified parameters", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewLoop(core.DoWhileLoop, core.NormalLoop, true)
				
				So(loop, ShouldNotBeNil)
				So(loop.Kind, ShouldEqual, core.DoWhileLoop)
				So(loop.Type, ShouldEqual, core.NormalLoop)
				So(loop.Distinct, ShouldBeTrue)
				So(loop.Allocate, ShouldBeTrue)
				So(loop.Dst, ShouldNotEqual, vm.NoopOperand)
			})

			Convey("Should handle temporal loops differently", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewLoop(core.ForInLoop, core.TemporalLoop, false)
				
				So(loop.Type, ShouldEqual, core.TemporalLoop)
				So(loop.Dst, ShouldEqual, vm.NoopOperand) // Temporal loops don't get result registers
			})

			Convey("Should handle nested PassThrough loops", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				// Create parent PassThrough loop
				parentLoop := lt.NewLoop(core.ForInLoop, core.PassThroughLoop, false)
				lt.Push(parentLoop)
				
				// Create child loop
				childLoop := lt.NewLoop(core.ForInLoop, core.NormalLoop, false)
				
				So(childLoop.Allocate, ShouldBeFalse) // Should inherit from PassThrough parent
			})
		})

		Convey(".Push", func() {
			Convey("Should add loop to stack", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				loop := lt.NewForInLoop(core.NormalLoop, false)
				
				lt.Push(loop)
				
				So(lt.Depth(), ShouldEqual, 1)
				So(lt.Current(), ShouldEqual, loop)
			})

			Convey("Should handle multiple loops", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop1 := lt.NewForInLoop(core.NormalLoop, false)
				loop2 := lt.NewForWhileLoop(core.NormalLoop, true)
				
				lt.Push(loop1)
				lt.Push(loop2)
				
				So(lt.Depth(), ShouldEqual, 2)
				So(lt.Current(), ShouldEqual, loop2)
			})
		})

		Convey(".Pop", func() {
			Convey("Should remove and return top loop", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop1 := lt.NewForInLoop(core.NormalLoop, false)
				loop2 := lt.NewForWhileLoop(core.NormalLoop, true)
				
				lt.Push(loop1)
				lt.Push(loop2)
				
				popped := lt.Pop()
				
				So(popped, ShouldEqual, loop2)
				So(lt.Depth(), ShouldEqual, 1)
				So(lt.Current(), ShouldEqual, loop1)
			})

			Convey("Should handle empty stack", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				popped := lt.Pop()
				
				So(popped, ShouldBeNil)
				So(lt.Depth(), ShouldEqual, 0)
			})

			Convey("Should handle popping until empty", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewForInLoop(core.NormalLoop, false)
				lt.Push(loop)
				
				popped1 := lt.Pop()
				popped2 := lt.Pop()
				
				So(popped1, ShouldEqual, loop)
				So(popped2, ShouldBeNil)
				So(lt.Depth(), ShouldEqual, 0)
			})
		})

		Convey(".FindParent", func() {
			Convey("Should find parent loop that allocates", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				// Create a loop that allocates
				allocatingLoop := lt.NewLoop(core.ForInLoop, core.NormalLoop, false)
				allocatingLoop.Allocate = true
				lt.Push(allocatingLoop)
				
				// Create a loop that doesn't allocate
				nonAllocatingLoop := lt.NewLoop(core.ForInLoop, core.PassThroughLoop, false)
				nonAllocatingLoop.Allocate = false
				lt.Push(nonAllocatingLoop)
				
				// Find parent from position 1 (current position)
				parent := lt.FindParent(1)
				
				So(parent, ShouldEqual, allocatingLoop)
			})

			Convey("Should return nil if no allocating parent found", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				// Create loops that don't allocate
				loop1 := lt.NewLoop(core.ForInLoop, core.PassThroughLoop, false)
				loop1.Allocate = false
				loop2 := lt.NewLoop(core.ForInLoop, core.PassThroughLoop, false)
				loop2.Allocate = false
				
				lt.Push(loop1)
				lt.Push(loop2)
				
				parent := lt.FindParent(1)
				
				So(parent, ShouldBeNil)
			})

			Convey("Should handle invalid positions", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				parent := lt.FindParent(0)
				So(parent, ShouldBeNil)
				
				parent2 := lt.FindParent(-1)
				So(parent2, ShouldBeNil)
			})
		})

		Convey(".Current", func() {
			Convey("Should return top of stack", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop := lt.NewForInLoop(core.NormalLoop, false)
				lt.Push(loop)
				
				current := lt.Current()
				
				So(current, ShouldEqual, loop)
			})

			Convey("Should return nil for empty stack", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				current := lt.Current()
				
				So(current, ShouldBeNil)
			})
		})

		Convey(".Depth", func() {
			Convey("Should return correct depth", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				So(lt.Depth(), ShouldEqual, 0)
				
				loop1 := lt.NewForInLoop(core.NormalLoop, false)
				lt.Push(loop1)
				So(lt.Depth(), ShouldEqual, 1)
				
				loop2 := lt.NewForWhileLoop(core.NormalLoop, true)
				lt.Push(loop2)
				So(lt.Depth(), ShouldEqual, 2)
				
				lt.Pop()
				So(lt.Depth(), ShouldEqual, 1)
			})
		})

		Convey(".DebugView", func() {
			Convey("Should return debug information", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				loop1 := lt.NewForInLoop(core.NormalLoop, false)
				loop2 := lt.NewForWhileLoop(core.PassThroughLoop, true)
				
				lt.Push(loop1)
				lt.Push(loop2)
				
				debug := lt.DebugView()
				
				So(debug, ShouldNotBeEmpty)
				So(debug, ShouldContainSubstring, "Loop[0]")
				So(debug, ShouldContainSubstring, "Loop[1]")
			})

			Convey("Should handle empty stack", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				debug := lt.DebugView()
				
				So(debug, ShouldEqual, "")
			})
		})

		Convey("Integration", func() {
			Convey("Should handle complex loop nesting", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				// Create nested loop structure
				outerLoop := lt.NewForInLoop(core.NormalLoop, false)
				lt.Push(outerLoop)
				
				innerLoop1 := lt.NewForWhileLoop(core.PassThroughLoop, true)
				lt.Push(innerLoop1)
				
				innerLoop2 := lt.NewLoop(core.DoWhileLoop, core.TemporalLoop, false)
				lt.Push(innerLoop2)
				
				// Test operations
				So(lt.Depth(), ShouldEqual, 3)
				So(lt.Current(), ShouldEqual, innerLoop2)
				
				// Find parent
				parent := lt.FindParent(2) // Should find first allocating loop
				So(parent, ShouldNotBeNil)
				
				// Pop and verify
				popped := lt.Pop()
				So(popped, ShouldEqual, innerLoop2)
				So(lt.Current(), ShouldEqual, innerLoop1)
				
				// Continue popping
				lt.Pop()
				lt.Pop()
				
				So(lt.Depth(), ShouldEqual, 0)
				So(lt.Current(), ShouldBeNil)
			})

			Convey("Should create loops with proper inheritance", func() {
				ra := core.NewRegisterAllocator()
				lt := core.NewLoopTable(ra)
				
				// Create PassThrough parent
				parentLoop := lt.NewLoop(core.ForInLoop, core.PassThroughLoop, false)
				lt.Push(parentLoop)
				
				// Create child - should not allocate due to PassThrough parent
				childLoop := lt.NewLoop(core.ForWhileLoop, core.NormalLoop, true)
				
				So(childLoop.Allocate, ShouldBeFalse)
				So(childLoop.Dst, ShouldEqual, parentLoop.Dst)
			})
		})
	})
}