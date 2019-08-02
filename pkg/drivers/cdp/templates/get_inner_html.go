package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

func GetInnerHTMLBySelector(selector string) string {
	return fmt.Sprintf(`
		(el) => {
			const selector = '%s';
			const found = el.querySelector(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			return found.innerHTML;
		}
	`, selector, drivers.ErrNotFound)
}

func GetInnerHTMLBySelectorAll(selector string) string {
	return fmt.Sprintf(`
		(el) => {
			const selector = '%s';
			const found = el.querySelectorAll(selector);
	
			if (found == null) {
				throw new Error('%s');
			}
	
			return Array.from(found).map(i => i.innerHTML);
		}
	`, selector, drivers.ErrNotFound)
}
