package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestDirectMutation(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(`
VAR count = 1
count += 1
RETURN count
`, 2, "VAR binding supports numeric +="),
		S(`
VAR total = 8
total -= 3
total *= 2
total /= 5
RETURN total
`, 2, "VAR binding supports numeric -=, *=, /="),
		S(`
VAR name = "Tim"
name += " Jr."
RETURN name
`, "Tim Jr.", "VAR binding supports string +="),
		Object(`
LET user = { name: "Alice", profile: { city: "Riga" } }
user.name = "Tim"
user.profile.city = "Berlin"
RETURN user
`, map[string]any{
			"name": "Tim",
			"profile": map[string]any{
				"city": "Berlin",
			},
		}, "Object field assignment replaces nested fields"),
		Object(`
LET user = { name: "Tim" }
user.nickname = "tim"
RETURN user
`, map[string]any{
			"name":     "Tim",
			"nickname": "tim",
		}, "Object field assignment creates missing final fields"),
		Object(`
LET obj = {}
LET key = "theme"
obj[key] = "dark"
RETURN obj
`, map[string]any{
			"theme": "dark",
		}, "Dynamic keyed assignment creates final object entries"),
		Array(`
LET arr = [1, 2, 3]
arr[1] = 20
RETURN arr
`, []any{1, 20, 3}, "Array index assignment replaces existing elements"),
		Object(`
LET data = { items: [{ name: "a", count: 1 }] }
data.items[0].name = "x"
data.items[0].count += 1
RETURN data
`, map[string]any{
			"items": []any{
				map[string]any{
					"name":  "x",
					"count": 2,
				},
			},
		}, "Mixed member and index paths support read-modify-write"),
		Object(`
LET user = { profile: { city: "Riga" } }
user?.profile?.city = "Berlin"
RETURN user
`, map[string]any{
			"profile": map[string]any{
				"city": "Berlin",
			},
		}, "Safe member assignment writes when every guarded hop is present"),
		Object(`
LET obj = { items: [{ count: 1 }] }
LET i = 0
obj?.items?.[i].count += 1
RETURN obj
`, map[string]any{
			"items": []any{
				map[string]any{
					"count": 2,
				},
			},
		}, "Safe computed assignment uses existing safe read syntax"),
		Object(`
LET user = { count: 1 }
FUNC inc() (
  user.count += 1
  RETURN 0
)
LET _ = inc()
RETURN user
`, map[string]any{
			"count": 2,
		}, "UDF path mutation captures the root without rebinding it"),
		Nil(`
LET user = NONE
user?.profile?.city = "Berlin"
RETURN user
`, "Safe assignment no-ops on missing root"),
		S(`
LET user = NONE
user?.profile?.city = FAIL()
RETURN 1
`, 1, "Safe assignment skips RHS on missing root"),
		S(`
LET user = {}
user?.profile?.city = FAIL()
RETURN 1
`, 1, "Safe assignment skips RHS on missing guarded intermediate"),
		S(`
LET user = NONE
user?.profile.count += FAIL()
RETURN 1
`, 1, "Safe augmented assignment skips RHS on missing guarded root"),
		S(`
LET user = { profile: {} }
user?.profile?.count += FAIL()
RETURN 1
`, 1, "Safe augmented assignment skips RHS on missing final value"),
		Error(`
LET user = {}
user.profile.city = "Paris"
RETURN user
`, "Strict missing intermediate fails"),
		Error(`
LET user = {}
user?.profile.city = "Paris"
RETURN user
`, "Mixed safe and strict path fails when a strict hop is missing"),
		Error(`
LET arr = [1, 2, 3]
arr[100] = 1
RETURN arr
`, "Array assignment rejects out-of-bounds indexes"),
		Error(`
LET arr = [1, 2, 3]
arr[-1] = 1
RETURN arr
`, "Array assignment rejects negative indexes"),
		Error(`
LET arr = [1, 2, 3]
arr["x"] = 1
RETURN arr
`, "Array assignment rejects non-integer indexes"),
		Error(`
LET arr = [1, 2, 3]
arr[1.5] = 1
RETURN arr
`, "Array assignment rejects float indexes"),
		Error(`
LET text = "abc"
text[0] = "x"
RETURN text
`, "Assignment rejects unsupported target mutation"),
	}, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("should not execute")
	}))
}
