package stdlib

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/path"
	stdlibstrings "github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/types"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/utils"
)

// Group identifies a standard library capability group.
type Group string

const (
	Types       Group = "types"
	Strings     Group = "strings"
	Math        Group = "math"
	Collections Group = "collections"
	DateTime    Group = "datetime"
	Arrays      Group = "arrays"
	Objects     Group = "objects"
	IO          Group = "io"
	FS          Group = "fs"
	NET         Group = "net"
	Path        Group = "path"
	Utils       Group = "utils"
	Testing     Group = "testing"
)

type groupRegistration struct {
	register func(runtime.Namespace)
	group    Group
}

var groupRegistrations = []groupRegistration{
	{Types, types.RegisterLib},
	{Strings, stdlibstrings.RegisterLib},
	{Math, math.RegisterLib},
	{Collections, collections.RegisterLib},
	{DateTime, datetime.RegisterLib},
	{Arrays, arrays.RegisterLib},
	{Objects, objects.RegisterLib},
	{FS, registerFS},
	{NET, registerNET},
	{Path, path.RegisterLib},
	{Utils, utils.RegisterLib},
	{Testing, testing.RegisterLib},
}

func expandGroup(group Group) ([]Group, bool) {
	switch group {
	case IO:
		return []Group{FS, NET}, true
	case Types, Strings, Math, Collections, DateTime, Arrays, Objects, FS, NET, Path, Utils, Testing:
		return []Group{group}, true
	default:
		return nil, false
	}
}

func invalidGroupsError(groups []Group) error {
	if len(groups) == 0 {
		return nil
	}

	names := make([]string, 0, len(groups))
	for _, group := range groups {
		names = append(names, string(group))
	}

	return fmt.Errorf("invalid stdlib group(s): %s", strings.Join(names, ", "))
}

func registerFS(ns runtime.Namespace) {
	fs.RegisterLib(ns.Namespace("IO"))
}

func registerNET(ns runtime.Namespace) {
	net.RegisterLib(ns.Namespace("IO"))
}
