package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const setInnerHTML = `(el, value) => {
	el.innerHTML = value;
}`

func SetInnerHTML(id runtime.RemoteObjectID, value values.String) *eval.Function {
	return eval.F(setInnerHTML).WithArgRef(id).WithArgValue(value)
}

const getInnerHTML = `(el) => {
	if (el.nodeType !== 9) {
		return el.innerHTML;
	}

	return document.documentElement.innerHTML;
}`

func GetInnerHTML(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getInnerHTML).WithArgRef(id)
}

var setInnerHTMLBySelector = fmt.Sprintf(`(el, selector, value) => {
	const found = el.querySelector(selector);

	if (found == null) {
		throw new Error(%s);
	}

	found.innerHTML = value;
}`, ParamErr(drivers.ErrNotFound))

func SetInnerHTMLBySelector(id runtime.RemoteObjectID, selector, value values.String) *eval.Function {
	return eval.F(setInnerHTMLBySelector).WithArgRef(id).WithArgValue(selector).WithArgValue(value)
}

var getInnerHTMLBySelector = fmt.Sprintf(`(el, selector) => {
	const found = el.querySelector(selector);

	if (found == null) {
		throw new Error(%s);
	}

	return found.innerHTML;
}`, ParamErr(drivers.ErrNotFound))

func GetInnerHTMLBySelector(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(getInnerHTMLBySelector).WithArgRef(id).WithArgValue(selector)
}

const getInnerHTMLBySelectorAll = `(el, selector) => {
	const found = el.querySelectorAll(selector);

	return Array.from(found).map(i => i.innerHTML);
}`

func GetInnerHTMLBySelectorAll(id runtime.RemoteObjectID, selector values.String) *eval.Function {
	return eval.F(getInnerHTMLBySelectorAll).WithArgRef(id).WithArgValue(selector)
}
