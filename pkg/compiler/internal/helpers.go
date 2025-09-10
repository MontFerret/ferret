package internal

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/parser/fql"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/vm"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func loadConstant(ctx *CompilerContext, value runtime.Value) vm.Operand {
	reg := ctx.Registers.Allocate(core.Temp)
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(value))
	return reg
}

func loadConstantTo(ctx *CompilerContext, constant runtime.Value, reg vm.Operand) {
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(constant))
}

//func loadIndex(ctx *CompilerContext, dst, arr vm.Operand, idx int) {
//	idxReg := loadConstant(ctx, runtime.Int(idx))
//	ctx.Emitter.EmitLoadIndex(dst, arr, idxReg)
//	ctx.Registers.Free(idxReg)
//}
//
//func loadKey(ctx *CompilerContext, dst, obj vm.Operand, key string) {
//	keyReg := loadConstant(ctx, runtime.String(key))
//	ctx.Emitter.EmitLoadKey(dst, obj, keyReg)
//}

func sortDirection(dir antlr.TerminalNode) runtime.SortDirection {
	if dir == nil {
		return runtime.SortDirectionAsc
	}

	if strings.ToLower(dir.GetText()) == "desc" {
		return runtime.SortDirectionDesc
	}

	return runtime.SortDirectionAsc
}

func getFunctionName(ctx fql.IFunctionCallContext) runtime.String {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	return runtime.NewString(strings.ToUpper(name))
}

//func copyFromNamespace(fns runtime.Builder, namespace string) error {
//	// In the name of the function "A::B::C", the namespace is "A::B",
//	// not "A::B::".
//	//
//	// So add "::" at the end.
//	namespace += "::"
//
//	// core.Builder cast every function name to upper case. Thus
//	// namespace also should be in upper case.
//	namespace = strings.ToUpper(namespace)
//
//	for _, name := range fns.Names() {
//		if !strings.HasPrefix(name, namespace) {
//			continue
//		}
//
//		noprefix := strings.Replace(name, namespace, "", 1)
//
//		if exists := fns.Has(noprefix); exists {
//			return errors.Errorf(
//				`collision occurred: "%s" already registered`,
//				noprefix,
//			)
//		}
//
//		if fn, exists := fns.F().Get(name); exists {
//			fns.F().Unset(name).Set(noprefix, fn)
//		} else if fn, exists := fns.F0().Get(name); exists {
//			fns.F0().Unset(name).Set(noprefix, fn)
//		} else if fn, exists := fns.F1().Get(name); exists {
//			fns.F1().Unset(name).Set(noprefix, fn)
//		} else if fn, exists := fns.F2().Get(name); exists {
//			fns.F2().Unset(name).Set(noprefix, fn)
//		} else if fn, exists := fns.F3().Get(name); exists {
//			fns.F3().Unset(name).Set(noprefix, fn)
//		} else if fn, exists := fns.F4().Get(name); exists {
//			fns.F4().Unset(name).Set(noprefix, fn)
//		} else {
//			return errors.Errorf(`function "%s" not found`, name)
//		}
//	}
//
//	return nil
//}
