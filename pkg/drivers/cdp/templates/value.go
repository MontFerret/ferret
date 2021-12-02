package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/runtime"
)

const getValue = `(el) => {
	return el.value
}`

func GetValue(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getValue).WithArgRef(id)
}

const setValue = `(el, value) => {
	el.value = value
}`

func SetValue(id runtime.RemoteObjectID, value core.Value) *eval.Function {
	return eval.F(setValue).WithArgRef(id).WithArgValue(value)
}
