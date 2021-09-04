package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

var querySelector = fmt.Sprintf(`
		(el, selector) => {
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error(%s);
			}
	
			return found;
		}
	`,
	ParamErr(drivers.ErrNotFound),
)

func QuerySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(querySelector).WithArgRef(id).WithArgValue(selector)
}

const querySelectorAll = `(el, selector) => return el.querySelectorAll(selector);`

func QuerySelectorAll(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(querySelectorAll).WithArgRef(id).WithArgValue(selector)
}

const existsBySelector = `
	(el, selector) => {
		const found = el.querySelector(selector);

		return found != null;
	}
`

func ExistsBySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(existsBySelector).WithArgRef(id).WithArgValue(selector)
}

const countBySelector = `
	(el, selector) => {
		const found = el.querySelectorAll(selector);

		return found.length;
	}
`

func CountBySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(countBySelector).WithArgRef(id).WithArgValue(selector)
}
