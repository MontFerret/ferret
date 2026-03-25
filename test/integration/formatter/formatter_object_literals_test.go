package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterObjectLiterals(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
LET    foo =     { a: 1, b: 2,  c: 3,    d: 4 }
 RETURN foo
`, `LET foo = { a: 1, b: 2, c: 3, d: 4 }
RETURN foo`),

		S(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: 4 }
 RETURN foo
`, `LET foo = { a: { e: 5 }, b: 2, c: 3, d: 4 }
RETURN foo`),

		S(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: [ 1, 2,3, 4, 5] }
 RETURN foo
`, `LET foo = { a: { e: 5 }, b: 2, c: 3, d: [1, 2, 3, 4, 5] }
RETURN foo`),

		S(`
LET    foo =      { a: { e: 5 }, 
                b: 2, c: 3, d: [ 1, 2,3, 4, 5], f: {
                  g: []
                } }
 RETURN foo
`, `LET foo = {
    a: { e: 5 },
    b: 2,
    c: 3,
    d: [1, 2, 3, 4, 5],
    f: {
        g: []
    }
}
RETURN foo`),

		S(`
LET    foo =     { 
// comment
a: 1, 
// comment 2
b: 2,  
c: 3,  
   d: 4 }
 RETURN foo
`, `LET foo = {
    // comment
    a: 1,
    // comment 2
    b: 2,
    c: 3,
    d: 4
}
RETURN foo`),
	})
}
