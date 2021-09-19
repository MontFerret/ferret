package eval

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/wI2L/jettison"
	"strconv"
	"strings"
)

type Function struct {
	exp        string
	ownerID    runtime.RemoteObjectID
	args       FunctionArguments
	returnType ReturnType
	async      bool
}

const defaultArgsCount = 5

func F(exp string) *Function {
	op := new(Function)
	op.exp = exp
	op.returnType = ReturnNothing

	return op
}

func (fn *Function) CallOn(id runtime.RemoteObjectID) *Function {
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

func (fn *Function) WithArgSelector(selector drivers.QuerySelector) *Function {
	return fn.WithArg(selector.String())
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

func (fn *Function) Length() int {
	return len(fn.args)
}

func (fn *Function) returnNothing() *Function {
	fn.returnType = ReturnNothing

	return fn
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
		fn.args = make([]runtime.CallArgument, 0, defaultArgsCount)
	}

	fn.args = append(fn.args, arg)

	return fn
}

func (fn *Function) call(ctx runtime.ExecutionContextID) *runtime.CallFunctionOnArgs {
	exp := fn.prepExp()

	call := runtime.NewCallFunctionOnArgs(exp).
		SetAwaitPromise(fn.async).
		SetReturnByValue(fn.returnType == ReturnValue)

	if fn.ownerID != EmptyObjectID {
		call.SetObjectID(fn.ownerID)
	} else if ctx != EmptyExecutionContextID {
		call.SetExecutionContextID(ctx)
	}

	if len(fn.args) > 0 {
		call.SetArguments(fn.args)
	}

	return call
}

func (fn *Function) compile(ctx runtime.ExecutionContextID) *runtime.CompileScriptArgs {
	exp := fn.precompileExp()

	call := runtime.NewCompileScriptArgs(exp, "", true)

	if ctx != EmptyExecutionContextID {
		call.SetExecutionContextID(ctx)
	}

	return call
}

func (fn *Function) prepExp() string {
	exp := strings.TrimSpace(fn.exp)

	if strings.HasPrefix(exp, "(") || strings.HasPrefix(exp, "function") {
		return exp
	}

	args := len(fn.args)

	if args == 0 {
		return "() => {\n" + exp + "\n}"
	}

	var buf strings.Builder
	lastIndex := args - 1

	for i := 0; i < args; i++ {
		buf.WriteString("arg")
		buf.WriteString(strconv.Itoa(i + 1))

		if i != lastIndex {
			buf.WriteString(",")
		}
	}

	return "(" + buf.String() + ") => {\n" + exp + "\n}"
}

func (fn *Function) precompileExp() string {
	exp := fn.prepExp()
	args := fn.args

	if len(args) == 0 {
		return exp
	}

	var buf strings.Builder
	var l = len(args)

	buf.WriteString("(function () {\n")
	buf.WriteString("const args = [")

	for i := 0; i < l; i++ {
		buf.WriteRune('\n')

		arg := args[i]

		if arg.Value != nil {
			buf.Write(arg.Value)
		} else if arg.ObjectID != nil {
			buf.WriteString("(() => { throw new Error('Reference values cannot be used in pre-compiled scrips')})()")
		}

		buf.WriteString(",")
	}
	buf.WriteRune('\n')
	buf.WriteString("];")

	buf.WriteRune('\n')
	buf.WriteString("const exp = ")
	buf.WriteString(exp)
	buf.WriteRune('\n')

	buf.WriteString("return exp(...args);")

	buf.WriteString("})")

	return buf.String()
}
