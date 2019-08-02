package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const setInnerTextTemplate = `
	(element, value) => {
		element.innerText = value;
	}
`

func SetInnerText() string {
	return setInnerTextTemplate
}

var setInnerTextBySelectorTemplate = fmt.Sprintf(`
		(el, selector, value) => {
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			found.innerText = value;
		}
	`,
	drivers.ErrNotFound,
)

func SetInnerTextBySelector() string {
	return setInnerTextBySelectorTemplate
}
