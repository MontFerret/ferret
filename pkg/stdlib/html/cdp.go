package html

import "context"

func fromContext(ctx context.Context) string {
	str, ok := ctx.Value("cdp").(string)

	if ok {
		return str
	}

	return ""
}
