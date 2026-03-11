package formatter_test

import "testing"

func TestFormatterLiterals(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
LET    foo =      10
 RETURN foo
`, `LET foo = 10
RETURN foo`),

		Case(`
LET    foo =      
10
 RETURN foo
`, `LET foo = 10
RETURN foo`),

		Case(`
LET    foo        =      "bar"
 RETURN foo
`, `LET foo = "bar"
RETURN foo`),

		Case(`
LET    foo        =    
     "bar"
 RETURN foo
`, `LET foo = "bar"
RETURN foo`),

		Case(`
LET    foo        =      [     ]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		Case(`
LET    foo        =      [     
]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		Case(`
LET    foo        =     
       [     
]
 RETURN foo
`, `LET foo = []
RETURN foo`),

		Case(`
LET    foo        =   {        }
 RETURN foo
`, `LET foo = {}
RETURN foo`),

		Case(`
LET    foo        =   {        

}
 RETURN foo
`, `LET foo = {}
RETURN foo`),

		Case(`
LET    foo        =   

      {        

}
 RETURN foo
`, `LET foo = {}
RETURN foo`),
	})
}
