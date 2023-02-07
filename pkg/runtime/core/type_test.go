package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Value struct {
		type_ core.Type
	}

	TypeA struct{}

	TypeB struct{}

	TypeC struct{}
)

func (t Value) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (t Value) Type() core.Type {
	return t.type_
}

func (t Value) String() string {
	return ""
}

func (t Value) Compare(other core.Value) int64 {
	return 0
}

func (t Value) Unwrap() interface{} {
	return nil
}

func (t Value) Hash() uint64 {
	return 0
}

func (t Value) Copy() core.Value {
	return t
}

func (t TypeA) ID() int64 {
	return 1
}

func (t TypeA) String() string {
	return "type_a"
}

func (t TypeA) Equals(other core.Type) bool {
	return t.ID() == other.ID()
}

func (t TypeB) ID() int64 {
	return 2
}

func (t TypeB) String() string {
	return "type_b"
}

func (t TypeB) Equals(other core.Type) bool {
	return t.ID() == other.ID()
}

func (t TypeC) ID() int64 {
	return 3
}

func (t TypeC) String() string {
	return "type_c"
}

func (t TypeC) Equals(other core.Type) bool {
	return t.ID() == other.ID()
}

func TestType(t *testing.T) {
	typeA := TypeA{}
	typeB := TypeB{}

	Convey("IsTypeOf", t, func() {
		Convey("Should return 'false' when types are different", func() {
			vA := Value{typeA}

			So(core.IsTypeOf(vA, typeB), ShouldBeFalse)
		})
	})
}

func TestValidateType(t *testing.T) {
	typeA := TypeA{}
	typeB := TypeB{}

	Convey("Should validate types", t, func() {
		vA := Value{typeA}
		vB := Value{typeB}

		So(core.ValidateType(vA, typeA), ShouldBeNil)
		So(core.ValidateType(vB, typeA), ShouldNotBeNil)
	})
}
