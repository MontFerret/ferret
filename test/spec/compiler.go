package spec

import "github.com/MontFerret/ferret/v2/pkg/compiler"

func mustNewCompiler(options ...compiler.Option) *compiler.Compiler {
	compilerInstance, err := compiler.New(options...)
	if err != nil {
		panic(err)
	}

	return compilerInstance
}
