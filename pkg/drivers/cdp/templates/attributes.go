package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const getAttribute = `(el, name) => {
	return el.getAttribute(name)
}`

func GetAttribute(id runtime.RemoteObjectID, name values.String) *eval.Function {
	return eval.F(getAttribute).WithArgRef(id).WithArgValue(name)
}

const getAttributes = `(el) => {
	return el.getAttributeNames().reduce((res, name) => {
		const out = res;
		out[name] = el.getAttribute(name);
		return out;
	}, {});
}`

func GetAttributes(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(getAttributes).WithArgRef(id)
}

const setAttribute = `(el, name, value) => {
	el.setAttribute(name, value)
}`

func SetAttribute(id runtime.RemoteObjectID, name, value values.String) *eval.Function {
	return eval.F(setAttribute).WithArgRef(id).WithArgValue(name).WithArgValue(value)
}

const setAttributes = `(el, values) => {
	Object.keys(values).forEach((name) => {
		const value = values[name];
		el.setAttribute(name, value)
	});
}`

func SetAttributes(id runtime.RemoteObjectID, values *values.Object) *eval.Function {
	return eval.F(setAttributes).WithArgRef(id).WithArgValue(values)
}

const removeAttribute = `(el, name) => {
	el.removeAttribute(name)
}`

func RemoveAttribute(id runtime.RemoteObjectID, name values.String) *eval.Function {
	return eval.F(removeAttribute).WithArgRef(id).WithArgValue(name)
}

const removeAttributes = `(el, names) => {
	names.forEach(name => el.removeAttribute(name));
}`

func RemoveAttributes(id runtime.RemoteObjectID, names []values.String) *eval.Function {
	return eval.F(removeAttributes).WithArgRef(id).WithArg(names)
}
