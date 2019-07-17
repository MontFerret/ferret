package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

func selectBase(values string) string {
	return fmt.Sprintf(`
		const values = %s;

		if (el.nodeName.toLowerCase() !== 'select') {
			throw new Error('element is not a <select> element.');
		}

		const options = Array.from(el.options);

		el.value = undefined;

		for (var option of options) {
			option.selected = values.includes(option.value);
		
			if (option.selected && !el.multiple) {
				break;
			}
		}

		el.dispatchEvent(new Event('input', { 'bubbles': true }));
		el.dispatchEvent(new Event('change', { 'bubbles': true }));
		
		return options.filter(option => option.selected).map(option => option.value);
	`, values,
	)
}

func Select(values string) string {
	return fmt.Sprintf(`
		(el) => {
			%s
		}
	`, selectBase(values),
	)
}

func SelectBySelector(selector, values string) string {
	return fmt.Sprintf(`
		const el = document.querySelector('%s');
		
		if (el == null) {
			throw new Error("%s")
		}

		%s
	`, selector, drivers.ErrNotFound, selectBase(values),
	)
}
