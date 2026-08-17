package publisher_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"

	"github.com/MontFerret/ferret/v2/tools/apipublish/internal/publisher"
)

func TestPublishCreatesPrereleaseIndexAndPreservesUnrelatedPagesFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".gitignore"), []byte("public\n"))
	writeFile(t, filepath.Join(root, "dev", "benchmarks.json"), []byte("{\"score\":1}\n"))
	reference := referenceData(t, "2.0.0-alpha.45", "montferret/core")
	catalog := catalogData(t, "2.0.0-alpha.45", "montferret/core")

	if err := publisher.Publish(root, reference, catalog); err != nil {
		t.Fatal(err)
	}

	index := readIndex(t, root)
	if index.Latest != "" {
		t.Fatalf("prerelease-only latest = %q, want omission", index.Latest)
	}

	if got, want := index.Versions, []api.IndexVersion{{Version: "2.0.0-alpha.45", Href: "./versions/2.0.0-alpha.45/api.json"}}; !reflect.DeepEqual(got, want) {
		t.Fatalf("versions = %#v, want %#v", got, want)
	}

	artifact, err := os.ReadFile(filepath.Join(root, "versions", "2.0.0-alpha.45", "api.json"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(artifact, reference) {
		t.Fatal("published artifact bytes differ from generated reference")
	}

	publishedCatalog, err := os.ReadFile(filepath.Join(root, "versions", "2.0.0-alpha.45", "catalog.json"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(publishedCatalog, catalog) {
		t.Fatal("published catalog bytes differ from generated catalog")
	}

	assertFile(t, filepath.Join(root, ".gitignore"), "public\n")
	assertFile(t, filepath.Join(root, "dev", "benchmarks.json"), "{\"score\":1}\n")
}

func TestPublishAddsVersionsInDeterministicOrderAndSetsStableLatest(t *testing.T) {
	root := t.TempDir()
	for _, version := range []string{"2.0.0-alpha.1", "1.9.0", "2.0.0-alpha.2", "2.0.0"} {
		if err := publisher.Publish(root, referenceData(t, version, "montferret/core"), catalogData(t, version, "montferret/core")); err != nil {
			t.Fatalf("publish %s: %v", version, err)
		}
	}

	index := readIndex(t, root)
	if index.Latest != "2.0.0" {
		t.Fatalf("latest = %q, want 2.0.0", index.Latest)
	}

	versions := make([]string, 0, len(index.Versions))
	for _, entry := range index.Versions {
		versions = append(versions, entry.Version)
		if entry.Href != "./versions/"+entry.Version+"/api.json" {
			t.Fatalf("non-authoritative href for %s: %s", entry.Version, entry.Href)
		}
	}

	want := []string{"2.0.0", "2.0.0-alpha.2", "2.0.0-alpha.1", "1.9.0"}
	if !reflect.DeepEqual(versions, want) {
		t.Fatalf("versions = %v, want %v", versions, want)
	}
}

func TestPublishRejectsImmutableVersionWithoutMutation(t *testing.T) {
	root := t.TempDir()
	reference := referenceData(t, "2.0.0-alpha.1", "montferret/core")
	catalog := catalogData(t, "2.0.0-alpha.1", "montferret/core")
	if err := publisher.Publish(root, reference, catalog); err != nil {
		t.Fatal(err)
	}

	before := treeSnapshot(t, root)
	err := publisher.Publish(root, reference, catalog)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("error = %v, want immutable-version rejection", err)
	}

	if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("tree changed after immutable rejection\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestPublishRejectsExistingVersionDirectoryEvenWithoutIndexEntry(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "versions", "2.0.0-alpha.1"), 0o755); err != nil {
		t.Fatal(err)
	}

	err := publisher.Publish(root, referenceData(t, "2.0.0-alpha.1", "montferret/core"), catalogData(t, "2.0.0-alpha.1", "montferret/core"))
	if err == nil || !strings.Contains(err.Error(), "unindexed API Reference version") {
		t.Fatalf("error = %v, want unindexed immutable directory rejection", err)
	}
}

func TestPublishRejectsInvalidIncomingIdentity(t *testing.T) {
	root := t.TempDir()
	err := publisher.Publish(root, referenceData(t, "2.0.0-alpha.1", "other/module"), catalogData(t, "2.0.0-alpha.1", "other/module"))
	if err == nil || !strings.Contains(err.Error(), `want "montferret/core"`) {
		t.Fatalf("error = %v, want identity rejection", err)
	}

	if entries, readErr := os.ReadDir(root); readErr != nil || len(entries) != 0 {
		t.Fatalf("invalid identity mutated root: entries=%v err=%v", entries, readErr)
	}
}

func TestPublishRejectsInvalidExistingStateBeforeMutation(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*testing.T, string)
		reason string
	}{
		{name: "malformed index", reason: "parse existing", mutate: func(t *testing.T, root string) {
			writeFile(t, filepath.Join(root, "index.json"), []byte("{}\n"))
		}},
		{name: "missing artifact", reason: "inspect existing API Reference", mutate: func(t *testing.T, root string) {
			if err := os.Remove(filepath.Join(root, "versions", "2.0.0-alpha.1", "api.json")); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "wrong artifact identity", reason: "identifies other/module", mutate: func(t *testing.T, root string) {
			writeFile(t, filepath.Join(root, "versions", "2.0.0-alpha.1", "api.json"), referenceData(t, "2.0.0-alpha.1", "other/module"))
		}},
		{name: "non-authoritative href", reason: "authoritative href", mutate: func(t *testing.T, root string) {
			index := readIndex(t, root)
			index.Versions[0].Href = "./api/2.0.0-alpha.1.json"
			writeJSON(t, filepath.Join(root, "index.json"), index)
		}},
		{name: "unindexed version", reason: "unindexed API Reference version", mutate: func(t *testing.T, root string) {
			if err := os.Mkdir(filepath.Join(root, "versions", "2.0.0-alpha.0"), 0o755); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "extra version file", reason: "optional catalog.json only", mutate: func(t *testing.T, root string) {
			writeFile(t, filepath.Join(root, "versions", "2.0.0-alpha.1", "notes.txt"), []byte("unexpected\n"))
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			if test.name != "malformed index" {
				if err := publisher.Publish(root, referenceData(t, "2.0.0-alpha.1", "montferret/core"), catalogData(t, "2.0.0-alpha.1", "montferret/core")); err != nil {
					t.Fatal(err)
				}
			}

			test.mutate(t, root)
			before := treeSnapshot(t, root)
			err := publisher.Publish(root, referenceData(t, "2.0.0-alpha.2", "montferret/core"), catalogData(t, "2.0.0-alpha.2", "montferret/core"))
			if err == nil || !strings.Contains(err.Error(), test.reason) {
				t.Fatalf("error = %v, want reason containing %q", err, test.reason)
			}

			if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
				t.Fatalf("invalid existing state mutated tree\nbefore: %#v\nafter: %#v", before, after)
			}
		})
	}
}

func TestPublishAcceptsExistingLegacyAPIOnlyVersion(t *testing.T) {
	root := t.TempDir()
	legacyVersion := "2.0.0-alpha.45"
	writeFile(t, filepath.Join(root, "versions", legacyVersion, "api.json"), referenceData(t, legacyVersion, "montferret/core"))
	writeJSON(t, filepath.Join(root, "index.json"), &api.Index{
		SchemaVersion: api.IndexSchemaVersion,
		Versions:      []api.IndexVersion{{Version: legacyVersion, Href: "./versions/" + legacyVersion + "/api.json"}},
	})

	newVersion := "2.0.0-alpha.47"
	if err := publisher.Publish(root, referenceData(t, newVersion, "montferret/core"), catalogData(t, newVersion, "montferret/core")); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(root, "versions", legacyVersion, "catalog.json")); !os.IsNotExist(err) {
		t.Fatalf("legacy version was backfilled: %v", err)
	}
}

func TestPublishRejectsInvalidIncomingCatalogWithoutMutation(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		catalog []byte
	}{
		{name: "malformed", catalog: []byte("{\n"), want: "parse generated API Catalog"},
		{name: "identity mismatch", catalog: catalogData(t, "2.0.0-alpha.47", "other/module"), want: "does not match API id"},
		{name: "version mismatch", catalog: catalogData(t, "2.0.0-alpha.46", "montferret/core"), want: "does not match API version"},
		{name: "unknown namespace", catalog: mutateCatalogData(t, "2.0.0-alpha.47", func(value *apicatalog.Catalog) {
			value.Categories[0].Functions[0].Namespace = "io::missing"
		}), want: "unknown API namespace"},
		{name: "unknown function", catalog: mutateCatalogData(t, "2.0.0-alpha.47", func(value *apicatalog.Catalog) {
			value.Categories[0].Functions[0].Name = "MISSING"
		}), want: "unknown function"},
		{name: "uncategorized namespaced function", catalog: mutateCatalogData(t, "2.0.0-alpha.47", func(value *apicatalog.Catalog) {
			value.Categories = value.Categories[1:]
		}), want: `function "io::fs::READ" is not assigned`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			err := publisher.Publish(root, referenceData(t, "2.0.0-alpha.47", "montferret/core"), test.catalog)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Publish error = %v, want %q", err, test.want)
			}

			if entries, readErr := os.ReadDir(root); readErr != nil || len(entries) != 0 {
				t.Fatalf("invalid catalog mutated root: entries=%v err=%v", entries, readErr)
			}
		})
	}
}

func referenceData(t *testing.T, version, id string) []byte {
	t.Helper()

	reference := &api.Reference{
		SchemaVersion: api.SchemaVersion,
		ID:            id,
		Version:       version,
		Namespaces: []api.Namespace{
			{
				Name: "",
				Functions: []api.Function{{
					Name: "PING",
					Signatures: []api.Signature{{
						Parameters:  []api.Parameter{},
						Description: "Returns a value.",
						Return:      &api.Return{Type: "String", Description: "Value."},
					}},
				}},
			},
			{
				Name: "io::fs",
				Functions: []api.Function{{
					Name: "READ",
					Signatures: []api.Signature{{
						Parameters:  []api.Parameter{},
						Description: "Reads a file.",
						Return:      &api.Return{Type: "String", Description: "Contents."},
					}},
				}},
			},
		},
	}

	data, err := json.MarshalIndent(reference, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	return append(data, '\n')
}

func catalogData(t *testing.T, version, id string) []byte {
	t.Helper()

	catalog := &apicatalog.Catalog{
		SchemaVersion: apicatalog.SchemaVersion,
		ID:            id,
		Version:       version,
		Categories: []apicatalog.Category{
			{
				ID:          "io",
				Title:       "I/O",
				Description: "Input and output functions.",
				Functions:   []apicatalog.FunctionRef{{Namespace: "io::fs", Name: "READ"}},
			},
			{
				ID:          "utils",
				Title:       "Utilities",
				Description: "General utility functions.",
				Functions:   []apicatalog.FunctionRef{{Namespace: "", Name: "PING"}},
			},
		},
	}

	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	return append(data, '\n')
}

func mutateCatalogData(t *testing.T, version string, mutate func(*apicatalog.Catalog)) []byte {
	t.Helper()

	catalog, err := apicatalog.Parse(catalogData(t, version, "montferret/core"))
	if err != nil {
		t.Fatal(err)
	}
	mutate(catalog)

	data, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	return append(data, '\n')
}

func readIndex(t *testing.T, root string) *api.Index {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(root, "index.json"))
	if err != nil {
		t.Fatal(err)
	}

	index, err := api.ParseIndex(data)
	if err != nil {
		t.Fatal(err)
	}

	return index
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	writeFile(t, path, append(data, '\n'))
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != want {
		t.Fatalf("%s = %q, want %q", path, data, want)
	}
}

func treeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()

	snapshot := make(map[string]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			snapshot[relative+"/"] = ""

			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		snapshot[relative] = string(data)

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	return snapshot
}
