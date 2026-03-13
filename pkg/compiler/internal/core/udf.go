package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	UDFTable struct {
		GlobalScope *UDFScope
		Functions   []*UDFInfo
	}

	UDFScope struct {
		Parent    *UDFScope
		Functions map[string]*UDFInfo
	}

	UDFInfo struct {
		Decl        fql.IFunctionDeclarationContext
		Scope       *UDFScope
		BodyScope   *UDFScope
		Name        string
		DisplayName string
		Params      []string
		Captures    []string
		ID          int
		Entry       int
		Registers   int
	}
)

func NewUDFScope(parent *UDFScope) *UDFScope {
	return &UDFScope{
		Parent:    parent,
		Functions: make(map[string]*UDFInfo),
	}
}

func NewUDFTable() *UDFTable {
	return &UDFTable{
		Functions: make([]*UDFInfo, 0),
	}
}

func (t *UDFTable) Metadata() []bytecode.UDF {
	if t == nil || len(t.Functions) == 0 {
		return nil
	}

	out := make([]bytecode.UDF, 0, len(t.Functions))
	for _, fn := range t.Functions {
		if fn == nil {
			continue
		}

		out = append(out, bytecode.UDF{
			Name:        fn.Name,
			DisplayName: fn.DisplayName,
			Entry:       fn.Entry,
			Registers:   fn.Registers,
			Params:      len(fn.Params) + len(fn.Captures),
		})
	}

	return out
}

func (t *UDFTable) Resolve(name string, scope *UDFScope) (*UDFInfo, bool) {
	if scope == nil {
		return nil, false
	}

	for s := scope; s != nil; s = s.Parent {
		if fn, ok := s.Functions[name]; ok {
			return fn, true
		}
	}

	return nil, false
}
