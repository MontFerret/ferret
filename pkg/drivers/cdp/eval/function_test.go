package eval

import (
	"testing"

	"github.com/mafredri/cdp/protocol/runtime"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestFunction(t *testing.T) {
	Convey("Function", t, func() {
		Convey(".AsAsync", func() {
			Convey("Should set async=true", func() {
				f := F("return 'foo'").AsAsync()
				args := f.call(EmptyExecutionContextID)

				So(*args.AwaitPromise, ShouldBeTrue)
			})
		})

		Convey(".AsSync", func() {
			Convey("Should set async=false", func() {
				f := F("return 'foo'").AsAsync()
				args := f.call(EmptyExecutionContextID)

				So(*args.AwaitPromise, ShouldBeTrue)

				args = f.AsSync().call(EmptyExecutionContextID)

				So(*args.AwaitPromise, ShouldBeFalse)
			})
		})

		Convey(".CallOn", func() {
			Convey("It should use a given ownerID over ContextID", func() {
				ownerID := runtime.RemoteObjectID("foo")
				contextID := runtime.ExecutionContextID(42)

				f := F("return 'foo'").CallOn(ownerID)
				call := f.call(contextID)

				So(call.ExecutionContextID, ShouldBeNil)
				So(call.ObjectID, ShouldNotBeNil)
				So(*call.ObjectID, ShouldEqual, ownerID)
			})

			Convey("It should use a given ContextID when ownerID is empty or nil", func() {
				ownerID := runtime.RemoteObjectID("")
				contextID := runtime.ExecutionContextID(42)

				f := F("return 'foo'").CallOn(ownerID)
				call := f.call(contextID)

				So(call.ExecutionContextID, ShouldNotBeNil)
				So(call.ObjectID, ShouldBeNil)
				So(*call.ExecutionContextID, ShouldEqual, contextID)
			})
		})

		Convey(".WithArgRef", func() {
			Convey("Should add argument with a given RemoteObjectID", func() {
				f := F("return 'foo'")
				id1 := runtime.RemoteObjectID("foo")
				id2 := runtime.RemoteObjectID("bar")
				id3 := runtime.RemoteObjectID("baz")

				f.WithArgRef(id1).WithArgRef(id2).WithArgRef(id3)

				So(f.Length(), ShouldEqual, 3)

				arg1 := f.args[0]
				arg2 := f.args[1]
				arg3 := f.args[2]

				So(*arg1.ObjectID, ShouldEqual, id1)
				So(arg1.Value, ShouldBeNil)
				So(arg1.UnserializableValue, ShouldBeNil)

				So(*arg2.ObjectID, ShouldEqual, id2)
				So(arg2.Value, ShouldBeNil)
				So(arg2.UnserializableValue, ShouldBeNil)

				So(*arg3.ObjectID, ShouldEqual, id3)
				So(arg3.Value, ShouldBeNil)
				So(arg3.UnserializableValue, ShouldBeNil)
			})
		})

		Convey(".WithArgValue", func() {
			Convey("Should add argument with a given Value", func() {
				f := F("return 'foo'")
				val1 := values.NewString("foo")
				val2 := values.NewInt(1)
				val3 := values.NewBoolean(true)

				f.WithArgValue(val1).WithArgValue(val2).WithArgValue(val3)

				So(f.Length(), ShouldEqual, 3)

				arg1 := f.args[0]
				arg2 := f.args[1]
				arg3 := f.args[2]

				So(arg1.ObjectID, ShouldBeNil)
				So(arg1.Value, ShouldResemble, values.MustMarshal(val1))
				So(arg1.UnserializableValue, ShouldBeNil)

				So(arg2.ObjectID, ShouldBeNil)
				So(arg2.Value, ShouldResemble, values.MustMarshal(val2))
				So(arg2.UnserializableValue, ShouldBeNil)

				So(arg3.ObjectID, ShouldBeNil)
				So(arg3.Value, ShouldResemble, values.MustMarshal(val3))
				So(arg3.UnserializableValue, ShouldBeNil)
			})
		})

		Convey(".WithArg", func() {
			Convey("Should add argument with a given any type", func() {
				f := F("return 'foo'")
				val1 := "foo"
				val2 := 1
				val3 := true

				f.WithArg(val1).WithArg(val2).WithArg(val3)

				So(f.Length(), ShouldEqual, 3)

				arg1 := f.args[0]
				arg2 := f.args[1]
				arg3 := f.args[2]

				So(arg1.ObjectID, ShouldBeNil)
				So(arg1.Value, ShouldResemble, values.MustMarshalAny(val1))
				So(arg1.UnserializableValue, ShouldBeNil)

				So(arg2.ObjectID, ShouldBeNil)
				So(arg2.Value, ShouldResemble, values.MustMarshalAny(val2))
				So(arg2.UnserializableValue, ShouldBeNil)

				So(arg3.ObjectID, ShouldBeNil)
				So(arg3.Value, ShouldResemble, values.MustMarshalAny(val3))
				So(arg3.UnserializableValue, ShouldBeNil)
			})
		})

		Convey(".WithArgSelector", func() {
			Convey("Should add argument with a given QuerySelector", func() {
				f := F("return 'foo'")
				val1 := drivers.NewCSSSelector(".foo-bar")
				val2 := drivers.NewCSSSelector("#submit")
				val3 := drivers.NewXPathSelector("//*[@id='q']")

				f.WithArgSelector(val1).WithArgSelector(val2).WithArgSelector(val3)

				So(f.Length(), ShouldEqual, 3)

				arg1 := f.args[0]
				arg2 := f.args[1]
				arg3 := f.args[2]

				So(arg1.ObjectID, ShouldBeNil)
				So(arg1.Value, ShouldResemble, values.MustMarshalAny(val1.String()))
				So(arg1.UnserializableValue, ShouldBeNil)

				So(arg2.ObjectID, ShouldBeNil)
				So(arg2.Value, ShouldResemble, values.MustMarshalAny(val2.String()))
				So(arg2.UnserializableValue, ShouldBeNil)

				So(arg3.ObjectID, ShouldBeNil)
				So(arg3.Value, ShouldResemble, values.MustMarshalAny(val3.String()))
				So(arg3.UnserializableValue, ShouldBeNil)
			})
		})

		Convey(".String", func() {
			Convey("It should return a function expression", func() {
				exp := "return 'foo'"
				f := F(exp)

				So(f.String(), ShouldEqual, exp)
			})
		})

		Convey(".returnNothing", func() {
			Convey("It should set return by value to false", func() {
				f := F("return 'foo'").returnNothing()
				call := f.call(EmptyExecutionContextID)

				So(*call.ReturnByValue, ShouldBeFalse)
			})
		})

		Convey(".returnValue", func() {
			Convey("It should set return by value to true", func() {
				f := F("return 'foo'").returnValue()
				call := f.call(EmptyExecutionContextID)

				So(*call.ReturnByValue, ShouldBeTrue)
			})
		})

		Convey(".returnRef", func() {
			Convey("It should set return by value to false", func() {
				f := F("return 'foo'").returnValue()
				call := f.call(EmptyExecutionContextID)

				So(*call.ReturnByValue, ShouldBeTrue)

				f.returnRef()

				call = f.call(EmptyExecutionContextID)

				So(*call.ReturnByValue, ShouldBeFalse)
			})
		})
	})

	//Convey("prepExp", t, func() {
	//	Convey("When a plain expression is passed", func() {
	//		exp := "return true"
	//		So(wrapExp(exp, 0), ShouldEqual, "() => {\n"+exp+"\n}")
	//	})
	//
	//	Convey("When a plain expression is passed with args > 0", func() {
	//		exp := "return true"
	//		So(wrapExp(exp, 3), ShouldEqual, "(arg1,arg2,arg3) => {\n"+exp+"\n}")
	//	})
	//})
}
