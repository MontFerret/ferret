package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const blur = `(el) => {
	el.blur()
}`

func Blur(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(blur).WithArgRef(id)
}

var blurBySelector = fmt.Sprintf(`
		(el, selector) => {
			const found = el.querySelector(selector);

			if (found == null) {
				throw new Error(%s)
			}

			found.blur();
		}
`, ParamErr(drivers.ErrNotFound))

func BlurBySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(blurBySelector).WithArgRef(id).WithArgValue(selector)
}
