package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const (
	waitFragment = `(el, expected, op, ...args) => {
	const actual = %s; // check

	// presence 
	if (op === 0) {
		if (actual === expected) {
			return true;
		}
	} else {
		if (actual !== expected) {
			return true;
		}
	}
	
	// null means we need to repeat
	return null;
}`

	waitBySelectorFragment = `(parent, selector, expected, op, ...args) => {
	const el = parent.querySelector(selector); // selector
	
	if (el == null) {
		return false;
	}

	const actual = %s; // check

	// presence 
	if (op === 0) {
		if (actual === expected) {
			return true;
		}
	} else {
		if (actual !== expected) {
			return true;
		}
	}
	
	// null means we need to repeat
	return null;
}`

	waitBySelectorAllFragment = `(parent, selector, expected, op, ...args) => {
	var elements = parent.querySelectorAll(selector); // selector
	
	if (elements == null || elements.length === 0) {
		return false;
	}
	
	var resultCount = 0;
	
	elements.forEach((el) => {
		var actual = %s; // check
	
		// when
		// presence 
		if (op === 0) {
			if (actual === expected) {
				resultCount++;
			}
		} else {
			if (actual !== expected) {
				resultCount++;
			}
		}
	});
	
	if (resultCount === elements.length) {
		return true;
	}
	
	// null means we need to repeat
	return null;
}`
)

func partialWait(id runtime.RemoteObjectID, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitFragment, fragment)).
		WithArgRef(id).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitAll(id runtime.RemoteObjectID, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitFragment, fragment)).
		WithArgRef(id).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitBySelector(id runtime.RemoteObjectID, selector values.String, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitBySelectorFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitBySelectorAll(id runtime.RemoteObjectID, selector values.String, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitBySelectorAllFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

const waitForElementFragment = `el.querySelector(args[0])`

func WaitForElement(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent) *eval.Function {
	return partialWait(id, values.None, when, waitForElementFragment).WithArgValue(selector)
}

const waitForElementAllFragment = `(function() {
const elements = el.querySelector(args[0]);

return elements.length;
})()`

func WaitForElementAll(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent) *eval.Function {
	return partialWait(id, values.ZeroInt, when, waitForElementAllFragment).WithArgValue(selector)
}

const waitForClassFragment = `el.className.split(' ').find(i => i === args[0]);`

func WaitForClass(id runtime.RemoteObjectID, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWait(id, values.None, when, waitForClassFragment).WithArgValue(class)
}

func WaitForClassBySelector(id runtime.RemoteObjectID, selector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelector(id, selector, values.None, when, waitForClassFragment).WithArgValue(class)
}

func WaitForClassBySelectorAll(id runtime.RemoteObjectID, selector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelectorAll(id, selector, values.None, when, waitForClassFragment).WithArgValue(class)
}

const waitForAttributeFragment = `el.attributes[name] != null ? el.attributes[name].value : null`

func WaitForAttribute(id runtime.RemoteObjectID, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWait(id, expected, when, waitForAttributeFragment).WithArgValue(name)
}

func WaitForAttributeBySelector(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelector(id, selector, expected, when, waitForAttributeFragment).WithArgValue(name)
}

func WaitForAttributeBySelectorAll(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelectorAll(id, selector, expected, when, waitForAttributeFragment).WithArgValue(name)
}

const waitForStyleFragment = `(function getStyles() {
	const styles = window.getComputedStyle(el);
	return styles[args[0]];
})()`

func WaitForStyle(id runtime.RemoteObjectID, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWait(id, expected, when, waitForStyleFragment).WithArgValue(name)
}

func WaitForStyleBySelector(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelector(id, selector, expected, when, waitForStyleFragment).WithArgValue(name)
}

func WaitForStyleBySelectorAll(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitBySelectorAll(id, selector, expected, when, waitForStyleFragment).WithArgValue(name)
}
