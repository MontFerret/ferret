package network

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/mafredri/cdp/protocol/page"
	"regexp"
)

type (
	Cookies map[string]*drivers.HTTPCookies

	Filter struct {
		Patterns []drivers.ResourceFilter
	}

	Options struct {
		Cookies Cookies
		Headers *drivers.HTTPHeaders
		Filter  *Filter
	}

	WaitEventOptions struct {
		FrameID page.FrameID
		URL     *regexp.Regexp
	}
)
