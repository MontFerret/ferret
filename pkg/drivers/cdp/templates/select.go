package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const selectFragment = `
	if (el.nodeName.toLowerCase() !== 'select') {
		throw new Error('element is not a <select> element.');
	}

	const options = Array.from(el.options);

	el.value = undefined;

	for (var option of options) {
		option.selected = values.includes(option.value);
	
		if (option.selected && !el.multiple) {
			break;
		}
	}

	el.dispatchEvent(new Event('input', { 'bubbles': true }));
	el.dispatchEvent(new Event('change', { 'bubbles': true }));
	
	return options.filter(option => option.selected).map(option => option.value);
`

const selec = `(el, values) => {` + selectFragment + `}`

func Select(id runtime.RemoteObjectID, inputs *values.Array) *eval.Function {
	return eval.F(selec).WithArgRef(id).WithArgValue(inputs)
}

var selectBySelector = fmt.Sprintf(`(parent, selector, values) => {
	const el = parent.querySelector(selector);
	
	if (el == null) {
		throw new Error(%s)
	}

	%s
}`, ParamErr(core.ErrNotFound), selectFragment)

func SelectBySelector(id runtime.RemoteObjectID, selector values.String, inputs *values.Array) *eval.Function {
	return eval.F(selectBySelector).WithArgRef(id).WithArgValue(selector).WithArgValue(inputs)
}
