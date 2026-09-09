package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func (e *Engine) compile(ctx context.Context, src source.Source, debug bool, setters []PlanOption) (*Plan, error) {
	if ctx == nil {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "context is required")
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if e.closed.Load() {
		return nil, runtime.Error(runtime.ErrInvalidOperation, "engine is closed")
	}

	level := e.optimizationLevel

	if debug {
		level = compiler.None
	}

	opts, err := newPlanConfig(level, debug, setters)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	selected := e.compiler
	if debug {
		selected = e.debugCompiler
	} else if opts.level != e.optimizationLevel {
		selected, err = compiler.New(compiler.WithOptimizationLevel(opts.level))
		if err != nil {
			return nil, err
		}
	}

	if err := e.hooks.PlanHooks.RunBeforeCompile(ctx); err != nil {
		err = fmt.Errorf("before compile hooks: %w", err)
		if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
			err = errors.Join(err, ctxErr)
		}

		return nil, err
	}

	var prog *bytecode.Program

	err = ctx.Err()
	if err == nil {
		prog, err = selected.Compile(ctx, src)
	}

	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if hookErr := e.hooks.PlanHooks.RunAfterCompile(ctx, err); hookErr != nil {
		err = errors.Join(err, fmt.Errorf("after compile hooks: %w", hookErr))
	}

	if ctxErr := ctx.Err(); ctxErr != nil && !errors.Is(err, ctxErr) {
		err = errors.Join(err, ctxErr)
	}

	if err != nil {
		return nil, err
	}

	return e.newPlan(prog)
}

func (e *Engine) newPlan(prog *bytecode.Program) (*Plan, error) {
	return &Plan{
		prog:         prog,
		closed:       make(chan struct{}),
		host:         e.host,
		hooks:        e.hooks.PlanHooks,
		sessionHooks: e.hooks.SessionHooks,
		limiter:      e.limiter,
		pool:         vm.NewPoolWithLimits(prog, e.idleCap, e.totalCap),
	}, nil
}
