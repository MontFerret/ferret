package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

type captureBindingInfo struct {
	Decl    antlr.ParserRuleContext
	Name    string
	ID      core.BindingID
	Mutable bool
}
