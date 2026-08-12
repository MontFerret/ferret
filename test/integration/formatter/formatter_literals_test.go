package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterLiterals(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
LET    foo =      10
 return foo
`, `let foo = 10
return foo`),

		S(`
LET    foo =      
10
 return foo
`, `let foo = 10
return foo`),

		S(`
LET    foo        =      "bar"
 return foo
`, `let foo = "bar"
return foo`),

		S(`
LET    foo        =    
     "bar"
 return foo
`, `let foo = "bar"
return foo`),

		S(`
LET    foo        =      [     ]
 return foo
`, `let foo = []
return foo`),

		S(`
LET    foo        =      [     
]
 return foo
`, `let foo = []
return foo`),

		S(`
LET    foo        =     
       [     
]
 return foo
`, `let foo = []
return foo`),

		S(`
LET    foo        =   {        }
 return foo
`, `let foo = {}
return foo`),

		S(`
LET    foo        =   {        

}
 return foo
`, `let foo = {}
return foo`),

		S(`
LET    foo        =   

      {        

}
 return foo
`, `let foo = {}
return foo`),
	})
}
