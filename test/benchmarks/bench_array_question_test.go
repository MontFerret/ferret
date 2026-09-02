package benchmarks_test

import "testing"

const bareArrayQuestionQuery = `RETURN @arr[?]`

const bareLocalArrayQuestionQuery = `
LET arr = [1, 2, 3, 4, 5, 6, 7, 8]
RETURN arr[?]`

var bareArrayQuestionValues = []any{1, 2, 3, 4, 5, 6, 7, 8}

func BenchmarkBareArrayQuestion_None(b *testing.B) {
	RunBenchmarkNone(b, bareArrayQuestionQuery, WithParam("arr", bareArrayQuestionValues))
}

func BenchmarkBareArrayQuestion_Full(b *testing.B) {
	RunBenchmarkFull(b, bareArrayQuestionQuery, WithParam("arr", bareArrayQuestionValues))
}

func BenchmarkBareLocalArrayQuestion_None(b *testing.B) {
	RunBenchmarkNone(b, bareLocalArrayQuestionQuery)
}

func BenchmarkBareLocalArrayQuestion_Full(b *testing.B) {
	RunBenchmarkFull(b, bareLocalArrayQuestionQuery)
}
