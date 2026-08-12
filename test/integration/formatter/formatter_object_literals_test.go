package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterObjectLiterals(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
LET    foo =     { a: 1, b: 2,  c: 3,    d: 4 }
 return foo
`, `let foo = { a: 1, b: 2, c: 3, d: 4 }
return foo`),

		S(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: 4 }
 return foo
`, `let foo = { a: { e: 5 }, b: 2, c: 3, d: 4 }
return foo`),

		S(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: [ 1, 2,3, 4, 5] }
 return foo
`, `let foo = { a: { e: 5 }, b: 2, c: 3, d: [1, 2, 3, 4, 5] }
return foo`),

		S(`
LET    foo =      { a: { e: 5 }, 
                b: 2, c: 3, d: [ 1, 2,3, 4, 5], f: {
                  g: []
                } }
 return foo
`, `let foo = {
    a: { e: 5 },
    b: 2,
    c: 3,
    d: [1, 2, 3, 4, 5],
    f: {
        g: []
    }
}
return foo`),

		S(`
LET    foo =     { 
// comment
a: 1, 
// comment 2
b: 2,  
c: 3,  
   d: 4 }
 return foo
`, `let foo = {
    // comment
    a: 1,
    // comment 2
    b: 2,
    c: 3,
    d: 4
}
return foo`),
	})
}
