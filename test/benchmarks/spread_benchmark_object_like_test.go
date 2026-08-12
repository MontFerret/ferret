package benchmarks_test

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type benchmarkObjectLike struct {
	runtime.Map
}

func (*benchmarkObjectLike) ObjectLike() {}
