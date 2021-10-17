package network

import (
	"github.com/MontFerret/ferret/pkg/drivers"
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

	EventOptions struct {
		FrameID string
		URL     *regexp.Regexp
	}
)
