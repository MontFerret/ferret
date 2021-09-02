package eval

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/runtime"
	"strings"
)

type (
	FunctionReturnType int

	Function struct {
		exp        string
		ownerID    *runtime.RemoteObjectID
		args       []runtime.CallArgument
		returnType FunctionReturnType
		async      bool
	}

	FunctionOption func(op *Function)
)

const (
	ReturnNothing FunctionReturnType = iota
	ReturnValue
	ReturnRef
)

func newFunction(exp string, opts []FunctionOption) *Function {
	op := new(Function)
	op.exp = exp
	op.returnType = ReturnNothing

	for _, opt := range opts {
		opt(op)
	}

	return op
}

func (fn *Function) Use(opt FunctionOption) {
	opt(fn)
}

func (fn *Function) toArgs(ctx runtime.ExecutionContextID) *runtime.CallFunctionOnArgs {
	exp := strings.TrimSpace(fn.exp)

	if !strings.HasPrefix(exp, "(") && !strings.HasPrefix(exp, "function") {
		exp = wrapExp(exp)
	}

	call := runtime.NewCallFunctionOnArgs(exp).
		SetAwaitPromise(fn.async)

	if fn.returnType == ReturnValue {
		call.SetReturnByValue(true)
	}

	if ctx != EmptyExecutionContextID {
		call.SetExecutionContextID(ctx)
	}

	if fn.ownerID != nil {
		call.SetObjectID(*fn.ownerID)
	}

	if len(fn.args) > 0 {
		call.SetArguments(fn.args)
	}

	return call
}

func withReturnRef() FunctionOption {
	return func(op *Function) {
		op.returnType = ReturnRef
	}
}

func withReturnValue() FunctionOption {
	return func(op *Function) {
		op.returnType = ReturnValue
	}
}

func WithArgs(args ...runtime.CallArgument) FunctionOption {
	return func(op *Function) {
		if op.args == nil {
			op.args = args
		} else {
			op.args = append(op.args, args...)
		}
	}
}

func WithArgValue(value core.Value) FunctionOption {
	raw, err := value.MarshalJSON()

	if err != nil {
		// we defer error
		return WithArgs(runtime.CallArgument{
			Value: []byte(err.Error()),
		})
	}

	return WithArgs(runtime.CallArgument{
		Value: raw,
	})
}

func WithArgRef(id runtime.RemoteObjectID) FunctionOption {
	return WithArgs(runtime.CallArgument{
		ObjectID: &id,
	})
}

func WithOwner(ctx *runtime.RemoteObjectID) FunctionOption {
	return func(op *Function) {
		op.ownerID = ctx
	}
}

func WithAsync() FunctionOption {
	return func(op *Function) {
		op.async = true
	}
}
