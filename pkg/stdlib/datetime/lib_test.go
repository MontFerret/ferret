package datetime_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestRegisterLib(t *testing.T) {
	Convey("When registering functions", t, func() {
		ns := runtime.NewLibrary()

		datetime.RegisterLib(ns)

		// Test that some key functions are registered by checking registered functions list
		registry, err := ns.Build()

		So(err, ShouldBeNil)

		registeredFunctions := registry.List()

		expectedFunctions := []string{
			"now", "date", "date_compare", "date_dayofweek", "date_year",
			"date_month", "date_day", "date_hour", "date_minute", "date_second",
			"date_millisecond", "date_dayofyear", "date_leapyear", "date_quarter",
			"date_days_in_month", "date_format", "date_add", "date_subtract", "date_diff",
		}

		foundFunctions := 0
		for _, expectedFunc := range expectedFunctions {
			for _, registeredFunc := range registeredFunctions {
				if registeredFunc == expectedFunc {
					foundFunctions++
					break
				}
			}
		}

		So(foundFunctions, ShouldEqual, len(expectedFunctions))
	})
}
