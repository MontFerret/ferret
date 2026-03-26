package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterLiterals(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
LET    foo =      10
 RETURN foo
`, `LET foo = 10
RETURN foo`),

		S(`
LET    foo =      
10
 RETURN foo
`, `LET foo = 10
RETURN foo`),

		S(`
LET    foo        =      "bar"
 RETURN foo
`, `LET foo = "bar"
RETURN foo`),

		S(`
LET    foo        =    
     "bar"
 RETURN foo
`, `LET foo = "bar"
RETURN foo`),

		S(`
LET    foo        =      [     ]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		S(`
LET    foo        =      [     
]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		S(`
LET    foo        =     
       [     
]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		S(`
LET    foo        =   {        }
 RETURN foo
`, `LET foo = {}
RETURN foo`),

		S(`
LET    foo        =   {        

}
 RETURN foo
`, `LET foo = {}
RETURN foo`),

		S(`
LET    foo        =   

      {        

}
 RETURN foo
`, `LET foo = {}
RETURN foo`),
	})
}
