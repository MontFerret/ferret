package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
)

const compilerMemberPipelineQuery = `
LET users = @users
RETURN users[* RETURN {
	name: .name,
	friends: .friends[* FILTER .age > 18][* RETURN .name]
}][* LIMIT 2].friends
`

const compilerMatchQueryMixQuery = `
LET doc = @doc
RETURN MATCH QUERY COUNT ".items" IN doc USING css (
	0 => "empty",
	count WHEN count > 2 => UPPER(QUERY VALUE ".featured" IN doc USING css ON ERROR RETURN "fallback"),
	_ => QUERY ANY ".items" IN doc USING css,
)
`

func BenchmarkCompilerCompileMemberPipeline_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerMemberPipelineQuery, compiler.O0)
}

func BenchmarkCompilerCompileMemberPipeline_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerMemberPipelineQuery, compiler.O1)
}

func BenchmarkCompilerCompileMatchQueryMix_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerMatchQueryMixQuery, compiler.O0)
}

func BenchmarkCompilerCompileMatchQueryMix_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerMatchQueryMixQuery, compiler.O1)
}
