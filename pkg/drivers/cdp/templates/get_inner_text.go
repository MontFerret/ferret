package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

func GetInnerTextBySelector(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelector(selector);

		if (found == null) {
			throw new Error('%s');
		}

		return found.innerText;
	`, selector, drivers.ErrNotFound)
}

func GetInnerTextBySelectorAll(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelectorAll(selector);

		if (found == null) {
			throw new Error('%s');
		}

		return Array.from(found).map(i => i.innerText);
	`, selector, drivers.ErrNotFound)
}
