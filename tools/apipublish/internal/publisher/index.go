package publisher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"

	"github.com/MontFerret/specs/pkg/api"
)

func nextIndex(existing *api.Index, version, href string) *api.Index {
	versions := make([]api.IndexVersion, 0, 1)
	if existing != nil {
		versions = append(versions, existing.Versions...)
	}

	versions = append(versions, api.IndexVersion{Version: version, Href: href})
	sort.Slice(versions, func(i, j int) bool {
		left, _ := semver.StrictNewVersion(versions[i].Version)
		right, _ := semver.StrictNewVersion(versions[j].Version)
		if left.GreaterThan(right) {
			return true
		}

		if left.LessThan(right) {
			return false
		}

		return versions[i].Version > versions[j].Version
	})

	index := &api.Index{
		SchemaVersion: api.IndexSchemaVersion,
		Versions:      versions,
	}

	for _, entry := range versions {
		parsed, _ := semver.StrictNewVersion(entry.Version)
		if parsed.Prerelease() == "" {
			index.Latest = entry.Version

			break
		}
	}

	return index
}

func canonicalHref(version string) string {
	return "./versions/" + version + "/api.json"
}

func encodeIndex(index *api.Index) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(index); err != nil {
		return nil, fmt.Errorf("encode API Reference index: %w", err)
	}

	if !bytes.HasSuffix(buffer.Bytes(), []byte("\n")) {
		return nil, fmt.Errorf("encoded API Reference index is not newline terminated")
	}

	return buffer.Bytes(), nil
}
