package http

import (
	"context"
	"errors"
	"io"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestClientDoSuccessDefaultsAndMaterializesResponse(t *testing.T) {
	var seenMethod string

	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		seenMethod = r.Method
		w.Header().Add("X-Result", "ok")
		w.WriteHeader(stdhttp.StatusCreated)
		_, _ = w.Write([]byte("created"))
	}))
	defer server.Close()

	res, err := New(WithAllowLocalhost(true)).Do(nil, &Request{URL: server.URL})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if seenMethod != stdhttp.MethodGet {
		t.Fatalf("expected default GET method, got %q", seenMethod)
	}

	if res.StatusCode != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", res.StatusCode)
	}

	if got := string(res.Body); got != "created" {
		t.Fatalf("expected response body %q, got %q", "created", got)
	}

	if got := res.Headers["X-Result"]; len(got) != 1 || got[0] != "ok" {
		t.Fatalf("expected copied response header, got %v", got)
	}
}

func TestClientDoSendsHeadersAndBodyWithPolicy(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if got := r.Header.Values("X-Token"); len(got) != 2 || got[0] != "a" || got[1] != "b" {
			t.Errorf("expected repeated X-Token headers, got %v", got)
			stdhttp.Error(w, "bad token header", stdhttp.StatusBadRequest)
			return
		}
		if got := r.Header.Get("X-Default"); got != "default" {
			t.Errorf("expected default header, got %q", got)
			stdhttp.Error(w, "bad default header", stdhttp.StatusBadRequest)
			return
		}
		if got := r.Header.Get("X-Blocked"); got != "" {
			t.Errorf("expected blocked header to be omitted, got %q", got)
			stdhttp.Error(w, "blocked header present", stdhttp.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			stdhttp.Error(w, "read body failed", stdhttp.StatusInternalServerError)
			return
		}
		if got := string(body); got != "payload" {
			t.Errorf("expected request body %q, got %q", "payload", got)
			stdhttp.Error(w, "bad body", stdhttp.StatusBadRequest)
			return
		}

		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client := New(
		WithDefaultHeader("X-Default", "default"),
		WithBlockedRequestHeaders("X-Blocked"),
		WithAllowLocalhost(true),
	)

	_, err := client.Do(context.Background(), &Request{
		Method: stdhttp.MethodPost,
		URL:    server.URL,
		Headers: Headers{
			"X-Token":   {"a", "b"},
			"X-Blocked": {"secret"},
		},
		Body: []byte("payload"),
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestClientDoRequestBodyLimit(t *testing.T) {
	_, err := New(WithMaxRequestSize(3)).Do(
		context.Background(),
		&Request{
			Method: stdhttp.MethodPost,
			URL:    "http://example.com",
			Body:   []byte("four"),
		},
	)
	if err == nil || !strings.Contains(err.Error(), "request body exceeds limit") {
		t.Fatalf("expected request body limit error, got %v", err)
	}
}

func TestClientDoResponseBodyLimit(t *testing.T) {
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		_, _ = w.Write([]byte("four"))
	}))
	defer server.Close()

	_, err := New(
		WithMaxResponseSize(3),
		WithAllowLocalhost(true),
	).Do(
		context.Background(),
		&Request{URL: server.URL},
	)
	if err == nil || !strings.Contains(err.Error(), "response body exceeds limit") {
		t.Fatalf("expected response body limit error, got %v", err)
	}
}

func TestClientDoRedirectPolicy(t *testing.T) {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("/start", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, "/middle", stdhttp.StatusFound)
	})
	mux.HandleFunc("/middle", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		stdhttp.Redirect(w, r, "/done", stdhttp.StatusFound)
	})
	mux.HandleFunc("/done", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		_, _ = w.Write([]byte("done"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	res, err := New(WithAllowLocalhost(true)).Do(
		context.Background(),
		&Request{URL: server.URL + "/start"},
	)
	if err != nil {
		t.Fatalf("default redirect request failed: %v", err)
	}
	if got := string(res.Body); got != "done" {
		t.Fatalf("expected followed redirect body %q, got %q", "done", got)
	}

	res, err = New(
		WithFollowRedirects(false),
		WithAllowLocalhost(true),
	).Do(
		context.Background(),
		&Request{URL: server.URL + "/start"},
	)
	if err != nil {
		t.Fatalf("no-follow redirect request failed: %v", err)
	}
	if res.StatusCode != stdhttp.StatusFound {
		t.Fatalf("expected redirect response status, got %d", res.StatusCode)
	}

	_, err = New(
		WithMaxRedirects(1),
		WithAllowLocalhost(true),
	).Do(
		context.Background(),
		&Request{URL: server.URL + "/start"},
	)
	if err == nil || !strings.Contains(err.Error(), "stopped after 1 redirect") {
		t.Fatalf("expected max redirect error, got %v", err)
	}
}

func TestClientDoValidatesRequest(t *testing.T) {
	tests := []struct {
		req  *Request
		name string
		want string
	}{
		{name: "nil", req: nil, want: ErrNilRequest.Error()},
		{name: "invalid method", req: &Request{Method: "BAD METHOD", URL: "http://example.com"}, want: "invalid method"},
		{name: "missing url", req: &Request{}, want: "url is required"},
		{name: "missing scheme", req: &Request{URL: "example.com"}, want: "url scheme is required"},
		{name: "missing host", req: &Request{URL: "http:///path"}, want: "url host is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New().Do(context.Background(), tt.req)
			if err == nil {
				t.Fatal("expected error")
			}
			if tt.req == nil {
				if !errors.Is(err, ErrNilRequest) {
					t.Fatalf("expected ErrNilRequest, got %v", err)
				}
				return
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestClientDoPolicyURLChecks(t *testing.T) {
	localServer := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer localServer.Close()

	tests := []struct {
		client Client
		name   string
		url    string
		want   string
	}{
		{
			name:   "scheme",
			client: New(WithAllowedSchemes("https")),
			url:    "http://example.com",
			want:   "scheme",
		},
		{
			name:   "blocked host",
			client: New(WithBlockedHosts("example.com")),
			url:    "http://example.com",
			want:   "blocked",
		},
		{
			name:   "allowed host",
			client: New(WithAllowedHosts("allowed.example")),
			url:    "http://other.example",
			want:   "not allowed",
		},
		{
			name:   "localhost",
			client: New(),
			url:    localServer.URL,
			want:   "localhost is not allowed",
		},
		{
			name:   "private network",
			client: New(),
			url:    "http://10.0.0.1",
			want:   "private network",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.client.Do(context.Background(), &Request{URL: tt.url})
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}

func TestClientDoDefaultRejectsLoopbackBeforeRequest(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		requests.Add(1)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	_, err := New().Do(context.Background(), &Request{URL: server.URL})
	if err == nil || !strings.Contains(err.Error(), "localhost is not allowed") {
		t.Fatalf("expected default loopback policy error, got %v", err)
	}
	if got := requests.Load(); got != 0 {
		t.Fatalf("expected blocked request not to reach server, got %d request(s)", got)
	}

	res, err := New(WithAllowLocalhost(true)).Do(
		context.Background(),
		&Request{URL: server.URL},
	)
	if err != nil {
		t.Fatalf("expected explicit localhost access to succeed, got %v", err)
	}
	if got := string(res.Body); got != "ok" {
		t.Fatalf("expected response body %q, got %q", "ok", got)
	}
	if got := requests.Load(); got != 1 {
		t.Fatalf("expected one explicitly allowed request, got %d", got)
	}
}
