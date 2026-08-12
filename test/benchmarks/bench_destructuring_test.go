package benchmarks_test

import "testing"

const (
	letDestructuringQuery = `
LET { profile: { score }, values: [first, _, third] } = {
    profile: { score: 10 },
    values: [1, 2, 3]
}
RETURN score + first + third
`

	letManualExtractionQuery = `
LET source = {
    profile: { score: 10 },
    values: [1, 2, 3]
}
LET score = source.profile.score
LET first = source.values[0]
LET third = source.values[2]
RETURN score + first + third
`

	varDestructuringQuery = `
VAR { count, step } = { count: 1, step: 2 }
count += step
step += 1
RETURN count + step
`

	varManualExtractionQuery = `
LET source = { count: 1, step: 2 }
VAR count = source.count
VAR step = source.step
count += step
step += 1
RETURN count + step
`

	forDestructuringQuery = `
RETURN FOR { left, right } IN [
    { left: 1, right: 2 },
    { left: 3, right: 4 },
    { left: 5, right: 6 },
    { left: 7, right: 8 }
]
    RETURN left + right
`

	forManualExtractionQuery = `
RETURN FOR item IN [
    { left: 1, right: 2 },
    { left: 3, right: 4 },
    { left: 5, right: 6 },
    { left: 7, right: 8 }
]
    LET left = item.left
    LET right = item.right
    RETURN left + right
`

	letIgnoredSubtreeQuery = `
LET { kept, ignored: { nested: [_, _] } } = @value
RETURN kept
`

	letDirectIgnoreQuery = `
LET { kept, ignored: _ } = @value
RETURN kept
`
)

var ignoredSubtreeBenchmarkValue = map[string]any{
	"kept": 42,
	"ignored": map[string]any{
		"nested": []any{1, 2},
	},
}

func BenchmarkLetDestructuring_O0(b *testing.B) {
	RunBenchmarkO0(b, letDestructuringQuery)
}

func BenchmarkLetDestructuring_O1(b *testing.B) {
	RunBenchmarkO1(b, letDestructuringQuery)
}

func BenchmarkLetManualExtraction_O0(b *testing.B) {
	RunBenchmarkO0(b, letManualExtractionQuery)
}

func BenchmarkLetManualExtraction_O1(b *testing.B) {
	RunBenchmarkO1(b, letManualExtractionQuery)
}

func BenchmarkVarDestructuring_O0(b *testing.B) {
	RunBenchmarkO0(b, varDestructuringQuery)
}

func BenchmarkVarDestructuring_O1(b *testing.B) {
	RunBenchmarkO1(b, varDestructuringQuery)
}

func BenchmarkVarManualExtraction_O0(b *testing.B) {
	RunBenchmarkO0(b, varManualExtractionQuery)
}

func BenchmarkVarManualExtraction_O1(b *testing.B) {
	RunBenchmarkO1(b, varManualExtractionQuery)
}

func BenchmarkForDestructuring_O0(b *testing.B) {
	RunBenchmarkO0(b, forDestructuringQuery)
}

func BenchmarkForDestructuring_O1(b *testing.B) {
	RunBenchmarkO1(b, forDestructuringQuery)
}

func BenchmarkForManualExtraction_O0(b *testing.B) {
	RunBenchmarkO0(b, forManualExtractionQuery)
}

func BenchmarkForManualExtraction_O1(b *testing.B) {
	RunBenchmarkO1(b, forManualExtractionQuery)
}

func BenchmarkLetIgnoredSubtree_O0(b *testing.B) {
	RunBenchmarkO0(b, letIgnoredSubtreeQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetIgnoredSubtree_O1(b *testing.B) {
	RunBenchmarkO1(b, letIgnoredSubtreeQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetDirectIgnore_O0(b *testing.B) {
	RunBenchmarkO0(b, letDirectIgnoreQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetDirectIgnore_O1(b *testing.B) {
	RunBenchmarkO1(b, letDirectIgnoreQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}
