package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	bindingPatternAccessKind uint8

	bindingPatternAccess struct {
		Key   string
		Index int
		Kind  bindingPatternAccessKind
	}

	bindingPatternLeaf struct {
		Context     antlr.ParserRuleContext
		Declaration antlr.ParserRuleContext
		Name        string
		Path        []bindingPatternAccess
		ID          core.BindingID
	}

	destructuringBindingTarget struct {
		Destination bytecode.Operand
		Load        bytecode.Operand
		Storage     core.BindingStorage
	}
)

const (
	bindingPatternObjectAccess bindingPatternAccessKind = iota
	bindingPatternArrayAccess
)

func structuredBindingPatternLeaves(pattern fql.IStructuredBindingPatternContext) []bindingPatternLeaf {
	return appendStructuredBindingPatternLeaves(nil, pattern, nil)
}

func declarationBindingPatternLeaves(ctx fql.IVariableDeclarationContext) []bindingPatternLeaf {
	if ctx == nil || ctx.StructuredBindingPattern() == nil {
		return nil
	}

	return structuredBindingPatternLeaves(ctx.StructuredBindingPattern())
}

func appendBindingPatternLeaves(
	leaves []bindingPatternLeaf,
	pattern fql.IBindingPatternContext,
	path []bindingPatternAccess,
) []bindingPatternLeaf {
	if pattern == nil {
		return leaves
	}

	if id := pattern.BindingIdentifier(); id != nil {
		ctx := id.(antlr.ParserRuleContext)
		return append(leaves, bindingPatternLeaf{
			Context:     ctx,
			Declaration: ctx,
			Name:        textOfBindingIdentifier(id),
			Path:        append([]bindingPatternAccess(nil), path...),
			ID:          bindingIDFromRule(ctx),
		})
	}

	if structured := pattern.StructuredBindingPattern(); structured != nil {
		return appendStructuredBindingPatternLeaves(leaves, structured, path)
	}

	return leaves
}

func appendStructuredBindingPatternLeaves(
	leaves []bindingPatternLeaf,
	pattern fql.IStructuredBindingPatternContext,
	path []bindingPatternAccess,
) []bindingPatternLeaf {
	if pattern == nil {
		return leaves
	}

	if object := pattern.ObjectBindingPattern(); object != nil {
		for _, entry := range object.AllObjectBindingEntry() {
			id := entry.BindingIdentifier()
			if id == nil {
				continue
			}

			key := textOfBindingIdentifier(id)
			entryPath := appendPatternAccess(path, bindingPatternAccess{Kind: bindingPatternObjectAccess, Key: key})
			if nested := entry.BindingPattern(); nested != nil {
				leaves = appendBindingPatternLeaves(leaves, nested, entryPath)
				continue
			}

			ctx := id.(antlr.ParserRuleContext)
			leaves = append(leaves, bindingPatternLeaf{
				Context:     ctx,
				Declaration: ctx,
				Name:        key,
				Path:        entryPath,
				ID:          bindingIDFromRule(ctx),
			})
		}

		return leaves
	}

	if array := pattern.ArrayBindingPattern(); array != nil {
		for index, child := range array.AllBindingPattern() {
			childPath := appendPatternAccess(path, bindingPatternAccess{Kind: bindingPatternArrayAccess, Index: index})
			leaves = appendBindingPatternLeaves(leaves, child, childPath)
		}
	}

	return leaves
}

func appendPatternAccess(path []bindingPatternAccess, access bindingPatternAccess) []bindingPatternAccess {
	out := make([]bindingPatternAccess, len(path), len(path)+1)
	copy(out, path)

	return append(out, access)
}

func duplicateBindingPatternLeaf(leaves []bindingPatternLeaf) (bindingPatternLeaf, bindingPatternLeaf, bool) {
	seen := make(map[string]bindingPatternLeaf, len(leaves))

	for _, leaf := range leaves {
		if first, ok := seen[leaf.Name]; ok {
			return leaf, first, true
		}

		seen[leaf.Name] = leaf
	}

	return bindingPatternLeaf{}, bindingPatternLeaf{}, false
}
