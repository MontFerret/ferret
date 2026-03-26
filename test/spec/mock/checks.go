package mock

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func ObservableReturnOneAndReads(obs *Observable, expectedReads int32) assert.Unary {
	return func(actual any) error {
		var ok bool
		switch v := actual.(type) {
		case float64:
			ok = v == 1
		case int:
			ok = v == 1
		case int64:
			ok = v == 1
		}

		if !ok {
			return fmt.Errorf("expected return value 1, got %v", actual)
		}

		if reads := obs.ReadCount(); reads != expectedReads {
			return fmt.Errorf("expected %d reads, got %d", expectedReads, reads)
		}

		return nil
	}
}
