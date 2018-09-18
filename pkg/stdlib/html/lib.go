package html

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"DOCUMENT":       Document,
		"DOCUMENT_PARSE": DocumentParse,
		"ELEMENT":        Element,
		"ELEMENTS":       Elements,
	}
}
