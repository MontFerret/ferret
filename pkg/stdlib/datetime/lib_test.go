package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestRegisterLib(t *testing.T) {
	t.Run("When registering functions", func(t *testing.T) {
		ns := runtime.NewRootNamespace()
		
		err := datetime.RegisterLib(ns)
		
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		
		// Test that some key functions are registered by checking registered functions list
		registeredFunctions := ns.RegisteredFunctions()
		
		expectedFunctions := []string{
			"NOW", "DATE", "DATE_COMPARE", "DATE_DAYOFWEEK", "DATE_YEAR", 
			"DATE_MONTH", "DATE_DAY", "DATE_HOUR", "DATE_MINUTE", "DATE_SECOND",
			"DATE_MILLISECOND", "DATE_DAYOFYEAR", "DATE_LEAPYEAR", "DATE_QUARTER",
			"DATE_DAYS_IN_MONTH", "DATE_FORMAT", "DATE_ADD", "DATE_SUBTRACT", "DATE_DIFF",
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
		
		if foundFunctions < len(expectedFunctions) {
			t.Errorf("not all expected functions were registered, found %d out of %d", foundFunctions, len(expectedFunctions))
		}
	})
}