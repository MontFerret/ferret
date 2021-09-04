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
	waitExistenceFragment = `(el, op, ...args) => {
	const actual = %s; // check

	// presence 
	if (op === 0) {
		if (actual != null) {
			return true;
		}
	} else {
		if (actual == null) {
			return true;
		}
	}
	
	// null means we need to repeat
	return null;
}`

	waitEqualityFragment = `(el, expected, op, ...args) => {
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

	waitExistenceBySelectorFragment = `(parent, selector, op, ...args) => {
	const el = parent.querySelector(selector); // selector
	
	if (el == null) {
		return false;
	}

	const actual = %s; // check

	// presence 
	if (op === 0) {
		if (actual != null) {
			return true;
		}
	} else {
		if (actual == null) {
			return true;
		}
	}
	
	// null means we need to repeat
	return null;
}`

	waitEqualityBySelectorFragment = `(parent, selector, expected, op, ...args) => {
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

	waitExistenceBySelectorAllFragment = `(parent, selector, op, ...args) => {
	const elements = parent.querySelectorAll(selector); // selector
	
	if (elements == null || elements.length === 0) {
		return false;
	}
	
	let resultCount = 0;
	
	elements.forEach((el) => {
		let actual = %s; // check
	
		// when
		// presence 
		if (op === 0) {
			if (actual != null) {
				resultCount++;
			}
		} else {
			if (actual == null) {
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

	waitEqualityBySelectorAllFragment = `(parent, selector, expected, op, ...args) => {
	const elements = parent.querySelectorAll(selector); // selector
	
	if (elements == null || elements.length === 0) {
		return false;
	}
	
	let resultCount = 0;

	elements.forEach((el) => {
		let actual = %s; // check
	
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

func partialWaitExistence(id runtime.RemoteObjectID, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitExistenceFragment, fragment)).
		WithArgRef(id).
		WithArg(int(when))
}

func partialWaitEquality(id runtime.RemoteObjectID, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitEqualityFragment, fragment)).
		WithArgRef(id).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitExistenceBySelector(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitExistenceBySelectorFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArg(int(when))
}

func partialWaitEqualityBySelector(id runtime.RemoteObjectID, selector values.String, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitEqualityBySelectorFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitExistenceBySelectorAll(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitExistenceBySelectorAllFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArg(int(when))
}

func partialWaitEqualityBySelectorAll(id runtime.RemoteObjectID, selector values.String, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	return eval.F(fmt.Sprintf(waitEqualityBySelectorAllFragment, fragment)).
		WithArgRef(id).
		WithArgValue(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

const waitForElementFragment = `el.querySelector(args[0])`

func WaitForElement(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistence(id, when, waitForElementFragment).WithArgValue(selector)
}

const waitForElementAllFragment = `(function() {
const elements = el.querySelector(args[0]);

return elements.length;
})()`

func WaitForElementAll(id runtime.RemoteObjectID, selector values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEquality(id, values.ZeroInt, when, waitForElementAllFragment).WithArgValue(selector)
}

const waitForClassFragment = `el.className.split(' ').find(i => i === args[0]);`

func WaitForClass(id runtime.RemoteObjectID, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistence(id, when, waitForClassFragment).WithArgValue(class)
}

func WaitForClassBySelector(id runtime.RemoteObjectID, selector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistenceBySelector(id, selector, when, waitForClassFragment).WithArgValue(class)
}

func WaitForClassBySelectorAll(id runtime.RemoteObjectID, selector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistenceBySelectorAll(id, selector, when, waitForClassFragment).WithArgValue(class)
}

const waitForAttributeFragment = `el.getAttribute(args[0])`

func WaitForAttribute(id runtime.RemoteObjectID, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEquality(id, expected, when, waitForAttributeFragment).WithArgValue(name)
}

func WaitForAttributeBySelector(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelector(id, selector, expected, when, waitForAttributeFragment).WithArgValue(name)
}

func WaitForAttributeBySelectorAll(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelectorAll(id, selector, expected, when, waitForAttributeFragment).WithArgValue(name)
}

const waitForStyleFragment = `(function getStyles() {
	const styles = window.getComputedStyle(el);
	return styles[args[0]];
})()`

func WaitForStyle(id runtime.RemoteObjectID, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEquality(id, expected, when, waitForStyleFragment).WithArgValue(name)
}

func WaitForStyleBySelector(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelector(id, selector, expected, when, waitForStyleFragment).WithArgValue(name)
}

func WaitForStyleBySelectorAll(id runtime.RemoteObjectID, selector, name, expected values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelectorAll(id, selector, expected, when, waitForStyleFragment).WithArgValue(name)
}
