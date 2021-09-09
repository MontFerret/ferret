package drivers

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Polyfill for Go 1.10
	SameSite int

	// HTTPCookie HTTPCookie object
	HTTPCookie struct {
		Name     string
		Value    string
		Path     string
		Domain   string
		Expires  time.Time
		MaxAge   int
		Secure   bool
		HTTPOnly bool
		SameSite SameSite
	}
)

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
)

func (s SameSite) String() string {
	switch s {
	case SameSiteLaxMode:
		return "Lax"
	case SameSiteStrictMode:
		return "Strict"
	default:
		return ""
	}
}

func (c HTTPCookie) Type() core.Type {
	return HTTPCookieType
}

func (c HTTPCookie) String() string {
	return fmt.Sprintf("%s=%s", c.Name, c.Value)
}

func (c HTTPCookie) Compare(other core.Value) int64 {
	if other.Type() != HTTPCookieType {
		return Compare(HTTPCookieType, other.Type())
	}

	oc := other.(HTTPCookie)

	if c.Name != oc.Name {
		return int64(strings.Compare(c.Name, oc.Name))
	}

	if c.Value != oc.Value {
		return int64(strings.Compare(c.Value, oc.Value))
	}

	if c.Path != oc.Path {
		return int64(strings.Compare(c.Path, oc.Path))
	}

	if c.Domain != oc.Domain {
		return int64(strings.Compare(c.Domain, oc.Domain))
	}

	if c.Expires.After(oc.Expires) {
		return 1
	} else if c.Expires.Before(oc.Expires) {
		return -1
	}

	if c.MaxAge > oc.MaxAge {
		return 1
	} else if c.MaxAge < oc.MaxAge {
		return -1
	}

	if c.Secure && !oc.Secure {
		return 1
	} else if !c.Secure && oc.Secure {
		return -1
	}

	if c.HTTPOnly && !oc.HTTPOnly {
		return 1
	} else if !c.HTTPOnly && oc.HTTPOnly {
		return -1
	}

	if c.SameSite > oc.SameSite {
		return 1
	} else if c.SameSite < oc.SameSite {
		return -1
	}

	return 0
}

func (c HTTPCookie) Unwrap() interface{} {
	return c.Value
}

func (c HTTPCookie) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(c.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(c.Name))
	h.Write([]byte(c.Value))
	h.Write([]byte(c.Path))
	h.Write([]byte(c.Domain))
	h.Write([]byte(c.Expires.String()))
	h.Write([]byte(strconv.Itoa(c.MaxAge)))
	h.Write([]byte(fmt.Sprintf("%t", c.Secure)))
	h.Write([]byte(fmt.Sprintf("%t", c.HTTPOnly)))
	h.Write([]byte(c.SameSite.String()))

	return h.Sum64()
}

func (c HTTPCookie) Copy() core.Value {
	return *(&c)
}

func (c HTTPCookie) MarshalJSON() ([]byte, error) {
	v := map[string]interface{}{
		"name":      c.Name,
		"value":     c.Value,
		"path":      c.Path,
		"domain":    c.Domain,
		"expires":   c.Expires,
		"max_age":   c.MaxAge,
		"secure":    c.Secure,
		"http_only": c.HTTPOnly,
		"same_site": c.SameSite.String(),
	}

	out, err := jettison.MarshalOpts(v, jettison.NoHTMLEscaping())

	if err != nil {
		return nil, err
	}

	return out, err
}

func (c HTTPCookie) GetIn(_ context.Context, path []core.Value) (core.Value, core.PathError) {
	if len(path) == 0 {
		return values.None, nil
	}

	segment := path[0]

	switch values.ToString(segment) {
	case "name":
		return values.NewString(c.Name), nil
	case "value":
		return values.NewString(c.Value), nil
	case "path":
		return values.NewString(c.Path), nil
	case "domain":
		return values.NewString(c.Domain), nil
	case "expires":
		return values.NewDateTime(c.Expires), nil
	case "maxAge":
		return values.NewInt(c.MaxAge), nil
	case "secure":
		return values.NewBoolean(c.Secure), nil
	case "httpOnly":
		return values.NewBoolean(c.HTTPOnly), nil
	case "sameSite":
		return values.NewString(c.SameSite.String()), nil
	default:
		return values.None, nil
	}
}
