package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterUDFs(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
func normalizePrice( value ){
return value
}
return normalizePrice(1)
`, `func normalizePrice(value) {
    return value
}
return normalizePrice(1)`),
		S(`
func unique( values ){
return distinct values
}
return unique([1, 1])
`, `func unique(values) {
    return distinct values
}
return unique([1, 1])`),
		S(`
func read( DISTINCT ){
return ( DISTINCT.values )
}
return ( DISTINCT.values )
`, `func read(DISTINCT) {
    return (DISTINCT.values)
}
return (DISTINCT.values)`),
	})
}
