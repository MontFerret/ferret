package assert

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	ShouldBeNil              = NewUnaryAssertion(Nil)
	ShouldBeError            = NewUnaryAssertion(Error)
	ShouldNotBeNil           = NewUnaryAssertion(NotNil)
	ShouldBeDiagnosticError  = NewBinaryAssertion(DiagnosticError)
	ShouldBeDiagnosticErrors = NewBinaryAssertion(DiagnosticErrors)
	ShouldEqual              = NewBinaryAssertion(Equal)
	ShouldHaveSameItems      = NewBinaryAssertion(HaveSameItems)
	ShouldEqualJSON          = NewBinaryAssertion(EqualJSON)
	ShouldNotReach           = func(t *testing.T, _ any, _ ...any) {
		t.Helper()

		t.Fatal("this assertion should not be reached")
	}
)

func IsSameAssertion(a, b Assertion) bool {
	if a == nil || b == nil {
		return false
	}

	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Kind() == reflect.Func && v2.Kind() == reflect.Func {
		return v1.Pointer() == v2.Pointer()
	}

	p1 := fmt.Sprintf("%T", a)
	p2 := fmt.Sprintf("%T", b)

	return p1 == p2
}
