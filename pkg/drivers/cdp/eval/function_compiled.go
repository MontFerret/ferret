package eval

import "github.com/mafredri/cdp/protocol/runtime"

type CompiledFunction struct {
	id  runtime.ScriptID
	src *Function
}

func CF(id runtime.ScriptID, src *Function) *CompiledFunction {
	op := new(CompiledFunction)
	op.id = id
	op.src = src

	return op
}

func (fn *CompiledFunction) returnNothing() *CompiledFunction {
	fn.src.returnNothing()

	return fn
}

func (fn *CompiledFunction) returnRef() *CompiledFunction {
	fn.src.returnRef()

	return fn
}

func (fn *CompiledFunction) returnValue() *CompiledFunction {
	fn.src.returnValue()

	return fn
}

func (fn *CompiledFunction) call(ctx runtime.ExecutionContextID) *runtime.RunScriptArgs {
	call := runtime.NewRunScriptArgs(fn.id).
		SetAwaitPromise(fn.src.async).
		SetReturnByValue(fn.src.returnType == ReturnValue)

	if ctx != EmptyExecutionContextID {
		call.SetExecutionContextID(ctx)
	}

	return call
}
