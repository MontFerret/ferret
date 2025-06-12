package base

import "fmt"

import (
	. "github.com/smartystreets/goconvey/convey"
)

func ArePtrsEqual(expected, actual any) bool {
	if expected == nil || actual == nil {
		return false
	}

	p1 := fmt.Sprintf("%v", expected)
	p2 := fmt.Sprintf("%v", actual)

	return p1 == p2
}

func ShouldHaveSameItems(actual any, expected ...any) string {
	wapper := expected[0].([]any)
	expectedArr := wapper[0].([]any)

	for _, item := range expectedArr {
		if err := ShouldContain(actual, item); err != "" {
			return err
		}
	}

	return ""
}

func ShouldBeCompilationError(actual any, _ ...any) string {
	// TODO: Expect a particular error message

	So(actual, ShouldBeError)

	return ""
}
