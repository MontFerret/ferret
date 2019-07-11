package templates

import "fmt"

func GetInnerHTMLBySelector(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelector(selector);

		if (found == null) {
			throw new Error('Element not found by selector: ' + selector);
		}

		return found.innerHTML;
	`, selector)
}

func GetInnerHTMLBySelectorAll(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelectorAll(selector);

		if (found == null) {
			throw new Error('Elements not found by selector: ' + selector);
		}

		return Array.from(found).map(i => i.innerHTML);
	`, selector)
}
