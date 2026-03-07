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
			"COUNT_DISTINCT",
			"COUNT",
			"INCLUDES",
			"REVERSE",
		}

		funcs, err := ns.Build()
		So(err, ShouldBeNil)
		registeredFunctions := funcs.List()

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
