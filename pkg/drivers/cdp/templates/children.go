package templates

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

const getChildren = "(el) => Array.from(el.children)"

func GetChildren(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getChildren).WithArgRef(id)
}

const getChildrenCount = "(el) => el.children.length"

func GetChildrenCount(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getChildrenCount).WithArgRef(id)
}

const getChildByIndex = "(el, idx) => el.children[idx]"

func GetChildByIndex(id runtime.RemoteObjectID, index core.Int) *eval.Function {
	return eval.F(getChildByIndex).WithArgRef(id).WithArgValue(index)
}
