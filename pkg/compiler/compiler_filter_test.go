package compiler_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestForFilter(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 2
				RETURN i
		`,
			[]any{3, 4, 3},
			ShouldEqualJSON,
		},
		{
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3 ]
				FILTER i > 1 AND i < 4
				RETURN i
		`,
			[]any{2, 3, 3},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					age: 31,
					gender: "m",
					name: "Josh"
				},
				{
					age: 29,
					gender: "f",
					name: "Mary"
				},
				{
					age: 36,
					gender: "m",
					name: "Peter"
				}
			]
			FOR u IN users
				FILTER u.name =~ "r"
				RETURN u
		`,
			[]any{map[string]any{"age": 29, "gender": "f", "name": "Mary"}, map[string]any{"age": 36, "gender": "m", "name": "Peter"}},
			ShouldEqualJSON,
		},
		{
			`
					LET users = [
						{
							active: true,
							age: 31,
							gender: "m"
						},
						{
							active: true,
							age: 29,
							gender: "f"
						},
						{
							active: true,
							age: 36,
							gender: "m"
						}
					]
					FOR u IN users
						FILTER u.active == true
						FILTER u.age < 35
						RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active
				RETURN u
				`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}, map[string]any{"active": true, "gender": "f", "age": 29}, map[string]any{"active": true, "gender": "m", "age": 36}},
			ShouldEqualJSON,
		},
		{
			`
			LET users = [
				{
					active: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					age: 29,
					gender: "f"
				},
				{
					active: true,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					age: 69,
					gender: "m"
				}
			]
			FOR u IN users
				FILTER u.active == true
				LIMIT 2
				FILTER u.gender == "m"
				RETURN u
		`,
			[]any{map[string]any{"active": true, "gender": "m", "age": 31}},
			ShouldEqualJSON,
		},
	})
}
