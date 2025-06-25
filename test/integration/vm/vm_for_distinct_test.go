package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestForDistinct(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(
			`FOR i IN [ 1, 2, 3, 4, 1, 3 ]
							RETURN DISTINCT i
		`,
			[]any{1, 2, 3, 4},
		),
		CaseArray(
			`FOR i IN ["foo", "bar", "qaz", "foo", "abc", "bar"]
							RETURN DISTINCT i
		`,
			[]any{"foo", "bar", "qaz", "abc"},
		),
		CaseArray(
			`FOR i IN [["foo"], ["bar"], ["qaz"], ["foo"], ["abc"], ["bar"]]
							RETURN DISTINCT i
		`,
			[]any{[]any{"foo"}, []any{"bar"}, []any{"qaz"}, []any{"abc"}},
		),
		CaseArray(`
LET strs = ["foo", "bar", "qaz", "foo", "abc", "bar"]

FOR s IN strs
	SORT s
	RETURN DISTINCT s
`, []any{"abc", "bar", "foo", "qaz"}, "Should sort and respect DISTINCT keyword"),
		CaseArray(
			`
			FOR i IN [ 1, 1, 2, 3, 4, 1, 3 ]
				LIMIT 2
				RETURN DISTINCT i
		`,
			[]any{1}),
		CaseArray(
			`
			FOR i IN [ 1, 1, 1, 3, 4, 1, 3 ]
				LIMIT 1, 2
				RETURN DISTINCT i
		`,
			[]any{1}),
		CaseArray(
			`
			FOR i IN [ 1, 2, 3, 4, 1, 3, 3, 4 ]
				FILTER i > 2
				RETURN DISTINCT i
		`,
			[]any{3, 4},
		),
		CaseArray(
			`LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR i IN users
				COLLECT gender = i.gender, age = i.age
				RETURN DISTINCT {gender}
		`, []any{
				map[string]any{"gender": "f"},
				map[string]any{"gender": "m"},
			}),
		CaseArray(
			`LET users = [
				{
					active: true,
					age: 31,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 25,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 36,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 69,
					gender: "m",
					married: true
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				}
			]
			FOR i IN users
				COLLECT gender = i.gender INTO genders = { active: i.active }
				RETURN DISTINCT genders[0]
		`, []any{
				map[string]any{"active": true},
			}),
		CaseArray(`
LET users = [
				{
					active: true,
					age: 39,
					gender: "f",
					married: false
				},
				{
					active: true,
					age: 45,
					gender: "f",
					married: true
				},
				{
					active: true,
					age: 39,
					gender: "m",
					married: false
				},
				{
					active: false,
					age: 45,
					gender: "m",
					married: true
				}
			]
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

  RETURN DISTINCT {
	minAge,
    maxAge
  }
`, []any{
			map[string]any{"maxAge": 45, "minAge": 39},
		}, "Should collect and aggregate values by a single key"),
	})
}
