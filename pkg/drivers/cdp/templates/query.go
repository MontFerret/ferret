package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/mafredri/cdp/protocol/runtime"
)

const (
	queryCSSSelectorFragment = "const found = el.querySelector(selector);"

	queryCSSSelectorAllFragment = "const found = el.querySelectorAll(selector);"
)

var (
	queryCSSSelector = `
		(el, selector) => {
			const found = el.querySelector(selector);
	
			return found;
		}
	`
	queryXPathSelector = fmt.Sprintf(`
		(el, selector) => {
			%s
	
			return found;
		}
	`,
		xpathAsElementFragment,
	)
)

func QuerySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, queryCSSSelector, queryXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector)
}

const queryCSSSelectorAll = `(el, selector) => {
	return el.querySelectorAll(selector);
}`

var queryXPathSelectorAll = fmt.Sprintf(`(el, selector) => {
	%s

	return found;
}`, xpathAsElementArrayFragment)

func QuerySelectorAll(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, queryCSSSelectorAll, queryXPathSelectorAll).
		WithArgRef(id).
		WithArgSelector(selector)
}

const existsByCSSSelector = `
	(el, selector) => {
		const found = el.querySelector(selector);

		return found != null;
	}
`

var existsByXPathSelector = fmt.Sprintf(`
	(el, selector) => {
		%s

		return found != null;
	}
`, xpathAsElementFragment)

func ExistsBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, existsByCSSSelector, existsByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector)
}

const countByCSSSelector = `
	(el, selector) => {
		const found = el.querySelectorAll(selector);

		return found.length;
	}
`

var countByXPathSelector = fmt.Sprintf(`
	(el, selector) => {
		%s

		return found.length;
	}
`, xpathAsElementArrayFragment)

func CountBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	return toFunction(selector, countByCSSSelector, countByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector)
}
