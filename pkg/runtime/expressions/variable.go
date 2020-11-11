package expressions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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
		return nil, core.Error(core.ErrMissedArgument, "variable name")
	}

	return &VariableExpression{src, name}, nil
}

func NewVariableDeclarationExpression(src core.SourceMap, name string, init core.Expression) (*VariableDeclarationExpression, error) {
	v, err := NewVariableExpression(src, name)

	if err != nil {
		return nil, err
	}

	if init == nil {
		return nil, core.Error(core.ErrMissedArgument, "missed variable initializer")
	}

	return &VariableDeclarationExpression{v, init}, nil
}

func (e *VariableExpression) Exec(_ context.Context, scope *core.Scope) (core.Value, error) {
	return scope.GetVariable(e.name)
}

func (e *VariableDeclarationExpression) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	val, err := e.init.Exec(ctx, scope)

	if err != nil {
		return values.None, err
	}

	return values.None, scope.SetVariable(e.name, val)
}
