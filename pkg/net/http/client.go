package http

import (
	"context"
	stdhttp "net/http"
)

type (
	// Client executes HTTP requests for Ferret host integrations.
	Client interface {
		Do(ctx context.Context, req *Request) (*Response, error)
	}

	defaultHTTPClient struct {
		policy    *Policies
		transport stdhttp.Client
	}
)

// New constructs an HTTP client with the provided policies.
func New(setters ...Policy) Client {
	policies := NewPolicies(setters...)

	return &defaultHTTPClient{
		policy: &policies,
	}
}

func (d *defaultHTTPClient) Do(ctx context.Context, req *Request) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if req == nil {
		return nil, ErrNilRequest
	}

	p := d.policy
	if p == nil {
		policies := NewPolicies()
		p = &policies
	}

	stdReq, err := toStdRequest(ctx, req, p)
	if err != nil {
		return nil, err
	}

	client := d.transport
	if p.Timeout > 0 {
		client.Timeout = p.Timeout
	}
	client.CheckRedirect = p.checkRedirect

	res, err := client.Do(stdReq)
	if err != nil {
		return nil, err
	}

	return fromStdResponse(res, p)
}
