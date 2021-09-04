package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const setInnerText = `(el, value) => {
	el.innerText = value;
}`

func SetInnerText(id runtime.RemoteObjectID, value values.String) *eval.Function {
	return eval.F(setInnerText).WithArgRef(id).WithArgValue(value)
}

const getInnerText = `(el) => {
	if (el.nodeType !== 9) {
		return el.innerText;
	}

	return document.documentElement.innerText;
}`

func GetInnerText(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(setInnerText).WithArgRef(id)
}

var setInnerTextBySelector = fmt.Sprintf(`
(el, selector, value) => {
	const found = el.querySelector(selector);

	if (found == null) {
		throw new Error(%s);
	}

	found.innerText = value;
}`, ParamErr(drivers.ErrNotFound))

func SetInnerTextBySelector(id runtime.RemoteObjectID, selector, value values.String) *eval.Function {
	return eval.F(setInnerTextBySelector).WithArgRef(id).WithArgValue(selector).WithArgValue(value)
}

var getInnerTextBySelector = fmt.Sprintf(`
(el, selector) => {
	const found = el.querySelector(selector);

	if (found == null) {
		throw new Error(%s);
	}

	return found.innerText;
}`, ParamErr(drivers.ErrNotFound))

func GetInnerTextBySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(getInnerTextBySelector).WithArgRef(id).WithArgValue(selector)
}

var getInnerTextBySelectorAll = fmt.Sprintf(`
(el, selector) => {
	const found = el.querySelectorAll(selector);

	if (found == null) {
		throw new Error(%s);
	}

	return Array.from(found).map(i => i.innerText);
}`, ParamErr(drivers.ErrNotFound))

func GetInnerTextBySelectorAll(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(getInnerTextBySelectorAll).WithArgRef(id).WithArgValue(selector)
}
