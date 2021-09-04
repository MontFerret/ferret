package eval

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/wI2L/jettison"
	"strings"
)

type (
	FunctionReturnType int

	Function struct {
		exp        string
		ownerID    runtime.RemoteObjectID
		args       []runtime.CallArgument
		returnType FunctionReturnType
		async      bool
	}
)

const (
	ReturnNothing FunctionReturnType = iota
	ReturnValue
	ReturnRef
)

func F(exp string) *Function {
	op := new(Function)
	op.exp = exp
	op.returnType = ReturnNothing

	return op
}

func (fn *Function) AsPartOf(id runtime.RemoteObjectID) *Function {
	fn.ownerID = id

	return fn
}

func (fn *Function) AsAsync() *Function {
	fn.async = true

	return fn
}

func (fn *Function) AsSync() *Function {
	fn.async = false

	return fn
}

func (fn *Function) WithArgRef(id runtime.RemoteObjectID) *Function {
	return fn.withArg(runtime.CallArgument{
		ObjectID: &id,
	})
}

func (fn *Function) WithArgValue(value core.Value) *Function {
	raw, err := value.MarshalJSON()

	if err != nil {
		panic(err)
	}

	return fn.withArg(runtime.CallArgument{
		Value: raw,
	})
}

func (fn *Function) WithArg(value interface{}) *Function {
	raw, err := jettison.MarshalOpts(value, jettison.NoHTMLEscaping())

	if err != nil {
		panic(err)
	}

	return fn.withArg(runtime.CallArgument{
		Value: raw,
	})
}

func (fn *Function) String() string {
	return fn.exp
}

func (fn *Function) returnRef() *Function {
	fn.returnType = ReturnRef

	return fn
}

func (fn *Function) returnValue() *Function {
	fn.returnType = ReturnValue

	return fn
}

func (fn *Function) withArg(arg runtime.CallArgument) *Function {
	if fn.args == nil {
		fn.args = make([]runtime.CallArgument, 0, 5)
	}

	fn.args = append(fn.args, arg)

	return fn
}

func (fn *Function) build(ctx runtime.ExecutionContextID) *runtime.CallFunctionOnArgs {
	exp := strings.TrimSpace(fn.exp)

	if !strings.HasPrefix(exp, "(") && !strings.HasPrefix(exp, "function") {
		exp = wrapExp(exp, len(fn.args))
	}

	call := runtime.NewCallFunctionOnArgs(exp).
		SetAwaitPromise(fn.async)

	if fn.returnType == ReturnValue {
		call.SetReturnByValue(true)
	}

	if ctx != EmptyExecutionContextID {
		call.SetExecutionContextID(ctx)
	}

	if fn.ownerID != "" {
		call.SetObjectID(fn.ownerID)
	}

	if len(fn.args) > 0 {
		call.SetArguments(fn.args)
	}

	return call
}
