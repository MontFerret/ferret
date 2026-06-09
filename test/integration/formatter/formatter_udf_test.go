package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterUDFs(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
FUNC normalizePrice( value )(
RETURN value
)
RETURN normalizePrice(1)
`, `FUNC normalizePrice(value) (
    RETURN value
)
RETURN normalizePrice(1)`),
		S(`
FUNC unique( values )(
RETURN DISTINCT values
)
RETURN unique([1, 1])
`, `FUNC unique(values) (
    RETURN DISTINCT values
)
RETURN unique([1, 1])`),
	})
}
