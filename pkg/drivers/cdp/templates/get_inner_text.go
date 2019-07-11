package templates

import "fmt"

func GetInnerTextBySelector(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelector(selector);

		if (found == null) {
			throw new Error('Element not found by selector: ' + selector);
		}

		return found.innerText;
	`, selector)
}

func GetInnerTextBySelectorAll(selector string) string {
	return fmt.Sprintf(`
		const selector = "%s";
		const found = document.querySelectorAll(selector);

		if (found == null) {
			throw new Error('Elements not found by selector: ' + selector);
		}

		return Array.from(found).map(i => i.innerText);
	`, selector)
}
