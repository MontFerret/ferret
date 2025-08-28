package collections_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

func TestLib(t *testing.T) {
	Convey("Should register all functions", t, func() {
		ns := runtime.NewRootNamespace()
		
		err := collections.RegisterLib(ns)
		
		So(err, ShouldBeNil)
		
		// Check that all expected functions are registered
		expectedFunctions := []string{
			"COUNT_DISTINCT",
			"COUNT", 
			"INCLUDES",
			"REVERSE",
		}
		
		registeredFunctions := ns.RegisteredFunctions()
		
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

	Convey("Should not return error on valid namespace", t, func() {
		ns := runtime.NewRootNamespace()
		
		err := collections.RegisterLib(ns)
		
		So(err, ShouldBeNil)
	})
}