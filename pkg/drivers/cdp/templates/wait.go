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

	waitExistenceBySelectorFragment = `(el, selector, op, ...args) => {
	// selector
	%s

	if (found == null) {
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

	waitEqualityBySelectorFragment = `(el, selector, expected, op, ...args) => {
	// selector
	%s

	if (found == null) {
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

	waitExistenceBySelectorAllFragment = `(el, selector, op, ...args) => {
	// selector
	%s
	
	if (found == null || !found || found.length === 0) {
		return false;
	}
	
	let resultCount = 0;
	
	found.forEach((el) => {
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
	
	if (resultCount === found.length) {
		return true;
	}
	
	// null means we need to repeat
	return null;
}`

	waitEqualityBySelectorAllFragment = `(el, selector, expected, op, ...args) => {
	// selector
	%s
	
	if (found == null || !found || found.length === 0) {
		return false;
	}
	
	let resultCount = 0;

	found.forEach((el) => {
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
	
	if (resultCount === found.length) {
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

func partialWaitExistenceBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, when drivers.WaitEvent, fragment string) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = fmt.Sprintf(waitExistenceBySelectorFragment, queryCSSSelectorFragment, fragment)
	} else {
		tmpl = fmt.Sprintf(waitExistenceBySelectorFragment, xpathAsElementFragment, fragment)
	}

	return eval.F(tmpl).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArg(int(when))
}

func partialWaitEqualityBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = fmt.Sprintf(waitEqualityBySelectorFragment, queryCSSSelectorFragment, fragment)
	} else {
		tmpl = fmt.Sprintf(waitEqualityBySelectorFragment, xpathAsElementFragment, fragment)
	}

	return eval.F(tmpl).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

func partialWaitExistenceBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, when drivers.WaitEvent, fragment string) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = fmt.Sprintf(waitExistenceBySelectorAllFragment, queryCSSSelectorAllFragment, fragment)
	} else {
		tmpl = fmt.Sprintf(waitExistenceBySelectorAllFragment, xpathAsElementArrayFragment, fragment)
	}

	return eval.F(tmpl).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArg(int(when))
}

func partialWaitEqualityBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, expected core.Value, when drivers.WaitEvent, fragment string) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = fmt.Sprintf(waitEqualityBySelectorAllFragment, queryCSSSelectorAllFragment, fragment)
	} else {
		tmpl = fmt.Sprintf(waitEqualityBySelectorAllFragment, xpathAsElementArrayFragment, fragment)
	}

	return eval.F(tmpl).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArgValue(expected).
		WithArg(int(when))
}

const waitForElementByCSSFragment = `el.querySelector(args[0])`

var waitForElementByXPathFragment = fmt.Sprintf(`(() => {
const selector = args[0];

%s

return found;
})()`, xpathAsElementFragment)

func WaitForElement(id runtime.RemoteObjectID, selector drivers.QuerySelector, when drivers.WaitEvent) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = waitForElementByCSSFragment
	} else {
		tmpl = waitForElementByXPathFragment
	}

	return partialWaitExistence(id, when, tmpl).WithArgSelector(selector)
}

const waitForElementAllByCSSFragment = `(function() {
const elements = el.querySelector(args[0]);

return elements.length;
})()`

var waitForElementAllByXPathFragment = fmt.Sprintf(`(function() {
const selector = args[0];

%s

return found;
})()`, xpathAsElementArrayFragment)

func WaitForElementAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, when drivers.WaitEvent) *eval.Function {
	var tmpl string

	if selector.Kind() == drivers.CSSSelector {
		tmpl = waitForElementAllByCSSFragment
	} else {
		tmpl = waitForElementAllByXPathFragment
	}

	return partialWaitEquality(id, values.ZeroInt, when, tmpl).WithArgSelector(selector)
}

const waitForClassFragment = `el.className.split(' ').find(i => i === args[0]);`

func WaitForClass(id runtime.RemoteObjectID, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistence(id, when, waitForClassFragment).WithArgValue(class)
}

const waitForClassBySelectorFragment = `found.className.split(' ').find(i => i === args[0]);`

func WaitForClassBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistenceBySelector(id, selector, when, waitForClassBySelectorFragment).WithArgValue(class)
}

func WaitForClassBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) *eval.Function {
	return partialWaitExistenceBySelectorAll(id, selector, when, waitForClassFragment).WithArgValue(class)
}

const waitForAttributeFragment = `el.getAttribute(args[0])`

func WaitForAttribute(id runtime.RemoteObjectID, name values.String, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEquality(id, expected, when, waitForAttributeFragment).WithArgValue(name)
}

const waitForAttributeBySelectorFragment = `found.getAttribute(args[0])`

func WaitForAttributeBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, name core.Value, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelector(id, selector, expected, when, waitForAttributeBySelectorFragment).WithArgValue(name)
}

func WaitForAttributeBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, name values.String, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelectorAll(id, selector, expected, when, waitForAttributeFragment).WithArgValue(name)
}

const waitForStyleFragment = `(function getStyles() {
	const styles = window.getComputedStyle(el);
	return styles[args[0]];
})()`

func WaitForStyle(id runtime.RemoteObjectID, name values.String, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEquality(id, expected, when, waitForStyleFragment).WithArgValue(name)
}

const waitForStyleBySelectorFragment = `(function getStyles() {
	const styles = window.getComputedStyle(found);
	return styles[args[0]];
})()`

func WaitForStyleBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, name values.String, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelector(id, selector, expected, when, waitForStyleBySelectorFragment).WithArgValue(name)
}

func WaitForStyleBySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector, name values.String, expected core.Value, when drivers.WaitEvent) *eval.Function {
	return partialWaitEqualityBySelectorAll(id, selector, expected, when, waitForStyleFragment).WithArgValue(name)
}
