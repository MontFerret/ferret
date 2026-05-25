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
	count WHEN count > 2 => UPPER(QUERY ONE ".featured" IN doc USING css),
	_ => QUERY ONE ".items" IN doc USING css,
)
`

const compilerQueryShorthandQuery = "\n" +
	"LET doc = @doc\n" +
	"RETURN {\n" +
	"  title: doc[~? css`h1`],\n" +
	"  cards: doc[~ css`.product-card`],\n" +
	"  next: doc[~? css`[data-testid=\"next-page\"]`],\n" +
	"  labels: doc[~ css`.product-card`][* RETURN .[~? css`.title`]]\n" +
	"}\n"

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

func BenchmarkCompilerCompileQueryShorthand_O0(b *testing.B) {
	benchmarkCompileQuery(b, compilerQueryShorthandQuery, compiler.O0)
}

func BenchmarkCompilerCompileQueryShorthand_O1(b *testing.B) {
	benchmarkCompileQuery(b, compilerQueryShorthandQuery, compiler.O1)
}
