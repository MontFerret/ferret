package internal

import "github.com/antlr4-go/antlr/v4"

type captureBindingInfo struct {
	Decl    antlr.ParserRuleContext
	Name    string
	Mutable bool
}
