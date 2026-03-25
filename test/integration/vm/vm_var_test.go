package vm_test

import (
	"testing"

	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestVarBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`
	VAR x = 1
	x = x + 1
	RETURN x
	`, 2, "Top-level VAR can be reassigned"),
		S(`
	VAR x = 1
	x += 2
	RETURN x
	`, 3, "Top-level VAR supports +="),
		S(`
	VAR x = 5
	x -= 2
	RETURN x
	`, 3, "Top-level VAR supports -="),
		S(`
	VAR x = 3
	x *= 4
	RETURN x
	`, 12, "Top-level VAR supports *="),
		S(`
	VAR x = 12
	x /= 3
	RETURN x
	`, 4, "Top-level VAR supports /="),
		S(`
	VAR x = 1
	x = "hello"
	RETURN x
`, "hello", "Top-level VAR can be reassigned across types"),
		S(`
VAR CURRENT = 1
CURRENT = CURRENT + 1
RETURN CURRENT
`, 2, "Safe reserved words remain valid mutable binding names"),
		S(`
FUNC run() (
  VAR total = 1
  total = total + 2
  RETURN total
)
RETURN run()
`, 3, "Function-block VAR can be reassigned"),
		Array(`
FOR item IN [1, 2]
  VAR current = item
  current = current + 1
  RETURN current
`, []any{2, 3}, "FOR body VAR can be reassigned"),
		Array(`
VAR total = 1
FUNC bump() (
  total = total + 1
  RETURN total
)
LET after = bump()
RETURN [after, total]
`, []any{2, 2}, "Nested UDF assignment mutates the captured outer VAR"),
		Array(`
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
		Array(`
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
		Array(`
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
		S(`
	FUNC run() (
	  VAR total = 1
	  FUNC setTotal(v) (
	    total = v
	    RETURN 10
	  )
	  total += setTotal(5)
	  RETURN total
	)
	RETURN run()
	`, 11, "Compound assignment snapshots the old VAR value before RHS side effects"),
		Array(`
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
