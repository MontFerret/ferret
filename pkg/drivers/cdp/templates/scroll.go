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

	scrollTemplate = `
		window.scrollTo({
			left: %s,
			top: %s,
			behavior: '%s',
			block: '%s', 
			inline: '%s'
  		});
`

	scrollTopTemplate = `
		window.scrollTo({
			left: 0,
			top: 0,
			behavior: '%s',
			block: '%s', 
			inline: '%s'
  		});
`

	scrollBottomTemplate = `
		window.scrollTo({
			left: 0,
			top: window.document.body.scrollHeight,
			behavior: '%s',
			block: '%s', 
			inline: '%s'
  		});
`

	scrollIntoViewTemplate = `
		(el) => {
			` + isElementInViewportTemplate + `

			if (!isInViewport(el)) {
				el.scrollIntoView({
					behavior: '%s',
					block: '%s', 
					inline: '%s'
				});
			}
	
			return true;
		}
`

	scrollIntoViewBySelectorTemplate = `
		const el = document.querySelector('%s');

		if (el == null) {
			throw new Error('%s');
		}

		` + isElementInViewportTemplate + `

		if (!isInViewport(el)) {
			el.scrollIntoView({
				behavior: '%s',
				block: '%s', 
				inline: '%s'
  			});
		}

		return true;
`
)

func Scroll(x, y string, options drivers.ScrollOptions) string {
	return fmt.Sprintf(
		scrollTemplate,
		x,
		y,
		options.Behavior,
		options.Block,
		options.Inline,
	)
}

func ScrollTop(options drivers.ScrollOptions) string {
	return fmt.Sprintf(
		scrollTopTemplate,
		options.Behavior,
		options.Block,
		options.Inline,
	)
}

func ScrollBottom(options drivers.ScrollOptions) string {
	return fmt.Sprintf(
		scrollBottomTemplate,
		options.Behavior,
		options.Block,
		options.Inline,
	)
}

func ScrollIntoView(options drivers.ScrollOptions) string {
	return fmt.Sprintf(
		scrollIntoViewTemplate,
		options.Behavior,
		options.Block,
		options.Inline,
	)
}

func ScrollIntoViewBySelector(selector string, options drivers.ScrollOptions) string {
	return fmt.Sprintf(
		scrollIntoViewBySelectorTemplate,
		selector,
		drivers.ErrNotFound,
		options.Behavior,
		options.Block,
		options.Inline,
	)
}
