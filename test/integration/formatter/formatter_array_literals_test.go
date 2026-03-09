package formatter_test

import "testing"

func TestFormatterArrayLiterals(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
LET    foo =      [ 1 , 2, 3,   4  ]
 RETURN foo
`, `LET foo = [1, 2, 3, 4]
RETURN foo`),

		Case(`
LET    foo =      [ 
1 , 
2, 
3,   
4  ]
 RETURN foo
`, `LET foo = [1, 2, 3, 4]
RETURN foo`),

		Case(`
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

		Case(`
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

		Case(`
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

		Case(`
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
