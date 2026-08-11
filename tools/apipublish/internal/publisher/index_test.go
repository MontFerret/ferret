package publisher

import (
	"reflect"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
)

func TestNextIndexUsesDescendingRawVersionTieBreaker(t *testing.T) {
	existing := &api.Index{
		SchemaVersion: api.IndexSchemaVersion,
		Latest:        "2.0.0+build.1",
		Versions: []api.IndexVersion{{
			Version: "2.0.0+build.1",
			Href:    canonicalHref("2.0.0+build.1"),
		}},
	}

	index := nextIndex(existing, "2.0.0+build.2", canonicalHref("2.0.0+build.2"))
	versions := []string{index.Versions[0].Version, index.Versions[1].Version}
	if want := []string{"2.0.0+build.2", "2.0.0+build.1"}; !reflect.DeepEqual(versions, want) {
		t.Fatalf("versions = %v, want %v", versions, want)
	}

	if index.Latest != "2.0.0+build.2" {
		t.Fatalf("latest = %q", index.Latest)
	}
}
