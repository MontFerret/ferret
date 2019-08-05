package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

var getInnerHTMLBySelectorTemplate = fmt.Sprintf(`
		(el, selector) => {
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			return found.innerHTML;
		}
	`, drivers.ErrNotFound,
)

func GetInnerHTMLBySelector() string {
	return getInnerHTMLBySelectorTemplate
}

var getInnerHTMLBySelectorAllTemplate = fmt.Sprintf(`
		(el, selector) => {
			const found = el.querySelectorAll(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			return Array.from(found).map(i => i.innerHTML);
		}
	`, drivers.ErrNotFound,
)

func GetInnerHTMLBySelectorAll() string {
	return getInnerHTMLBySelectorAllTemplate
}
