package templates

import (
	"fmt"

	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

const (
	isElementInViewportFragment = `function isInViewport(i) {
	var bounding = i.getBoundingClientRect();
	
	return (
		bounding.top >= 0 &&
		bounding.left >= 0 &&
		bounding.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
		bounding.right <= (window.innerWidth || document.documentElement.clientWidth)
	);
}`

	scroll = `(opts) =>
	window.scrollTo({
		left: opts.left,
		top: opts.top,
		behavior: opts.behavior,
		block: opts.block, 
		inline: opts.inline
	});
}`

	scrollTop = `(opts) => {
	window.scrollTo({
		left: 0,
		top: 0,
		behavior: opts.behavior,
		block: opts.block, 
		inline: opts.inline
	});
}`

	scrollBottom = `(opts) => {
	window.scrollTo({
		left: 0,
		top: window.document.body.scrollHeight,
		behavior: opts.behavior,
		block: opts.block, 
		inline: opts.inline
	});
}`
)

var (
	scrollIntoView = fmt.Sprintf(`(el, opts) => {
	%s

	if (!isInViewport(el)) {
		el.scrollIntoView({
			behavior: opts.behavior,
			block: opts.block, 
			inline: opts.inline
		});
	}

	return true;
}`, isElementInViewportFragment)

	scrollIntoViewByCSSSelector = fmt.Sprintf(`(el, selector, opts) => {
		const found = el.querySelector(selector);

		%s

		%s

		if (!isInViewport(found)) {
			found.scrollIntoView({
				behavior: opts.behavior,
				block: opts.block, 
				inline: opts.inline
  			});
		}

		return true;
}`, notFoundErrorFragment, isElementInViewportFragment)

	scrollIntoViewByXPathSelector = fmt.Sprintf(`(el, selector, opts) => {
		%s

		%s

		%s

		if (!isInViewport(found)) {
			found.scrollIntoView({
				behavior: opts.behavior,
				block: opts.block, 
				inline: opts.inline
  			});
		}

		return true;
}`, xpathAsElementFragment, notFoundErrorFragment, isElementInViewportFragment)
)

func Scroll(options drivers.ScrollOptions) *eval.Function {
	return eval.F(scroll).WithArg(options)
}

func ScrollTop(options drivers.ScrollOptions) *eval.Function {
	return eval.F(scrollTop).WithArg(options)
}

func ScrollBottom(options drivers.ScrollOptions) *eval.Function {
	return eval.F(scrollBottom).WithArg(options)
}

func ScrollIntoView(id runtime.RemoteObjectID, options drivers.ScrollOptions) *eval.Function {
	return eval.F(scrollIntoView).WithArgRef(id).WithArg(options)
}

func ScrollIntoViewBySelector(id runtime.RemoteObjectID, selector drivers.QuerySelector, options drivers.ScrollOptions) *eval.Function {
	return toFunction(selector, scrollIntoViewByCSSSelector, scrollIntoViewByXPathSelector).
		WithArgRef(id).
		WithArgSelector(selector).
		WithArg(options)
}
