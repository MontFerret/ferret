package scripts_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEveryPagesWriterUsesRepositoryWideSerialization(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	benchmark := readWorkflow(t, root, "benchmark.yml")
	publication := readWorkflow(t, root, "publish-core-api.yml")
	for name, workflow := range map[string]string{"benchmark": benchmark, "publication": publication} {
		if !strings.Contains(workflow, "group: gh-pages-writer\n") || !strings.Contains(workflow, "cancel-in-progress: false\n") {
			t.Fatalf("%s workflow does not use the shared non-cancelling Pages writer key", name)
		}
	}

	if strings.Contains(benchmark, "auto-push: true") {
		t.Fatal("benchmark workflow still contains an independent auto-push writer")
	}

	if !strings.Contains(benchmark, "needs:\n      - main-baseline-unit\n      - main-baseline-integration") {
		t.Fatal("benchmark persistence does not wait for both parallel baseline jobs")
	}
}

func TestReleasePublicationChecksPublishedExactTagAndUsesNormalPushWrapper(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	workflow := readWorkflow(t, root, "publish-core-api.yml")
	for _, required := range []string{
		"types: [published]",
		"workflow_dispatch:",
		"published_at != null",
		"ref: refs/tags/${{ steps.release.outputs.tag }}",
		"contents: write",
		"go -C tools/apiref run .",
		`-catalog "$RUNNER_TEMP/montferret-core-catalog.json"`,
		"./scripts/publish-core-api.sh",
	} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("release publication workflow is missing %q", required)
		}
	}
}

func readWorkflow(t *testing.T, root, name string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", name))
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}
