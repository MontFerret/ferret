package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/mafredri/cdp/protocol/runtime"
)

const getPreviousElementSibling = "(el) => el.previousElementSibling"
const getNextElementSibling = "(el) => el.nextElementSibling"

func GetPreviousElementSibling(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getPreviousElementSibling).WithArgRef(id)
}

func GetNextElementSibling(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getNextElementSibling).WithArgRef(id)
}
