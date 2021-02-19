package network

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/mafredri/cdp/protocol/fetch"
)

type (
	Cookies map[string]drivers.HTTPCookies

	Filter struct {
		Patterns []drivers.ResourceFilter
	}

	Options struct {
		Cookies Cookies
		Headers drivers.HTTPHeaders
		Filter  Filter
	}
)

func toFetchArgs(filterPatterns []drivers.ResourceFilter) *fetch.EnableArgs {
	patterns := make([]fetch.RequestPattern, 0, len(filterPatterns))

	for _, pattern := range filterPatterns {
		rt := toResourceType(pattern.Type)

		patterns = append(patterns, fetch.RequestPattern{
			URLPattern:   &pattern.URL,
			ResourceType: &rt,
		})
	}

	return &fetch.EnableArgs{
		Patterns: patterns,
	}
}
