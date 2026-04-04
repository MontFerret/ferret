package internal

import (
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
