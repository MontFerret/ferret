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

func BenchmarkLetDestructuring_None(b *testing.B) {
	RunBenchmarkNone(b, letDestructuringQuery)
}

func BenchmarkLetDestructuring_Full(b *testing.B) {
	RunBenchmarkFull(b, letDestructuringQuery)
}

func BenchmarkLetManualExtraction_None(b *testing.B) {
	RunBenchmarkNone(b, letManualExtractionQuery)
}

func BenchmarkLetManualExtraction_Full(b *testing.B) {
	RunBenchmarkFull(b, letManualExtractionQuery)
}

func BenchmarkVarDestructuring_None(b *testing.B) {
	RunBenchmarkNone(b, varDestructuringQuery)
}

func BenchmarkVarDestructuring_Full(b *testing.B) {
	RunBenchmarkFull(b, varDestructuringQuery)
}

func BenchmarkVarManualExtraction_None(b *testing.B) {
	RunBenchmarkNone(b, varManualExtractionQuery)
}

func BenchmarkVarManualExtraction_Full(b *testing.B) {
	RunBenchmarkFull(b, varManualExtractionQuery)
}

func BenchmarkForDestructuring_None(b *testing.B) {
	RunBenchmarkNone(b, forDestructuringQuery)
}

func BenchmarkForDestructuring_Full(b *testing.B) {
	RunBenchmarkFull(b, forDestructuringQuery)
}

func BenchmarkForManualExtraction_None(b *testing.B) {
	RunBenchmarkNone(b, forManualExtractionQuery)
}

func BenchmarkForManualExtraction_Full(b *testing.B) {
	RunBenchmarkFull(b, forManualExtractionQuery)
}

func BenchmarkLetIgnoredSubtree_None(b *testing.B) {
	RunBenchmarkNone(b, letIgnoredSubtreeQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetIgnoredSubtree_Full(b *testing.B) {
	RunBenchmarkFull(b, letIgnoredSubtreeQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetDirectIgnore_None(b *testing.B) {
	RunBenchmarkNone(b, letDirectIgnoreQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}

func BenchmarkLetDirectIgnore_Full(b *testing.B) {
	RunBenchmarkFull(b, letDirectIgnoreQuery, WithParam("value", ignoredSubtreeBenchmarkValue))
}
