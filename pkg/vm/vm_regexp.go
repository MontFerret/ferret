package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func (exec *execState) regexpCached(pc int, value runtime.Value) (*data.Regexp, error) {
	// We compare patterns to ensure that the cached regexp is the same as the one we're trying to use.
	// This is necessary because the same compiled function can be used in different places with different regexps,
	// and we want to avoid caching a regexp that doesn't match the current pattern.
	switch v := value.(type) {
	case *data.Regexp:
		pattern := v.String()

		if cached := exec.vm.cache.Regexps[pc]; cached == nil || cached.Pattern != pattern {
			exec.vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: v}
		}

		return v, nil
	case runtime.String:
		pattern := v.String()

		if cached := exec.vm.cache.Regexps[pc]; cached != nil && cached.Pattern == pattern {
			return cached.Regexp, nil
		}

		r, err := data.NewRegexp(v)
		if err != nil {
			return nil, err
		}

		exec.vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: r}

		return r, nil
	default:
		return nil, runtime.TypeErrorOf(value, runtime.TypeString, data.TypeRegexp)
	}
}
