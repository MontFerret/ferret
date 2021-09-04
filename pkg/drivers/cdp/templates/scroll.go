package templates

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

const (
	isElementInViewportFragment = `(i) => {
	var bounding = i.getBoundingClientRect();
	
	return (
		bounding.top >= 0 &&
		bounding.left >= 0 &&
		bounding.bottom <= (window.innerHeight || document.documentElement.clientHeight) &&
		bounding.right <= (window.innerWidth || document.documentElement.clientWidth)
	);
};`

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

	scrollIntoViewBySelector = fmt.Sprintf(`(parent, selector, opts) => {
		const el = parent.querySelector(selector);

		if (el == null) {
			throw new Error(%s);
		}

		%s

		if (!isInViewport(el)) {
			el.scrollIntoView({
				behavior: opts.behavior,
				block: opts.block, 
				inline: opts.inline
  			});
		}

		return true;
}`, ParamErr(core.ErrNotFound), isElementInViewportFragment)
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

func ScrollIntoViewBySelector(id runtime.RemoteObjectID, selector values.String, options drivers.ScrollOptions) *eval.Function {
	return eval.F(scrollIntoViewBySelector).WithArgRef(id).WithArgValue(selector).WithArg(options)
}
