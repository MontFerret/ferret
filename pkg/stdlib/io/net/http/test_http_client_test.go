package http_test

import (
	"bytes"
	"context"
	"io"
	stdhttp "net/http"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type testHTTPClient struct{}

func (testHTTPClient) Do(ctx context.Context, req *ferrethttp.Request) (*ferrethttp.Response, error) {
	stdReq, err := stdhttp.NewRequestWithContext(
		ctx,
		req.Method,
		req.URL,
		bytes.NewReader(req.Body),
	)
	if err != nil {
		return nil, err
	}

	for key, values := range req.Headers {
		for _, value := range values {
			stdReq.Header.Add(key, value)
		}
	}

	res, err := stdhttp.DefaultClient.Do(stdReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &ferrethttp.Response{
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Headers:    ferrethttp.Headers(res.Header),
		Body:       body,
	}, nil
}
