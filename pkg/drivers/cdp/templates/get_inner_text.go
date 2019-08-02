package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

var getInnerTextBySelectorTemplate = fmt.Sprintf(`
	(el, selector) => {
		const found = el.querySelector(selector);

		if (found == null) {
			throw new Error("%s");
		}

		return found.innerText;
	}
	`, drivers.ErrNotFound,
)

func GetInnerTextBySelector() string {
	return getInnerTextBySelectorTemplate
}

var getInnerTextBySelectorAllTemplate = fmt.Sprintf(`
	(el, selector) => {
		const found = el.querySelectorAll(selector);

		if (found == null) {
			throw new Error("%s");
		}

		return Array.from(found).map(i => i.innerText);
	}
	`, drivers.ErrNotFound,
)

func GetInnerTextBySelectorAll() string {
	return getInnerTextBySelectorAllTemplate
}
