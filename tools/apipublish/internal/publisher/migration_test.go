package publisher_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"

	"github.com/MontFerret/ferret/v2/tools/apipublish/internal/publisher"
)

func TestMigrateTypesCheckApplyAndIdempotency(t *testing.T) {
	root := t.TempDir()
	versions := []string{"2.0.0-alpha.2", "2.0.0-alpha.1"}
	writeLegacyHistory(t, root, versions)
	before := treeSnapshot(t, root)

	err := publisher.MigrateTypes(root, true)
	if !errors.Is(err, publisher.ErrTypeMigrationRequired) {
		t.Fatalf("check error = %v, want ErrTypeMigrationRequired", err)
	}

	if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("check mode changed the Pages tree\nbefore: %#v\nafter: %#v", before, after)
	}

	indexBefore := readFile(t, filepath.Join(root, "index.json"))
	catalogsBefore := make(map[string][]byte, len(versions))
	for _, version := range versions {
		catalogsBefore[version] = readFile(t, filepath.Join(root, "versions", version, "catalog.json"))
	}

	if err := publisher.MigrateTypes(root, false); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	if got := readFile(t, filepath.Join(root, "index.json")); !reflect.DeepEqual(got, indexBefore) {
		t.Fatal("migration changed index.json")
	}

	for _, version := range versions {
		catalogPath := filepath.Join(root, "versions", version, "catalog.json")
		if got := readFile(t, catalogPath); !reflect.DeepEqual(got, catalogsBefore[version]) {
			t.Fatalf("migration changed catalog for %s", version)
		}

		data := readFile(t, filepath.Join(root, "versions", version, "api.json"))
		reference, err := api.Parse(data)
		if err != nil {
			t.Fatalf("parse migrated API Reference %s: %v", version, err)
		}

		if reference.SchemaVersion != 1 || countSignatures(reference) != 2 {
			t.Fatalf("migrated API Reference %s changed schema or signature count: %#v", version, reference)
		}

		parameterType := reference.Namespaces[0].Functions[0].Signatures[0].Parameters[0].Type
		if parameterType.Kind != api.TypeKindUnion || len(parameterType.Types) != 2 {
			t.Fatalf("migrated parameter type %s = %#v", version, parameterType)
		}

		returnType := reference.Namespaces[0].Functions[0].Signatures[0].Return.Type
		if returnType.Kind != api.TypeKindList || returnType.Element.Name != "String" {
			t.Fatalf("migrated return type %s = %#v", version, returnType)
		}
	}

	after := treeSnapshot(t, root)
	if err := publisher.MigrateTypes(root, true); err != nil {
		t.Fatalf("check migrated history: %v", err)
	}

	if err := publisher.MigrateTypes(root, false); err != nil {
		t.Fatalf("reapply migration: %v", err)
	}

	if repeated := treeSnapshot(t, root); !reflect.DeepEqual(repeated, after) {
		t.Fatalf("idempotent migration changed tree\nfirst: %#v\nsecond: %#v", after, repeated)
	}
}

func TestMigrateTypesRejectsMalformedHistoryWithoutMutation(t *testing.T) {
	root := t.TempDir()
	version := "2.0.0-alpha.1"
	writeLegacyHistory(t, root, []string{version})
	artifactPath := filepath.Join(root, "versions", version, "api.json")
	data := legacyReferenceData(t, version)
	data = replaceFirstLegacyType(t, data, "String |")
	writeFile(t, artifactPath, data)
	before := treeSnapshot(t, root)

	err := publisher.MigrateTypes(root, false)
	if err == nil || !strings.Contains(err.Error(), "parse legacy type expression") {
		t.Fatalf("migration error = %v, want malformed type expression", err)
	}

	if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("malformed migration changed tree\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestMigrateTypesRollsBackWriteFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permission failure is platform-specific")
	}

	root := t.TempDir()
	versions := []string{"2.0.0-alpha.2", "2.0.0-alpha.1"}
	writeLegacyHistory(t, root, versions)
	before := treeSnapshot(t, root)
	blockedDirectory := filepath.Join(root, "versions", versions[1])
	if err := os.Chmod(blockedDirectory, 0o555); err != nil {
		t.Fatal(err)
	}

	err := publisher.MigrateTypes(root, false)
	if chmodErr := os.Chmod(blockedDirectory, 0o755); chmodErr != nil {
		t.Fatal(chmodErr)
	}

	if err == nil || !strings.Contains(err.Error(), "write migrated API Reference") {
		t.Fatalf("migration error = %v, want write failure", err)
	}

	if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("failed migration was not rolled back\nbefore: %#v\nafter: %#v", before, after)
	}
}

func TestPublishSucceedsAfterTypeMigration(t *testing.T) {
	root := t.TempDir()
	writeLegacyHistory(t, root, []string{"2.0.0-alpha.1"})
	if err := publisher.MigrateTypes(root, false); err != nil {
		t.Fatalf("migrate history: %v", err)
	}

	if err := publisher.Publish(
		root,
		referenceData(t, "2.0.0-alpha.2", "montferret/core"),
		catalogData(t, "2.0.0-alpha.2", "montferret/core"),
	); err != nil {
		t.Fatalf("publish after migration: %v", err)
	}

	if got := len(readIndex(t, root).Versions); got != 2 {
		t.Fatalf("published index has %d versions, want 2", got)
	}
}

func writeLegacyHistory(t *testing.T, root string, versions []string) {
	t.Helper()

	index := &api.Index{SchemaVersion: api.IndexSchemaVersion}
	for _, version := range versions {
		writeFile(t, filepath.Join(root, "versions", version, "api.json"), legacyReferenceData(t, version))
		writeFile(t, filepath.Join(root, "versions", version, "catalog.json"), catalogData(t, version, "montferret/core"))
		index.Versions = append(index.Versions, api.IndexVersion{Version: version, Href: "./versions/" + version + "/api.json"})
	}

	writeJSON(t, filepath.Join(root, "index.json"), index)
}

func legacyReferenceData(t *testing.T, version string) []byte {
	t.Helper()

	var document any
	if err := json.Unmarshal(referenceData(t, version, "montferret/core"), &document); err != nil {
		t.Fatal(err)
	}

	convertTypesToLegacy(document)
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	return append(data, '\n')
}

func convertTypesToLegacy(value any) {
	switch current := value.(type) {
	case map[string]any:
		for key, child := range current {
			if key == "type" {
				if structured, ok := child.(map[string]any); ok {
					current[key] = renderLegacyType(structured)

					continue
				}
			}

			convertTypesToLegacy(child)
		}
	case []any:
		for _, child := range current {
			convertTypesToLegacy(child)
		}
	}
}

func renderLegacyType(value map[string]any) string {
	switch value["kind"] {
	case "named":
		return value["name"].(string)
	case "union":
		members := value["types"].([]any)
		rendered := make([]string, 0, len(members))
		for _, member := range members {
			rendered = append(rendered, renderLegacyType(member.(map[string]any)))
		}

		return strings.Join(rendered, " | ")
	case "list":
		return "[" + renderLegacyType(value["element"].(map[string]any)) + "]"
	default:
		panic("unexpected structured type kind")
	}
}

func replaceFirstLegacyType(t *testing.T, data []byte, replacement string) []byte {
	t.Helper()

	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}

	namespaces := document["namespaces"].([]any)
	functions := namespaces[0].(map[string]any)["functions"].([]any)
	signatures := functions[0].(map[string]any)["signatures"].([]any)
	parameters := signatures[0].(map[string]any)["parameters"].([]any)
	parameters[0].(map[string]any)["type"] = replacement
	result, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	return append(result, '\n')
}

func countSignatures(reference *api.Reference) int {
	count := 0
	for _, namespace := range reference.Namespaces {
		for _, function := range namespace.Functions {
			count += len(function.Signatures)
		}
	}

	return count
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return data
}
