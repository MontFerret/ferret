package drivers

import (
	"context"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// HTTPResponse HTTP response object.
type (
	HTTPResponse struct {
		URL          string
		StatusCode   int
		Status       string
		Headers      *HTTPHeaders
		Body         []byte
		ResponseTime float64
	}

	// responseMarshal is a structure that repeats HTTPResponse. It allows
	// easily Marshal the HTTPResponse object.
	responseMarshal struct {
		URL          string       `json:"url"`
		StatusCode   int          `json:"status_code"`
		Status       string       `json:"status"`
		Headers      *HTTPHeaders `json:"headers"`
		Body         []byte       `json:"body"`
		ResponseTime float64      `json:"response_time"`
	}
)

func (resp *HTTPResponse) Type() core.Type {
	return HTTPResponseType
}

func (resp *HTTPResponse) String() string {
	return resp.Status
}

func (resp *HTTPResponse) Compare(other core.Value) int64 {
	otherResp, ok := other.(*HTTPResponse)

	if !ok {
		return CompareTypes(HTTPResponseType, core.Reflect(other))
	}

	comp := resp.Headers.Compare(otherResp.Headers)

	if comp != 0 {
		return comp
	}

	// it makes no sense to compare Status strings
	// because they are always equal if StatusCode's are equal
	return values.NewInt(resp.StatusCode).
		Compare(values.NewInt(resp.StatusCode))
}

func (resp *HTTPResponse) Unwrap() interface{} {
	return resp
}

func (resp *HTTPResponse) Copy() core.Value {
	cop := *resp
	return &cop
}

func (resp *HTTPResponse) Hash() uint64 {
	return values.Parse(resp).Hash()
}

func (resp *HTTPResponse) MarshalJSON() ([]byte, error) {
	if resp == nil {
		return values.None.MarshalJSON()
	}

	return jettison.MarshalOpts(responseMarshal(*resp), jettison.NoHTMLEscaping())
}

func (resp *HTTPResponse) GetByKey(_ context.Context, key string) (core.Value, error) {
	if len(key) == 0 {
		return resp, nil
	}

	switch key {
	case "url", "URL":
		return values.NewString(resp.URL), nil
	case "status":
		return values.NewString(resp.Status), nil
	case "statusCode":
		return values.NewInt(resp.StatusCode), nil
	case "headers":
		return resp.Headers, nil
	case "body":
		return values.NewBinary(resp.Body), nil
	case "responseTime":
		return values.NewFloat(resp.ResponseTime), nil
	}

	return values.None, nil
}
