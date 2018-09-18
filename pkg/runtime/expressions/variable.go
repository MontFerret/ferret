package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/pkg/errors"
)

type (
	VariableExpression struct {
		src  core.SourceMap
		name string
	}

	VariableDeclarationExpression struct {
		*VariableExpression
		init core.Expression
	}
)

func NewVariableExpression(src core.SourceMap, name string) (*VariableExpression, error) {
	if name == "" {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed variable name")
	}

	return &VariableExpression{src, name}, nil
}

func NewVariableDeclarationExpression(src core.SourceMap, name string, init core.Expression) (*VariableDeclarationExpression, error) {
	v, err := NewVariableExpression(src, name)

	if err != nil {
		return nil, err
	}

	if core.IsNil(init) {
		return nil, errors.Wrap(core.ErrMissedArgument, "missed variable initializer")
	}

	return &VariableDeclarationExpression{v, init}, nil
}

func (e *VariableExpression) Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error) {
	value, err := e.Exec(ctx, scope)

	if err != nil {
		return nil, core.SourceError(e.src, err)
	}

	iter, err := collections.ToIterator(value)

	if err != nil {
		return nil, core.SourceError(e.src, err)
	}

	return iter, nil
}

func (e *VariableExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	return scope.GetVariable(e.name)
}

func (e *VariableDeclarationExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	val, err := e.init.Exec(ctx, scope)

	if err != nil {
		return values.None, err
	}

	return values.None, scope.SetVariable(e.name, val)
}
