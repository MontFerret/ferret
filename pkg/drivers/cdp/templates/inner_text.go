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
	return eval.F(getInnerText).WithArgRef(id)
}

var (
	setInnerTextByCSSSelector = fmt.Sprintf(`
(el, selector, value) => {
	const found = el.querySelector(selector);

	%s

	found.innerText = value;
}`, notFoundErrorFragment)

	setInnerTextByXPathSelector = fmt.Sprintf(`
(el, selector, value) => {
	%s

	%s

	found.innerText = value;
}`, xpathAsElementFragment, notFoundErrorFragment)
)

func SetInnerTextBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, value values.String) *eval.Function {
	return toFunction(selector, setInnerTextByCSSSelector, setInnerTextByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArgValue(value)
}

var (
	getInnerTextByCSSSelector = fmt.Sprintf(`
(el, selector) => {
	const found = el.querySelector(selector);

	%s

	return found.innerText;
}`, notFoundErrorFragment)

	getInnerTextByXPathSelector = fmt.Sprintf(`
(el, selector) => {
	%s

	%s

	return found.innerText;
}`, xpathAsElementFragment, notFoundErrorFragment)
)

func GetInnerTextBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, getInnerTextByCSSSelector, getInnerTextByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector)
}

var (
	getInnerTextByCSSSelectorAll = fmt.Sprintf(`
(el, selector) => {
	const found = el.querySelectorAll(selector);

	%s

	return Array.from(found).map(i => i.innerText);
}`, notFoundErrorFragment)

	getInnerTextByXPathSelectorAll = fmt.Sprintf(`
(el, selector) => {
	%s

	%s

	return found.map(i => i.innerText);
}`, xpathAsElementArrayFragment, notFoundErrorFragment)
)

func GetInnerTextBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, getInnerTextByCSSSelectorAll, getInnerTextByXPathSelectorAll).
		WithArgRef(id).
		WithArgSelector(selector)
}
