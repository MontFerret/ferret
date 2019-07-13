package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const (
	scrollTopTemplate = `
		window.scrollTo({
			left: 0,
			top: 0,
    		behavior: 'instant'
  		});
`

	scrollBottomTemplate = `
		window.scrollTo({
			left: 0,
			top: window.document.body.scrollHeight,
    		behavior: 'instant'
  		});
`

	scrollIntoViewTemplate = `
		(el) => {
			el.scrollIntoView({
				behavior: 'instant'
			});
	
			return true;
		}
`
)

func Scroll(x, y string) string {
	return fmt.Sprintf(`
			window.scrollBy({
				top: %s,
				left: %s,
				behavior: 'instant'
			});
`, x, y)
}

func ScrollTop() string {
	return scrollTopTemplate
}

func ScrollBottom() string {
	return scrollBottomTemplate
}

func ScrollIntoView() string {
	return scrollIntoViewTemplate
}

func ScrollIntoViewBySelector(selector string) string {
	return fmt.Sprintf(`
		const el = document.querySelector('%s');

		if (el == null) {
			throw new Error('%s');
		}

		el.scrollIntoView({
    		behavior: 'instant'
  		});

		return true;
`, selector, drivers.ErrNotFound)
}
