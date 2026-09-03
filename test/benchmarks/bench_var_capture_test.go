package benchmarks_test

import "testing"

const (
	varCaptureSingleCellQuery = `
FUNC run() {
  VAR total = 0
  FUNC bump(v) {
    total = total + v
    RETURN total
  }

  RETURN (
    FOR i IN 1..500
      RETURN bump(i)
  )
}
RETURN run()
`

	varCaptureMultiCellQuery = `
FUNC run() {
  VAR a = 0
  VAR b = 1
  VAR c = 2
  VAR d = 3

  FUNC bump(v) {
    a = a + v
    b = b + v
    c = c + v
    d = d + v
    RETURN a + b + c + d
  }

  RETURN (
    FOR i IN 1..500
      RETURN bump(i)
  )
}
RETURN run()
`
)

func BenchmarkVarCapture_SingleCell_None(b *testing.B) {
	RunBenchmarkNone(b, varCaptureSingleCellQuery)
}

func BenchmarkVarCapture_SingleCell_Basic(b *testing.B) {
	RunBenchmarkBasic(b, varCaptureSingleCellQuery)
}

func BenchmarkVarCapture_SingleCell_Full(b *testing.B) {
	RunBenchmarkFull(b, varCaptureSingleCellQuery)
}

func BenchmarkVarCapture_MultiCell_None(b *testing.B) {
	RunBenchmarkNone(b, varCaptureMultiCellQuery)
}

func BenchmarkVarCapture_MultiCell_Basic(b *testing.B) {
	RunBenchmarkBasic(b, varCaptureMultiCellQuery)
}

func BenchmarkVarCapture_MultiCell_Full(b *testing.B) {
	RunBenchmarkFull(b, varCaptureMultiCellQuery)
}
