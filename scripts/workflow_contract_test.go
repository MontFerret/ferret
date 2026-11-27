package scripts_test

import (
	"fmt"
	"os"
	"os/exec"
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

	for _, required := range []string{"needs: main-baseline", "suite: [unit, integration]", "sample: [1, 2, 3, 4, 5]"} {
		if !strings.Contains(benchmark, required) {
			t.Fatalf("baseline publication is missing %q", required)
		}
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

func TestPublishedV2ReleaseNotifiesWebsiteThroughSharedWorkflow(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	workflow := readWorkflow(t, root, "notify-website.yml")
	for _, required := range []string{
		"types: [published]",
		"contents: read",
		"if: startsWith(github.event.release.tag_name, 'v2.')",
		"uses: MontFerret/.github/.github/workflows/notify-website-release.yml@main",
		"tag: ${{ github.event.release.tag_name }}",
		"release-url: ${{ github.event.release.html_url }}",
		"commit-sha: ${{ github.sha }}",
		"FERRET_RELEASE_APP_CLIENT_ID: ${{ secrets.FERRET_RELEASE_APP_CLIENT_ID }}",
		"FERRET_RELEASE_APP_PRIVATE_KEY: ${{ secrets.FERRET_RELEASE_APP_PRIVATE_KEY }}",
	} {
		if !strings.Contains(workflow, required) {
			t.Fatalf("website notification workflow is missing %q", required)
		}
	}

	for _, forbidden := range []string{
		"actions/checkout",
		"actions/create-github-app-token",
		"data/versions.yaml",
		"ferret-release-published",
		"gh pr",
		"workflow_dispatch:",
	} {
		if strings.Contains(workflow, forbidden) {
			t.Fatalf("website notification workflow duplicates downstream behavior %q", forbidden)
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

func TestCombineBenchmarkSamples(t *testing.T) {
	script, err := filepath.Abs("combine-benchmark-samples.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, fifth string }{
		{"complete", "BenchmarkObject 5 32 ns/op\n"},
		{"missing", ""}, {"empty", ""}, {"log-only", "ok github.com/MontFerret/ferret/v2/pkg/runtime\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			output := filepath.Join(dir, "benchmark-unit-output.txt")
			var want strings.Builder
			for i := 1; i <= 5; i++ {
				path := filepath.Join(dir, fmt.Sprintf("benchmark-output-main-unit-%d", i))
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
				data := fmt.Sprintf("BenchmarkObject %d 32 ns/op\n", i)
				if i == 5 {
					if tc.name == "missing" {
						continue
					}
					data = tc.fifth
				}
				if err := os.WriteFile(filepath.Join(path, filepath.Base(output)), []byte(data), 0o644); err != nil {
					t.Fatal(err)
				}
				want.WriteString(data)
			}
			log, err := exec.Command("bash", script, dir, "unit", output).CombinedOutput()
			if tc.name != "complete" {
				if err == nil || !strings.Contains(string(log), "unit sample 5") {
					t.Fatalf("incomplete suite returned %v, %s", err, log)
				}
				if _, err := os.Stat(output); !os.IsNotExist(err) {
					t.Fatalf("partial baseline published: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("combine: %v, %s", err, log)
			}
			got, err := os.ReadFile(output)
			if err != nil || string(got) != want.String() {
				t.Fatalf("combined samples = %q, %v; want %q", got, err, want.String())
			}
		})
	}
}
