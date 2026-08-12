package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestUDFReturnedForResult(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
FUNC project(items) {
  RETURN FOR item IN items {
    RETURN item * 2
  }
}
RETURN project([])
`, []any{}, "empty input"),
		Array(`
FUNC project(items) {
  RETURN FOR item IN items {
    RETURN item * 2
  }
}
RETURN project([3])
`, []any{6}, "single input"),
		Array(`
FUNC project(items) {
  RETURN FOR item IN items {
    FILTER item > 1
    SORT item DESC
    LIMIT 2
    RETURN item * 2
  }
}
RETURN project([1, 4, 2, 3])
`, []any{8, 6}, "filter sort and limit"),
		Array(`
FUNC countByParity(items) {
  RETURN FOR item IN items {
    COLLECT parity = item % 2 WITH COUNT INTO count
    SORT parity
    RETURN [parity, count]
  }
}
RETURN countByParity([1, 2, 3, 4])
`, []any{[]any{0, 2}, []any{1, 2}}, "collect"),
		Array(`
FUNC expand(items) {
  RETURN FOR item IN items {
    LET expanded = (
      FOR value IN [item, item + 10]
        RETURN value
    )
    RETURN expanded
  }
}
RETURN expand([1, 2])
`, []any{[]any{1, 11}, []any{2, 12}}, "nested mixed loop styles"),
		Array(`
LET item = 100
FUNC project(items) {
  RETURN FOR item IN items {
    RETURN item + 1
  }
}
RETURN [project([1, 2]), item]
`, []any{[]any{2, 3}, 100}, "loop binding is not an outer capture"),
		Array(`
VAR total = 0
FUNC runningTotals(items) {
  RETURN FOR item IN items {
    total += item
    RETURN total
  }
}
RETURN [runningTotals([1, 2]), total]
`, []any{[]any{1, 3}, 3}, "mutable outer cell capture"),
		Array(`
FUNC generate() {
  VAR value = 0
  RETURN FOR WHILE value < 3 {
    value += 1
    RETURN value
  }
}
RETURN generate()
`, []any{1, 2, 3}, "returned WHILE loop"),
		Array(`
FUNC once() {
  RETURN FOR value DO WHILE false {
    RETURN value
  }
}
RETURN once()
`, []any{0}, "returned DO WHILE loop"),
		Array(`
LET items = [1, 2, 3]
FUNC project(values) {
  RETURN FOR value IN values {
    RETURN value * 2
  }
}
RETURN [project(items), (FOR value IN items RETURN value * 2)]
`, []any{[]any{2, 4, 6}, []any{2, 4, 6}}, "matches equivalent top-level loop result"),
		Error(`
FUNC divide(items) {
  RETURN FOR item IN items {
    RETURN 10 / item
  }
}
RETURN divide([2, 0])
`, "runtime errors propagate"),
	})
}
