package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func textOfBindingIdentifier(id fql.IBindingIdentifierContext) string {
	if id == nil {
		return ""
	}

	return id.GetText()
}

func textOfLoopVariable(loopVar fql.ILoopVariableContext) string {
	if loopVar == nil {
		return ""
	}

	return loopVar.GetText()
}

func tokenOfLoopVariable(loopVar fql.ILoopVariableContext) antlr.Token {
	if loopVar == nil {
		return nil
	}

	return loopVar.GetStart()
}
