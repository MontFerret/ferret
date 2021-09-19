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
	name       string
	ownerID    runtime.RemoteObjectID
	args       FunctionArguments
	returnType ReturnType
	async      bool
}

const (
	defaultArgsCount = 5
	expName          = "$exp"
	compiledExpName  = "$$exp"
)

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

func (fn *Function) AsAnonymous() *Function {
	fn.name = ""

	return fn
}

func (fn *Function) AsNamed(name string) *Function {
	if name != "" {
		fn.name = name
	}

	return fn
}

func (fn *Function) WithArgRemoteValue(value RemoteValue) *Function {
	return fn.WithArgRef(value.RemoteID())
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

func (fn *Function) eval(ctx runtime.ExecutionContextID) *runtime.CallFunctionOnArgs {
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
	var invoke bool
	exp := strings.TrimSpace(fn.exp)
	name := fn.name

	// If the given expression is either an arrow or plain function
	if strings.HasPrefix(exp, "(") || strings.HasPrefix(exp, "function") {
		// And if this function must be an anonymous
		// we just pass the expression as is without wrapping it.
		if name == "" {
			return exp
		}

		// But if the function must be identified (named)
		// we need to wrap the given function with a named one/
		// And then eval it with passing available arguments.
		invoke = true
	}

	// Start building a wrapper
	var buf strings.Builder
	buf.WriteString("function")

	// Name the function if the name is set
	if name != "" {
		buf.WriteString(" ")
		buf.WriteString(name)
	}

	buf.WriteString("(")

	// If the given expression is a function then we do not need to define wrapper's function arguments.
	// Any available arguments will be passed down via 'arguments' runtime variable.
	// Otherwise, we define a list of arguments as argN, so the given expression could access them by name.
	if !invoke {
		args := len(fn.args)
		lastIndex := args - 1

		for i := 0; i < args; i++ {
			buf.WriteString("arg")
			buf.WriteString(strconv.Itoa(i + 1))

			if i != lastIndex {
				buf.WriteString(",")
			}
		}
	}

	buf.WriteString(") {\n")

	if !invoke {
		buf.WriteString(exp)
	} else {
		buf.WriteString("const ")
		buf.WriteString(expName)
		buf.WriteString(" = ")
		buf.WriteString(exp)
		buf.WriteString(";\n")
		buf.WriteString("return ")
		buf.WriteString(expName)
		buf.WriteString(".apply(this, arguments);")
	}

	buf.WriteString("\n}")

	return buf.String()
}

func (fn *Function) precompileExp() string {
	exp := fn.prepExp()
	args := fn.args

	var buf strings.Builder
	var l = len(args)

	buf.WriteString("const args = [")

	if l > 0 {
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
	}

	buf.WriteString("];")

	buf.WriteRune('\n')
	buf.WriteString("const ")
	buf.WriteString(compiledExpName)
	buf.WriteString(" = ")
	buf.WriteString(exp)
	buf.WriteString(";\n")
	buf.WriteString(compiledExpName)
	buf.WriteString(".apply(this, args);")
	buf.WriteString("\n")

	return buf.String()
}
