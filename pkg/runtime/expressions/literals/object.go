package literals

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	ObjectPropertyAssignment struct {
		name  core.Expression
		value core.Expression
	}

	ObjectLiteral struct {
		properties []*ObjectPropertyAssignment
	}
)

func NewObjectPropertyAssignment(name, value core.Expression) (*ObjectPropertyAssignment, error) {
	if name == nil {
		return nil, core.Error(core.ErrMissedArgument, "property name expression")
	}

	if value == nil {
		return nil, core.Error(core.ErrMissedArgument, "property value expression")
	}

	return &ObjectPropertyAssignment{name, value}, nil
}

func NewObjectLiteral(props []*ObjectPropertyAssignment) *ObjectLiteral {
	return &ObjectLiteral{props}
}

func NewObjectLiteralWith(props ...*ObjectPropertyAssignment) *ObjectLiteral {
	return NewObjectLiteral(props)
}

func (l *ObjectLiteral) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	obj := values.NewObject()

	for _, el := range l.properties {
		name, err := el.name.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		val, err := el.value.Exec(ctx, scope)

		if err != nil {
			return values.None, err
		}

		if name.Type() != types.String {
			return values.None, core.TypeError(name.Type(), types.String)
		}

		obj.Set(name.(values.String), val)
	}

	return obj, nil
}
