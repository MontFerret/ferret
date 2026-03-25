package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestArrayOperators(t *testing.T) {
	RunSpecs(t, []Spec{
		Array(`
					LET friends = [
						{ name: "John", age: 30 },
						{ name: "Mary", age: 25 },
						{ name: "Bob", age: 50 },
						{ name: "Alice", age: 28 },
						{ name: "Tom", age: 35 },
						{ name: "Jane", age: 32 }
					]
					LET users = [
						{ 
							name: "John", 
							age: 30,
							friends: [	
								{ name: "Alice", age: 28 },
							]
						},
						{
							name: "Mary", 
							age: 25,
							friends: [
								{ name: "Tom", age: 35 },
								{ name: "Jane", age: 32 }
							]
						},
						{ 
							name: "Bob", 
							age: 50,
							friends: []
						}
					]

					FOR user IN users
						RETURN user.friends[*].name
				`,
			[]any{[]any{"Alice"}, []any{"Tom", "Jane"}, []any{}}),
		Array(`
					LET users = [
						{ name: "John", age: 30 },
						{ name: "Mary", age: 25 },
						{ name: "Bob", age: 50 }
					]

					RETURN users[*].name
				`,
			[]any{"John", "Mary", "Bob"}),
		Array(`
					LET users = [
						{ name: [ { num: [1, 2] }, { num: [4] } ] },
						{ name: [ { num: [5] } ] }
					]

					RETURN users[*].name[*].num
				`,
			[]any{
				[]any{
					[]any{1, 2},
					[]any{4},
				},
				[]any{
					[]any{5},
				},
			}),
		Array(`
					LET users = [
						{ name: "Ann", age: 20 },
						{ name: "Bob", age: 45 },
						{ name: "Cat", age: 50 }
					]

					RETURN users[*][* FILTER .age > 40].age
				`,
			[]any{45, 50}),
		Array(`
					LET users = [
						{ name: "Ann", age: 20 },
						{ name: "Bob", age: 35 },
						{ name: "Cat", age: 45 },
						{ name: "Dan", age: 55 }
					]

					RETURN users[* FILTER .age > 20][* FILTER .age < 50].name
				`,
			[]any{"Bob", "Cat"}),
		Array(`
					LET users = [
						{ name: "Ann", age: 20 },
						{ name: "Bob", age: 35 },
						{ name: "Cat", age: 45 },
						{ name: "Dan", age: 55 }
					]

					RETURN users[* FILTER .age > 30][*].name
				`,
			[]any{"Bob", "Cat", "Dan"}),
		Array(`
					LET arr = [[1, 2], 3, [4, 5], 6]

					RETURN arr[**]
				`,
			[]any{1, 2, 3, 4, 5, 6}),
		Array(`
					LET arr = [[[1], [2]], [[3]]]

					RETURN arr[***]
				`,
			[]any{1, 2, 3}),
		Array(`
					LET users = [
						[{ name: "Ann" }, { name: "Ben" }],
						[{ name: "Cat" }]
					]

					RETURN users[**].name
				`,
			[]any{"Ann", "Ben", "Cat"}),
		Array(`
LET users = [
						{ 
							name: "John", 
							age: 30,
							friends: [	
								{ name: "Alice", age: 28 },
							]
						},
						{
							name: "Mary", 
							age: 25,
							friends: [
								{ name: "Tom", age: 35 },
								{ name: "Jane", age: 32 }
							]
						},
						{ 
							name: "Bob", 
							age: 50,
							friends: []
						}
					]

					RETURN (
					  FOR u IN users RETURN u.friends[*].name
					)[**]
				`,
			[]any{"Alice", "Tom", "Jane"}),
		Array(`
LET arr = [ [ 1, 2 ], 3, [ 4, 5 ], 6 ]
RETURN arr[** FILTER . % 2 == 0]`, []any{2, 4, 6}),
		Array(`
LET values = [1, 2, 3, 4]
RETURN values[* LIMIT 2]`, []any{1, 2}),
		Array(`
LET values = [1, 2, 3, 4]
RETURN values[* LIMIT 1, 2]`, []any{2, 3}),
		Array(`
LET values = [1, 2, 3]
RETURN values[* RETURN . * 2]`, []any{2, 4, 6}),
		Array(`
LET values = [1, 2, 3]
RETURN values[* RETURN .]`, []any{1, 2, 3}),
		Array(`
LET users = [
	{ name: "Ann", friends: [ { name: "Bob", age: 20 }, { name: "Cat", age: 17 } ] },
	{ name: "Dan", friends: [ { name: "Eve", age: 19 } ] },
	{ name: "Liz", friends: [ { name: "Max", age: 16 } ] }
]
RETURN users[* RETURN {
	name: .name,
	friends: .friends[* FILTER .age > 18][* RETURN .name]
}]`, []any{
			map[string]any{"name": "Ann", "friends": []any{"Bob"}},
			map[string]any{"name": "Dan", "friends": []any{"Eve"}},
			map[string]any{"name": "Liz", "friends": []any{}},
		}),
		Array(`
LET users = [
	{ name: "Ann", friends: [ { name: "Bob", age: 20 }, { name: "Cat", age: 17 } ] },
	{ name: "Dan", friends: [ { name: "Eve", age: 19 } ] },
	{ name: "Liz", friends: [ { name: "Max", age: 16 } ] }
]
RETURN users[* RETURN {
	name: .name,
	friends: .friends[* RETURN .name],
	again: .name
}]`, []any{
			map[string]any{"name": "Ann", "friends": []any{"Bob", "Cat"}, "again": "Ann"},
			map[string]any{"name": "Dan", "friends": []any{"Eve"}, "again": "Dan"},
			map[string]any{"name": "Liz", "friends": []any{"Max"}, "again": "Liz"},
		}),
		Array(`
LET users = [
	{ name: "Ann", friends: [ { name: "Bob", age: 20 }, { name: "Cat", age: 17 } ] },
	{ name: "Dan", friends: [ { name: "Eve", age: 19 } ] },
	{ name: "Liz", friends: [ { name: "Max", age: 16 } ] }
]
RETURN users[* RETURN {
	name: .name,
	hasAdult: .friends[? FILTER .age > 18]
}]`, []any{
			map[string]any{"name": "Ann", "hasAdult": true},
			map[string]any{"name": "Dan", "hasAdult": true},
			map[string]any{"name": "Liz", "hasAdult": false},
		}),
		Array(`
LET values = [1, 2, 3, 4]
RETURN values[* FILTER . > 2 RETURN . * 10]`, []any{30, 40}),
		S(`
LET arr = [1, 2, 3, 4]
RETURN arr[? 2 FILTER . % 2 == 0]`, true),
		S(`
LET arr = [1, 3, 5]
RETURN arr[? FILTER . % 2 == 0]`, false),
		S(`
LET arr = [1]
RETURN arr[?]`, true),
		S(`
LET arr = []
RETURN arr[?]`, false),
		S(`
LET arr = [1, 2]
RETURN arr[? ANY FILTER . > 1]`, true),
		S(`
LET arr = [1, 3, 5]
RETURN arr[? NONE FILTER . % 2 == 0]`, true),
		S(`
LET arr = [2, 4]
RETURN arr[? ALL FILTER . % 2 == 0]`, true),
		S(`
LET arr = [2, 4, 6]
RETURN arr[? AT LEAST (2) FILTER . % 2 == 0]`, true),
		S(`
LET arr = [2, 4, 6]
RETURN arr[? 2..3 FILTER . % 2 == 0]`, true),
		Array(`
LET docs = [
	{ name: "A", dimensions: [{ type: "height", value: 45 }] },
	{ name: "B", dimensions: [{ type: "height", value: 35 }] },
	{ name: "C", dimensions: [{ type: "width", value: 50 }] }
]

FOR doc IN docs
	FILTER doc.dimensions[? FILTER .type == "height" AND .value > 40]
	RETURN doc.name`, []any{"A"}),
		Array(`
LET docs = [
	{ name: "A", dimensions: [{ part: "frame", measurements: [{ type: "width", value: 80 }] }] },
	{ name: "B", dimensions: [{ part: "frame", measurements: [{ type: "width", value: 90 }] }] },
	{ name: "C", dimensions: [{ part: "wheel", measurements: [{ type: "width", value: 70 }] }] }
]

FOR doc IN docs
	FILTER doc.dimensions[? FILTER .part == "frame" AND
		.measurements[? FILTER .type == "width" AND .value <= 80]]
	RETURN doc.name`, []any{"A"}),
		Array(`
LET users = [
	{ name: "Ann", age: 20 },
	{ name: "Bob", age: 30 }
]
RETURN users[* FILTER .age > 20].name`, []any{"Bob"}),
		Array(`
LET users = [
	{ name: "Ann", age: 20 },
	{ name: "Bob", age: 30 }
]
RETURN users[* RETURN .name]`, []any{"Ann", "Bob"}),
		Array(`
LET users = [
	{ age: 10 },
	{ name: "Bob" }
]
RETURN users[* RETURN ?.name]`, []any{nil, "Bob"}),
		Array(`
LET users = [
	[1, 2],
	[3]
]
RETURN users[* RETURN .[0]]`, []any{1, 3}),
		Array(`
LET users = [
	{ name: "Ann", age: 20 },
	{ name: "Bob", age: 30 }
]
RETURN users[* FILTER .age > 20][* RETURN .name]`, []any{"Bob"}),
		Array(`
LET users = [
	[{ active: true }, { active: false }],
	[{ active: false }]
]
RETURN users[* FILTER .[? FILTER .active]]`, []any{
			[]any{
				map[string]any{"active": true},
				map[string]any{"active": false},
			},
		}),
		Array(`
LET users = [
	[{ name: "Ann" }],
	[{ name: "Bob" }, { name: "Cat" }]
]
RETURN users[* RETURN .[*].name]`, []any{
			[]any{"Ann"},
			[]any{"Bob", "Cat"},
		}),
		Array(`
LET groups = [
	[[1, 2], [3]],
	[[4]]
]
RETURN groups[* RETURN .[**]]`, []any{
			[]any{1, 2, 3},
			[]any{4},
		}),
		Array(`
LET users = [
	{ name: "Ann" },
	{ name: "Bob" },
	{ name: "Cat" }
]
RETURN users[* LIMIT 2].name`, []any{"Ann", "Bob"}),
		Array(`
LET users = [
	{ name: "Ann" },
	{ name: "Bob" }
]
RETURN users[* RETURN { n: .name }].n`, []any{"Ann", "Bob"}),
		Array(`
LET groups = [
	[{ name: "Ann", age: 20 }],
	[{ name: "Bob", age: 30 }]
]
RETURN groups[** FILTER .age > 20].name`, []any{"Bob"}),
		Array(`
LET users = [
						{ 
							name: "john", 
							age: 30,
							friends: [	
								{ name: "tina", age: 43 },
								{ name: "tom", age: 35 },
								{ name: "helga", age: 52 }
							]
						},
						{
							name: "sandra", 
							age: 25,
							friends: [
								{ name: "elena", age: 48 },
								{ name: "maria", age: 38 }
							]
						}
					]

			FOR u IN users
				RETURN {
					name: u.name,
					friends: u.friends[* FILTER CONTAINS(.name, "a") AND .age > 40
						LIMIT 2
						RETURN CONCAT(.name, " is ", .age)
					]
				}
				`,
			[]any{map[string]any{
				"name": "john",
				"friends": []any{
					"tina is 43",
					"helga is 52",
				}},
				map[string]any{
					"name": "sandra",
					"friends": []any{
						"elena is 48",
					}},
			}),
		Error(`
					LET value = 1

					RETURN value[**]
				`),
	})
}
