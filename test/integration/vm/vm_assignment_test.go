package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
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
		Object(`
LET item = { id: 1, deprecated: true }
DELETE item.deprecated
RETURN item
`, map[string]any{
			"id": 1,
		}, "DELETE removes an object property"),
		Object(`
LET key = "debug"
LET payload = { id: 1, debug: { trace: true } }
DELETE payload[key]
RETURN payload
`, map[string]any{
			"id": 1,
		}, "DELETE removes a computed key"),
		Object(`
LET item = { meta: { deprecated: true, keep: true } }
DELETE item.meta.deprecated
RETURN item
`, map[string]any{
			"meta": map[string]any{
				"keep": true,
			},
		}, "DELETE removes the final member from a nested parent"),
		Object(`
LET item = {}
DELETE item.missing
RETURN item
`, map[string]any{}, "DELETE missing property is a no-op"),
		Object(`
LET assigned = { foo: 1 }
assigned.foo = NONE
LET deleted = { foo: 1 }
DELETE deleted.foo
RETURN { assigned: assigned, deleted: deleted }
`, map[string]any{
			"assigned": map[string]any{
				"foo": nil,
			},
			"deleted": map[string]any{},
		}, "DELETE removes a property instead of assigning NONE"),
		Object(`
LET item = {}
DELETE item.meta?.deprecated
RETURN item
`, map[string]any{}, "Safe DELETE no-ops on missing optional parent"),
		Nil(`
LET value = NONE
DELETE value?.foo
RETURN value
`, "Safe DELETE no-ops on NONE parent"),
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
LET user = { profile: {} }
user?.profile?.count = FAIL()
RETURN 1
`, 1, "Safe assignment skips RHS on missing final member"),
		S(`
LET arr = []
arr?.[0] = FAIL()
RETURN 1
`, 1, "Safe assignment skips RHS on missing final index"),
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
		Error(`
LET value = NONE
DELETE value.foo
RETURN value
`, "Strict DELETE from NONE parent errors"),
		Error(`
LET value = 42
DELETE value.foo
RETURN value
`, "DELETE from non-object parent errors"),
		Array(`
LET arr = [1, 2, 3]
DELETE arr[1]
RETURN arr
`, []any{1, 3}, "DELETE removes an array index"),
		Array(`
LET arr = [1, 2, 3]
DELETE arr[100]
RETURN arr
`, []any{1, 2, 3}, "DELETE missing array index is a no-op"),
		Error(`
LET arr = [1, 2, 3]
DELETE arr[-1]
RETURN arr
`, "DELETE rejects negative array index"),
		Error(`
LET arr = [1, 2, 3]
DELETE arr[1.5]
RETURN arr
`, "DELETE rejects float array index"),
		Error(`
LET arr = [1, 2, 3]
DELETE arr["1"]
RETURN arr
`, "DELETE rejects non-integer array index"),
	}, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, errors.New("should not execute")
	}))
}

func TestDirectMutationProxyRemovable(t *testing.T) {
	RunSpecFactory(t, func() []spec.Spec {
		return []spec.Spec{
			Array(`
LET arr = @arr
DELETE arr[1]
RETURN arr
`, []any{1, 3}, "DELETE removes a proxy slice index").Env(vm.WithParam("arr", sdk.NewProxySlice([]int{1, 2, 3}))),
			Object(`
LET obj = @obj
DELETE obj.one
RETURN obj
`, map[string]any{
				"two": 2,
			}, "DELETE removes a proxy map key").Env(vm.WithParam("obj", sdk.NewProxyMap(map[string]int{
				"one": 1,
				"two": 2,
			}))),
		}
	})
}
