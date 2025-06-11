package internal

import (
	"github.com/MontFerret/ferret/pkg/parser/fql"
)

type CollectCompiler struct {
	Ctx *FuncContext
}

func NewCollectCompiler(ctx *FuncContext) *CollectCompiler {
	return &CollectCompiler{Ctx: ctx}
}

func (cc *CollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	//loop := cc.Ctx.Loops.Current()
	//if loop == nil {
	//	panic("COLLECT clause must appear inside a loop")
	//}
	//
	//// Grouping by key
	//if group := ctx.CollectGrouping(); group != nil {
	//	// Example: COLLECT key = expr
	//	keyName := group.Variable().GetText()
	//	keyExpr := group.Expression()
	//	keyReg := cc.Ctx.ExprCompiler.Compile(keyExpr)
	//
	//	loop.Result = cc.Ctx.Registers.Allocate(Result)
	//
	//	cc.Ctx.Emitter.EmitABC(vm.OpCollect, loop.Result, keyReg, keyReg) // src1=key, src2=key (single-group)
	//	cc.Ctx.Symbols.DeclareLocal(keyName)
	//}
	//
	//// Aggregation
	//if agg := ctx.CollectAggregator(); agg != nil {
	//	for _, part := range agg.AllCollectGroupVariable() {
	//		name := part.Variable().GetText()
	//		expr := part.Expression()
	//
	//		src := cc.Ctx.ExprCompiler.Compile(expr)
	//		dst := cc.Ctx.Registers.Allocate(Result)
	//
	//		cc.Ctx.Emitter.EmitABC(vm.OpCollect, dst, src, src)
	//		cc.Ctx.Symbols.DeclareLocal(name)
	//	}
	//}
	//
	//// Optional counter
	//if counter := ctx.CollectCounter(); counter != nil {
	//	name := counter.Variable().GetText()
	//	dst := cc.Ctx.Registers.Allocate(Result)
	//
	//	cc.Ctx.Emitter.EmitAB(vm.OpCount, dst, loop.Value)
	//	cc.Ctx.Symbols.DeclareLocal(name)
	//}
}
