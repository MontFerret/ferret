package internal

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

// ProgramContext holds state that lives for the entire compilation of one
// source file. It is shared across the main body and all UDF compilations.
//
// Everything in ProgramContext is program-wide: writes from any function
// compile (top-level body or a nested UDF body) are observable globally and
// persist for the whole compilation. The final bytecode.Program reads its
// metadata almost entirely from here.
type ProgramContext struct {
	Emitter             *core.Emitter
	Constants           *core.ConstantPool
	CatchTable          *core.CatchStack
	UDFs                *core.UDFTable
	HostParams          *core.HostParamTable
	HostFunctions       *core.HostFunctionTable
	Errors              *diagnostics.ErrorHandler
	Source              *source.Source
	UseAliases          map[string]string
	aggregatePlanByHash map[uint64][]int
	aggregatePlans      []*bytecode.AggregatePlan
	OptimizationLevel   optimization.Level
}

// FunctionContext holds state that is local to a single function body
// compilation. A fresh FunctionContext is created for the top-level body
// and for every UDF body.
//
// Function-local state is NEVER inherited from an enclosing function: a
// nested UDF starts with an empty RegisterAllocator, an empty SymbolTable
// (locals and globals), a fresh TypeTracker, a fresh LoopTable, and its
// own UDFScope pointing at the UDF's body scope for inner UDF name
// lookup. Carrying any of these over from a parent function would produce
// wrong bytecode.
//
// Anything that must be visible across function boundaries (host params,
// host function refs, constants, catch table, UDF metadata, the emitter)
// lives on ProgramContext instead.
type FunctionContext struct {
	Registers *core.RegisterAllocator
	Symbols   *core.SymbolTable
	Types     *core.TypeTracker
	Loops     *core.LoopTable
	UDFScope  *core.UDFScope
}

// CompilationSession is the thin coordinator passed to all compilers.
// It provides access to both program-wide and function-local state.
// The Function pointer is swapped by withFunctionCompileState during UDF
// compilation; the Program pointer never changes during a single Compile.
type CompilationSession struct {
	Program  *ProgramContext
	Function *FunctionContext
}

// NewFunctionContext creates a fresh function-local compilation state.
func NewFunctionContext(constants *core.ConstantPool) *FunctionContext {
	fc := &FunctionContext{
		Registers: core.NewRegisterAllocator(),
		Types:     core.NewTypeTracker(),
	}
	fc.Symbols = core.NewSymbolTable(fc.Registers, constants)
	fc.Loops = core.NewLoopTable(fc.Registers)

	return fc
}

// NewCompilationSession initializes a new CompilationSession with default values.
func NewCompilationSession(src *source.Source, errors *diagnostics.ErrorHandler, level optimization.Level) *CompilationSession {
	program := &ProgramContext{
		Source:            src,
		Errors:            errors,
		OptimizationLevel: level,

		Emitter:       core.NewEmitter(),
		Constants:     core.NewConstantPool(),
		CatchTable:    core.NewCatchStack(),
		HostParams:    core.NewHostParamTable(),
		HostFunctions: core.NewHostFunctionTable(),

		UseAliases: make(map[string]string),

		aggregatePlans:      make([]*bytecode.AggregatePlan, 0),
		aggregatePlanByHash: make(map[uint64][]int),
	}

	return &CompilationSession{
		Program:  program,
		Function: NewFunctionContext(program.Constants),
	}
}

func (c *ProgramContext) AddAggregatePlan(plan *bytecode.AggregatePlan) int {
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

func (c *ProgramContext) AggregatePlans() []bytecode.AggregatePlan {
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

	if len(a.Keys) != len(b.Keys) || a.TrackGroupValues != b.TrackGroupValues {
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

	if plan.TrackGroupValues {
		h.Write([]byte{1})
	} else {
		h.Write([]byte{0})
	}

	for i, key := range plan.Keys {
		h.Write([]byte(key))
		h.Write([]byte{0})

		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(plan.Kinds[i]))
		h.Write(buf[:])
	}

	return h.Sum64()
}
