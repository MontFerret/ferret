package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/test/integration/base"
)

func TestCollectAggregate(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
LET users = [
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
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

  RETURN {
    genderGroup,
    minAge,
    maxAge
  }
`, []any{
			map[string]any{"genderGroup": "f", "minAge": 25, "maxAge": 45},
			map[string]any{"genderGroup": "m", "minAge": 31, "maxAge": 69},
		}, "Should collect and aggregate values by a single key"),
		CaseArray(`
			LET users = [
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
			FOR u IN users
  				COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
    				minAge, 
    				maxAge 
  				}
		`, []any{
			map[string]any{"minAge": 25, "maxAge": 69},
		}, "Should collect and aggregate values without grouping"),
		SkipCaseArray(`
LET users = [
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
FOR u IN users
  COLLECT genderGroup = u.gender
   AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)

  RETURN {
    genderGroup,
    minAge,
    maxAge
  }
`, []any{
			map[string]any{"genderGroup": "f", "minAge": 25, "maxAge": 45},
			map[string]any{"genderGroup": "m", "minAge": 31, "maxAge": 69},
		}, "Should collect and aggregate values by a single key"),
	})
}
