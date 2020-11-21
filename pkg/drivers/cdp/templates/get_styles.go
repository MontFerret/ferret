package templates

var getStylesTemplate = `
	(el) => {
		const out = {};
		const styles = window.getComputedStyle(el);
	
		Object.keys(styles).forEach((key) => {
			if (!isNaN(parseFloat(key))) {
				const name = styles[key];
				const value = styles.getPropertyValue(name);
				out[name] = value;
			}
		});

		return out;
	}
`

func GetStyles() string {
	return getStylesTemplate
}
