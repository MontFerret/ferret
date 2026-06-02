package internal

import "github.com/antlr4-go/antlr/v4"

func reportVariableNotFound(ctx *CompilationSession, token antlr.Token, name string) {
	if ctx == nil || ctx.Program == nil || ctx.Program.Errors == nil {
		return
	}

	if index := ctx.Program.ForwardBindings; index != nil {
		if declaration, ok := index.Lookup(token, name); ok {
			ctx.Program.Errors.VariableUsedBeforeDeclaration(token, name, declaration)
			return
		}
	}

	ctx.Program.Errors.VariableNotFound(token, name)
}
