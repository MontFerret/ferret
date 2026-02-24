package runtime_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type typedOnly struct{}

func (typedOnly) Type() runtime.Type {
	return runtime.TypeMap
}

func (typedOnly) String() string {
	return "typedOnly"
}

func (typedOnly) Hash() uint64 {
	return 0
}

func (typedOnly) Copy() runtime.Value {
	return typedOnly{}
}

func TestValidateType(t *testing.T) {
	type testCase struct {
		Name     string
		Input    runtime.Value
		Expected runtime.Type
		Failure  bool
	}

	tests := []testCase{
		{
			Name:     "None",
			Input:    runtime.None,
			Expected: runtime.TypeNone,
			Failure:  false,
		},
		{
			Name:     "True",
			Input:    runtime.True,
			Expected: runtime.TypeBoolean,
			Failure:  false,
		},
		{
			Name:     "Int",
			Input:    runtime.NewInt(42),
			Expected: runtime.TypeInt,
			Failure:  false,
		},
		{
			Name:     "Float",
			Input:    runtime.NewFloat(42),
			Expected: runtime.TypeFloat,
			Failure:  false,
		},
		{
			Name:     "String",
			Input:    runtime.NewString("hello"),
			Expected: runtime.TypeString,
			Failure:  false,
		},
		{
			Name:     "DateTime",
			Input:    runtime.NewCurrentDateTime(),
			Expected: runtime.TypeDateTime,
			Failure:  false,
		},
		{
			Name:     "Binary",
			Input:    runtime.NewBinary([]byte{1, 2, 3}),
			Expected: runtime.TypeBinary,
			Failure:  false,
		},
		// TODO: Where is Regexp?
		//{
		//	Input:    runtime.NewRegexp(".*"),
		//	Expected: runtime.TypeRegexp,
		//	Failure:  false,
		//},
		{
			Name:     "Array",
			Input:    runtime.NewArrayWith(runtime.NewInt(1)),
			Expected: runtime.TypeArray,
			Failure:  false,
		},
		{
			Name:     "Array as List",
			Input:    runtime.NewArrayWith(runtime.NewInt(1)),
			Expected: runtime.TypeList,
			Failure:  false,
		},
		{
			Name:     "Object as Map",
			Input:    runtime.NewObject(),
			Expected: runtime.TypeMap,
			Failure:  false,
		},
		{
			Name:     "Array as Iterable",
			Input:    runtime.NewArrayWith(runtime.NewInt(1)),
			Expected: runtime.TypeIterable,
			Failure:  false,
		},
		{
			Name:     "Typed-only Map",
			Input:    typedOnly{},
			Expected: runtime.TypeMap,
			Failure:  false,
		},
		{
			Name:     "Array as Map",
			Input:    runtime.NewArrayWith(runtime.NewInt(1)),
			Expected: runtime.TypeMap,
			Failure:  true,
		},
	}

	Convey("ValidateType should correctly validate types", t, func() {
		for _, tCase := range tests {
			Convey(tCase.Name, func() {
				err := runtime.ValidateType(tCase.Input, tCase.Expected)

				if tCase.Failure {
					SoMsg(fmt.Sprintf("expected failure for input %v, but got no error", tCase.Input), err, ShouldNotBeNil)
				} else {
					SoMsg(fmt.Sprintf("expected success for input %v, but got an error", tCase.Input), err, ShouldBeNil)
				}
			})
		}
	})
}

func TestTypeOfTypedOverride(t *testing.T) {
	Convey("TypeOf should respect Typed overrides", t, func() {
		So(runtime.SameType(runtime.TypeOf(typedOnly{}), runtime.TypeMap), ShouldBeTrue)
	})

	Convey("IsType should use TypeOf before interface checks", t, func() {
		So(runtime.IsType(typedOnly{}, runtime.TypeMap), ShouldBeTrue)
		So(runtime.IsType(typedOnly{}, runtime.TypeList), ShouldBeFalse)
	})
}
