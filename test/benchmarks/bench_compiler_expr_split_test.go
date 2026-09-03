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
RETURN MATCH QUERY COUNT ".items" IN doc USING css {
	0 => "empty",
	count WHEN count > 2 => UPPER(QUERY ONE ".featured" IN doc USING css),
	_ => QUERY ONE ".items" IN doc USING css,
}
`

const compilerQueryShorthandQuery = "\n" +
	"LET doc = @doc\n" +
	"RETURN {\n" +
	"  title: doc[~? css`h1`],\n" +
	"  cards: doc[~ css`.product-card`],\n" +
	"  next: doc[~? css`[data-testid=\"next-page\"]`],\n" +
	"  labels: doc[~ css`.product-card`][* RETURN .[~? css`.title`]]\n" +
	"}\n"

func BenchmarkCompilerCompileMemberPipeline_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerMemberPipelineQuery, compiler.None)
}

func BenchmarkCompilerCompileMemberPipeline_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerMemberPipelineQuery, compiler.Basic)
}

func BenchmarkCompilerCompileMemberPipeline_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerMemberPipelineQuery, compiler.Full)
}

func BenchmarkCompilerCompileMatchQueryMix_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerMatchQueryMixQuery, compiler.None)
}

func BenchmarkCompilerCompileMatchQueryMix_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerMatchQueryMixQuery, compiler.Basic)
}

func BenchmarkCompilerCompileMatchQueryMix_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerMatchQueryMixQuery, compiler.Full)
}

func BenchmarkCompilerCompileQueryShorthand_None(b *testing.B) {
	benchmarkCompileQuery(b, compilerQueryShorthandQuery, compiler.None)
}

func BenchmarkCompilerCompileQueryShorthand_Basic(b *testing.B) {
	benchmarkCompileQuery(b, compilerQueryShorthandQuery, compiler.Basic)
}

func BenchmarkCompilerCompileQueryShorthand_Full(b *testing.B) {
	benchmarkCompileQuery(b, compilerQueryShorthandQuery, compiler.Full)
}
