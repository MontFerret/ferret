package collections_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
)

func TestLib(t *testing.T) {
	Convey("Should register all functions", t, func() {
		ns := runtime.NewLibrary()

		collections.RegisterLib(ns)

		// Check that all expected functions are registered
		expectedFunctions := []string{
			"collections::count_distinct",
			"collections::count",
			"collections::includes",
			"collections::reverse",
			"count_distinct",
			"count",
			"includes",
			"reverse",
		}

		funcs, err := ns.Build()
		So(err, ShouldBeNil)
		registeredFunctions := funcs.List()
		So(len(registeredFunctions), ShouldEqual, len(expectedFunctions))

		for _, funcName := range expectedFunctions {
			found := false
			for _, registered := range registeredFunctions {
				if registered == funcName {
					found = true
					break
				}
			}
			So(found, ShouldBeTrue)
		}
	})
}
