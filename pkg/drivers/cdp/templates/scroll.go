package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
)

const (
	isElementInViewportTemplate = `
		function isInViewport(elem) {
			var bounding = elem.getBoundingClientRect();
			
			return (
				bounding.top >= 0 &&
				bounding.left >= 0 &&
				bounding.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
				bounding.right <= (window.innerWidth || document.documentElement.clientWidth)
			);
		};
	`

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
			` + isElementInViewportTemplate + `
			

			if (!isInViewport(el)) {
				el.scrollIntoView({
					block: 'center', 
					inline: 'center',
					behavior: 'instant'
				});
			}
	
			return true;
		}
`
)

func Scroll(x, y string) string {
	return fmt.Sprintf(`
			window.scrollBy(%s, %s);
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

		%s

		if (!isInViewport(el)) {
			el.scrollIntoView({
				block: 'center', 
				inline: 'center',
    			behavior: 'instant'
  			});
		}

		return true;
`, selector, drivers.ErrNotFound, isElementInViewportTemplate)
}
