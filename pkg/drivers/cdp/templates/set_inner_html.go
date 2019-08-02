package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const setInnerHTMLTemplate = `
	(element, value) => {
		element.innerHTML = value;
	}
`

func SetInnerHTML() string {
	return setInnerHTMLTemplate
}

var setInnerHTMLBySelectorTemplate = fmt.Sprintf(`
		(el, selector, value) => {
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			found.innerHTML = value;
		}
	`,
	drivers.ErrNotFound,
)

func SetInnerHTMLBySelector() string {
	return setInnerHTMLBySelectorTemplate
}
