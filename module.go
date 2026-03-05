package ferret

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	ModuleRegistry struct {
		ns         runtime.Namespace
		encoding   *encoding.Registry
		decorators []ContextDecorator
	}

	ContextDecorator func(ctx context.Context) (context.Context, error)

	Module interface {
		Name() string
		Register(registry *ModuleRegistry) error
	}
)

func (mr *ModuleRegistry) Functions() runtime.Namespace {
	return mr.ns
}

func (mr *ModuleRegistry) Encoding() *encoding.Registry {
	return mr.encoding
}

func (mr *ModuleRegistry) WithContext(extender ContextDecorator) {
	mr.decorators = append(mr.decorators, extender)
}
