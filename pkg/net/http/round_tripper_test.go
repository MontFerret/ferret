package http

import stdhttp "net/http"

type testRoundTripper func(*stdhttp.Request) (*stdhttp.Response, error)

func (fn testRoundTripper) RoundTrip(req *stdhttp.Request) (*stdhttp.Response, error) {
	return fn(req)
}
