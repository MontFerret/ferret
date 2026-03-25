package formatter_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterArrayLiterals(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
LET    foo =      [ 1 , 2, 3,   4  ]
 RETURN foo
`, `LET foo = [1, 2, 3, 4]
RETURN foo`),

		S(`
LET    foo =      [ 
1 , 
2, 
3,   
4  ]
 RETURN foo
`, `LET foo = [1, 2, 3, 4]
RETURN foo`),

		S(`
// comment
LET    foo =      [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40]
 RETURN foo
`, `// comment
LET foo = [
    1,
    2,
    3,
    4,
    5,
    6,
    7,
    8,
    9,
    10,
    11,
    12,
    13,
    14,
    15,
    16,
    17,
    18,
    19,
    20,
    21,
    22,
    23,
    24,
    25,
    26,
    27,
    28,
    29,
    30,
    31,
    32,
    33,
    34,
    35,
    36,
    37,
    38,
    39,
    40
]
RETURN foo`),

		S(`
LET    foo =      [ 
{ a: 1, b: 2}, 
{a: 3, b: 4}
]
 RETURN foo
`, `LET foo = [
    { a: 1, b: 2 },
    { a: 3, b: 4 }
]
RETURN foo`),

		S(`
LET    foo =      [ 
// comment
{ a: 1, b: 2},
// comment 2
{a: 3, b: 4}
]
 RETURN foo
`, `LET foo = [
    // comment
    { a: 1, b: 2 },
    // comment 2
    { a: 3, b: 4 }
]
RETURN foo`),

		S(`
LET    foo =      [ 
// comment
{ a: 1, b: 2},
// comment 2
{a: 3, b: 4}, [ 
// comment 3
{ a: 1, b: 2},
// comment 4
{a: 3, b: 4}
]
]
 RETURN foo
`, `LET foo = [
    // comment
    { a: 1, b: 2 },
    // comment 2
    { a: 3, b: 4 },
    [
        // comment 3
        { a: 1, b: 2 },
        // comment 4
        { a: 3, b: 4 }
    ]
]
RETURN foo`),
	})
}
