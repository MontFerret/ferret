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

var (
	setInnerHTMLByCSSSelector = fmt.Sprintf(`(el, selector, value) => {
	const found = el.querySelector(selector);

	%s

	found.innerHTML = value;
}`, notFoundErrorFragment)

	setInnerHTMLByXPathSelector = fmt.Sprintf(`(el, selector, value) => {
	%s

	%s

	found.innerHTML = value;
}`, xpathAsElementFragment, notFoundErrorFragment)
)

func SetInnerHTMLBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, value values.String) *eval.Function {
	return toFunction(selector, setInnerHTMLByCSSSelector, setInnerHTMLByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArgValue(value)
}

var (
	getInnerHTMLByCSSSelector = fmt.Sprintf(`(el, selector) => {
	const found = el.querySelector(selector);

	%s

	return found.innerHTML;
}`, notFoundErrorFragment)

	getInnerHTMLByXPathSelector = fmt.Sprintf(`(el, selector) => {
	%s

	%s

	return found.innerHTML;
}`, xpathAsElementFragment, notFoundErrorFragment)
)

func GetInnerHTMLBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, getInnerHTMLByCSSSelector, getInnerHTMLByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector)
}

const getInnerHTMLByCSSSelectorAll = `(el, selector) => {
	const found = el.querySelectorAll(selector);

	return Array.from(found).map(i => i.innerHTML);
}`

var getInnerHTMLByXPathSelectorAll = fmt.Sprintf(`(el, selector) => {
	%s

	%s

	return found.map(i => i.innerHTML);
}`, xpathAsElementArrayFragment, notFoundErrorFragment)

func GetInnerHTMLBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, getInnerHTMLByCSSSelectorAll, getInnerHTMLByXPathSelectorAll).
		WithArgRef(id).
		WithArgSelector(selector)
}
