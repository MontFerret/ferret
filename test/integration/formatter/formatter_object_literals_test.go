package formatter_test

import "testing"

func TestFormatterObjectLiterals(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
LET    foo =     { a: 1, b: 2,  c: 3,    d: 4 }
 RETURN foo
`, `LET foo = { a: 1, b: 2, c: 3, d: 4 }
RETURN foo`),

		Case(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: 4 }
 RETURN foo
`, `LET foo = { a: { e: 5 }, b: 2, c: 3, d: 4 }
RETURN foo`),

		Case(`
LET    foo =      { a: { e: 5 },     
b: 2, c: 3, d: [ 1, 2,3, 4, 5] }
 RETURN foo
`, `LET foo = { a: { e: 5 }, b: 2, c: 3, d: [1, 2, 3, 4, 5] }
RETURN foo`),

		Case(`
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

		Case(`
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
