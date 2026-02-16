package internal

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

// CompilerContext holds the context for the compilation process, including various compilers and allocators.
type CompilerContext struct {
	Source     *file.Source
	Emitter    *core.Emitter
	Registers  *core.RegisterAllocator
	Symbols    *core.SymbolTable
	Types      *core.TypeTracker
	Loops      *core.LoopTable
	CatchTable *core.CatchStack
	Errors     *diagnostics.ErrorHandler
	UseAliases map[string]string

	aggregatePlans      []*bytecode.AggregatePlan
	aggregatePlanByHash map[uint64][]int

	ExprCompiler        *ExprCompiler
	LiteralCompiler     *LiteralCompiler
	StmtCompiler        *StmtCompiler
	LoopCompiler        *LoopCompiler
	LoopSortCompiler    *LoopSortCompiler
	LoopCollectCompiler *LoopCollectCompiler
	WaitCompiler        *WaitCompiler
	DispatchCompiler    *DispatchCompiler
}

// NewCompilerContext initializes a new CompilerContext with default values.
func NewCompilerContext(src *file.Source, errors *diagnostics.ErrorHandler) *CompilerContext {
	ctx := &CompilerContext{
		Source:              src,
		Errors:              errors,
		Emitter:             core.NewEmitter(),
		Registers:           core.NewRegisterAllocator(),
		Symbols:             nil, // set later
		Loops:               nil, // set later
		CatchTable:          core.NewCatchStack(),
		UseAliases:          make(map[string]string),
		aggregatePlans:      make([]*bytecode.AggregatePlan, 0),
		aggregatePlanByHash: make(map[uint64][]int),
	}

	ctx.Symbols = core.NewSymbolTable(ctx.Registers)
	ctx.Types = core.NewTypeTracker()
	ctx.Loops = core.NewLoopTable(ctx.Registers)

	ctx.ExprCompiler = NewExprCompiler(ctx)
	ctx.LiteralCompiler = NewLiteralCompiler(ctx)
	ctx.StmtCompiler = NewStmtCompiler(ctx)
	ctx.LoopCompiler = NewLoopCompiler(ctx)
	ctx.LoopSortCompiler = NewLoopSortCompiler(ctx)
	ctx.LoopCollectCompiler = NewCollectCompiler(ctx)
	ctx.WaitCompiler = NewWaitCompiler(ctx)
	ctx.DispatchCompiler = NewDispatchCompiler(ctx)

	return ctx
}

func (c *CompilerContext) AddAggregatePlan(plan *bytecode.AggregatePlan) int {
	if plan == nil {
		return -1
	}

	hash := aggregatePlanHash(plan)

	if existing, ok := c.aggregatePlanByHash[hash]; ok {
		for _, idx := range existing {
			if idx >= 0 && idx < len(c.aggregatePlans) && areAggregatePlansEqual(c.aggregatePlans[idx], plan) {
				return idx
			}
		}
	}

	idx := len(c.aggregatePlans)
	c.aggregatePlans = append(c.aggregatePlans, plan)
	c.aggregatePlanByHash[hash] = append(c.aggregatePlanByHash[hash], idx)

	return idx
}

func (c *CompilerContext) AggregatePlans() []bytecode.AggregatePlan {
	plans := make([]bytecode.AggregatePlan, len(c.aggregatePlans))

	for i, p := range c.aggregatePlans {
		if p != nil {
			plans[i] = *p
		}
	}

	return plans
}

func areAggregatePlansEqual(a, b *bytecode.AggregatePlan) bool {
	if len(a.Keys) == 0 || len(b.Keys) == 0 {
		return a == b
	}

	if len(a.Keys) != len(b.Keys) {
		return false
	}

	for i := 0; i < len(a.Keys); i++ {
		if a.Keys[i] != b.Keys[i] || a.Kinds[i] != b.Kinds[i] {
			return false
		}
	}

	return true
}

func aggregatePlanHash(plan *bytecode.AggregatePlan) uint64 {
	h := fnv.New64a()
	h.Write([]byte("aggregate_plan:"))

	for i, key := range plan.Keys {
		h.Write([]byte(key))
		h.Write([]byte{0})

		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(plan.Kinds[i]))
		h.Write(buf[:])
	}

	return h.Sum64()
}
