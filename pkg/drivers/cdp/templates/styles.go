package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const getStyles = `(el) => {
	const out = {};
	const styles = window.getComputedStyle(el);

	Object.keys(styles).forEach((key) => {
		if (!isNaN(parseFloat(key))) {
			const name = styles[key];
			const value = styles.getPropertyValue(name);
			out[name] = value;
		}
	});

	return out;
}`

func GetStyles(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getStyles).WithArgRef(id)
}

const getStyle = `(el, name) => {
	const styles = window.getComputedStyle(el);

	return styles[name];
}`

func GetStyle(id runtime.RemoteObjectID, name values.String) *eval.Function {
	return eval.F(getStyle).WithArgRef(id).WithArgValue(name)
}

const setStyle = `(el, name, value) => {
	el.style[name] = value;
}`

func SetStyle(id runtime.RemoteObjectID, name, value values.String) *eval.Function {
	return eval.F(setStyle).WithArgRef(id).WithArgValue(name).WithArgValue(value)
}

const setStyles = `(el, values) => {
	Object.keys(values).forEach((key) => {
		el.style[key] = values[key]
	});
}`

func SetStyles(id runtime.RemoteObjectID, values *values.Object) *eval.Function {
	return eval.F(setStyles).WithArgRef(id).WithArgValue(values)
}

const removeStyles = `(el, names) => {
	const style = el.style;
	names.forEach((name) => { style[name] = "" })
}`

func RemoveStyles(id runtime.RemoteObjectID, names []values.String) *eval.Function {
	return eval.F(removeStyles).WithArgRef(id).WithArg(names)
}

const removeStylesAll = `(el) => {
	el.removeAttribute("style");
}`

func RemoveStylesAll(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(removeStylesAll).WithArgRef(id)
}
