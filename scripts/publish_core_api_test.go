package scripts_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
)

func TestPublishCoreAPIRejectsStaleNonFastForwardPush(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	temporary := t.TempDir()
	remote := filepath.Join(temporary, "pages.git")
	git(t, temporary, "init", "--bare", remote)

	seed := filepath.Join(temporary, "seed")
	git(t, temporary, "clone", remote, seed)
	git(t, seed, "switch", "--orphan", "gh-pages")
	write(t, filepath.Join(seed, ".gitignore"), "public\n")
	configureAuthor(t, seed)
	git(t, seed, "add", ".gitignore")
	git(t, seed, "commit", "-m", "Initialize Pages")
	git(t, seed, "push", "origin", "gh-pages")

	first := filepath.Join(temporary, "first")
	second := filepath.Join(temporary, "second")
	git(t, temporary, "clone", "--branch", "gh-pages", remote, first)
	git(t, temporary, "clone", "--branch", "gh-pages", remote, second)
	configureAuthor(t, first)
	configureAuthor(t, second)

	firstReference := filepath.Join(temporary, "first-api.json")
	secondReference := filepath.Join(temporary, "second-api.json")
	writeReference(t, firstReference, "2.0.0-alpha.1")
	writeReference(t, secondReference, "2.0.0-alpha.2")

	publish := filepath.Join(root, "scripts", "publish-core-api.sh")
	runPublish(t, root, publish, firstReference, first, remote, true)
	remoteAfterFirst := strings.TrimSpace(git(t, temporary, "--git-dir", remote, "rev-parse", "refs/heads/gh-pages"))

	runPublish(t, root, publish, secondReference, second, remote, false)
	remoteAfterStale := strings.TrimSpace(git(t, temporary, "--git-dir", remote, "rev-parse", "refs/heads/gh-pages"))
	if remoteAfterStale != remoteAfterFirst {
		t.Fatalf("stale push changed published history: before=%s after=%s", remoteAfterFirst, remoteAfterStale)
	}

	assertContent(t, filepath.Join(first, ".gitignore"), "public\n")
	assertContent(t, filepath.Join(first, "versions", "2.0.0-alpha.1", "api.json"), string(read(t, firstReference)))
}

func runPublish(t *testing.T, root, script, reference, pages, remote string, wantSuccess bool) {
	t.Helper()

	command := exec.Command(script, reference, pages, remote, "gh-pages")
	command.Dir = root
	command.Env = append(os.Environ(), "GOWORK=off")
	output, err := command.CombinedOutput()
	if wantSuccess && err != nil {
		t.Fatalf("publish failed: %v\n%s", err, output)
	}

	if !wantSuccess && err == nil {
		t.Fatalf("stale publication unexpectedly succeeded:\n%s", output)
	}

	if !wantSuccess && !strings.Contains(string(output), "non-fast-forward") && !strings.Contains(string(output), "fetch first") {
		t.Fatalf("stale publication did not report non-fast-forward rejection:\n%s", output)
	}
}

func writeReference(t *testing.T, path, version string) {
	t.Helper()

	reference := &api.Reference{
		SchemaVersion: api.SchemaVersion,
		ID:            "montferret/core",
		Version:       version,
		Namespaces: []api.Namespace{{
			Name: "",
			Functions: []api.Function{{
				Name: "PING",
				Signatures: []api.Signature{{
					Parameters:  []api.Parameter{},
					Description: "Returns a value.",
					Return:      &api.Return{Type: "String", Description: "Value."},
				}},
			}},
		}},
	}

	data, err := json.MarshalIndent(reference, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func configureAuthor(t *testing.T, repository string) {
	t.Helper()

	git(t, repository, "config", "user.name", "Ferret Test")
	git(t, repository, "config", "user.email", "ferret-test@example.com")
}

func git(t *testing.T, directory string, arguments ...string) string {
	t.Helper()

	command := exec.Command("git", arguments...)
	command.Dir = directory
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(arguments, " "), err, output)
	}

	return string(output)
}

func write(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

func assertContent(t *testing.T, path, want string) {
	t.Helper()

	if got := string(read(t, path)); got != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
}
