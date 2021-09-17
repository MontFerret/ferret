package templates

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/mafredri/cdp/protocol/runtime"
)

const blur = `(el) => {
	el.blur()
}`

func Blur(id runtime.RemoteObjectID) *eval.Function {
	return eval.F(blur).WithArgRef(id)
}

var (
	blurByCSSSelector = fmt.Sprintf(`
		(el, selector) => {
			const found = el.querySelector(selector);

			%s

			found.blur();
		}
`, notFoundErrorFragment)

	blurByXPathSelector = fmt.Sprintf(`
		(el, selector) => {
			%s

			%s

			found.blur();
		}
`, xpathAsElementFragment, notFoundErrorFragment)
)

func BlurBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector) *eval.Function {
	var f *eval.Function

	if selector.Kind() == drivers.CSSSelector {
		f = eval.F(blurByCSSSelector)
	} else {
		f = eval.F(blurByXPathSelector)
	}

	return f.WithArgRef(id).WithArgSelector(selector)
}
