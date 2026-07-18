package http

import "testing"

func FuzzParseRequestURL(f *testing.F) {
	for _, seed := range []string{
		"https://example.com/path",
		"HTTP://127.1:8080",
		"https://user:password@example.com",
		"https://xn--xample-9ua.com",
		"https://K.example",
		"",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, rawURL string) {
		u, err := parseRequestURL(rawURL)
		if err != nil {
			return
		}
		if u.Scheme == "" || u.Host == "" {
			t.Fatalf("successful parse returned incomplete URL: %#v", u)
		}
	})
}

func FuzzParseIPAddress(f *testing.F) {
	for _, seed := range []string{
		"127.1",
		"2130706433",
		"8.8.8.8",
		"::ffff:127.0.0.1",
		"64:ff9b::8.8.8.8",
		"not-an-address",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, hostname string) {
		addr, ok := parseIPAddress(hostname)
		if !ok {
			return
		}
		if !addr.IsValid() {
			t.Fatalf("successful parse returned invalid address for %q", hostname)
		}
		if addr != addr.Unmap() {
			t.Fatalf("successful parse returned mapped address %q", addr)
		}
	})
}

func FuzzNormalizeHostValue(f *testing.F) {
	for _, seed := range []string{
		"Example.COM",
		"example.com:8443",
		"[::ffff:8.8.8.8]:443",
		"127.1",
		"K.example",
		"",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(_ *testing.T, host string) {
		_ = normalizeHostValue(host)
	})
}
