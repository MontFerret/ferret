package vm_test

import "testing"

func TestVarBindings(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`
VAR x = 1
x = x + 1
RETURN x
`, 2, "Top-level VAR can be reassigned"),
		Case(`
VAR x = 1
x = "hello"
RETURN x
`, "hello", "Top-level VAR can be reassigned across types"),
		Case(`
VAR CURRENT = 1
CURRENT = CURRENT + 1
RETURN CURRENT
`, 2, "Safe reserved words remain valid mutable binding names"),
		Case(`
FUNC run() (
  VAR total = 1
  total = total + 2
  RETURN total
)
RETURN run()
`, 3, "Function-block VAR can be reassigned"),
		CaseArray(`
FOR item IN [1, 2]
  VAR current = item
  current = current + 1
  RETURN current
`, []any{2, 3}, "FOR body VAR can be reassigned"),
		CaseArray(`
VAR total = 1
FUNC bump() (
  total = total + 1
  RETURN total
)
LET after = bump()
RETURN [after, total]
`, []any{2, 2}, "Nested UDF assignment mutates the captured outer VAR"),
		CaseArray(`
FUNC run() (
  VAR carried = 0
  FUNC setCarried(v) (
    carried = v
    RETURN carried
  )

  RETURN (
    FOR item IN [3, 1]
      LET _ = setCarried(item * 10)
      SORT item
      RETURN { item, carried }
  )
)

RETURN run()
`, []any{
			map[string]any{"item": 1, "carried": 10},
			map[string]any{"item": 3, "carried": 30},
		}, "SORT snapshots the current value of a promoted VAR"),
		CaseArray(`
FUNC run() (
  VAR carried = 0
  FUNC setCarried(v) (
    carried = v
    RETURN carried
  )

  RETURN (
    FOR item IN [1, 2, 3]
      LET _ = setCarried(item * 10)
      COLLECT parity = item % 2 INTO groups KEEP carried
      RETURN { parity, groups }
  )
)

RETURN run()
`, []any{
			map[string]any{
				"parity": 0,
				"groups": []any{
					map[string]any{"carried": 20},
				},
			},
			map[string]any{
				"parity": 1,
				"groups": []any{
					map[string]any{"carried": 10},
					map[string]any{"carried": 30},
				},
			},
		}, "COLLECT KEEP snapshots the current value of a promoted VAR"),
		CaseArray(`
FUNC run() (
  VAR carried = 0
  FUNC setCarried(v) (
    carried = v
    RETURN carried
  )

  RETURN (
    FOR item IN [1, 2, 3]
      LET _ = setCarried(item * 10)
      COLLECT parity = item % 2 INTO groups = { carried: carried }
      RETURN { parity, groups }
  )
)

RETURN run()
`, []any{
			map[string]any{
				"parity": 0,
				"groups": []any{
					map[string]any{"carried": 20},
				},
			},
			map[string]any{
				"parity": 1,
				"groups": []any{
					map[string]any{"carried": 10},
					map[string]any{"carried": 30},
				},
			},
		}, "COLLECT INTO projection snapshots the current value of a promoted VAR"),
		CaseArray(`
FUNC outer() (
  VAR total = 1
  FUNC middle(v) (
    FUNC inner() => total
    total = total + v
    RETURN inner()
  )
  RETURN [middle(2), total]
)
RETURN outer()
`, []any{3, 3}, "Deeply nested UDF reads updated VAR from grandparent scope"),
	})
}
