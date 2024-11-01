package drivers

import (
	"context"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// HTTPRequest HTTP request object.
type (
	HTTPRequest struct {
		URL     string
		Method  string
		Headers *HTTPHeaders
		Body    []byte
	}
	// requestMarshal is a structure that repeats HTTPRequest. It allows
	// easily Marshal the HTTPRequest object.
	requestMarshal struct {
		URL     string       `json:"url"`
		Method  string       `json:"method"`
		Headers *HTTPHeaders `json:"headers"`
		Body    []byte       `json:"body"`
	}
)

func (req *HTTPRequest) MarshalJSON() ([]byte, error) {
	if req == nil {
		return values.None.MarshalJSON()
	}

	return jettison.MarshalOpts(requestMarshal(*req), jettison.NoHTMLEscaping())
}

func (req *HTTPRequest) Type() core.Type {
	return HTTPRequestType
}

func (req *HTTPRequest) String() string {
	return req.URL
}

func (req *HTTPRequest) Compare(other core.Value) int64 {
	otherReq, ok := other.(*HTTPRequest)

	if !ok {
		return CompareTypes(HTTPRequestType, core.Reflect(other))
	}

	comp := req.Headers.Compare(otherReq.Headers)

	if comp != 0 {
		return comp
	}

	comp = values.NewString(req.Method).Compare(values.NewString(otherReq.Method))

	if comp != 0 {
		return comp
	}

	return values.NewString(req.URL).
		Compare(values.NewString(otherReq.URL))
}

func (req *HTTPRequest) Unwrap() interface{} {
	return req
}

func (req *HTTPRequest) Hash() uint64 {
	return values.Parse(req).Hash()
}

func (req *HTTPRequest) Copy() core.Value {
	cop := *req
	return &cop
}

func (req *HTTPRequest) GetByKey(ctx context.Context, key string) (core.Value, error) {
	if len(key) == 0 {
		return req, nil
	}

	switch key {
	case "url", "URL":
		return values.NewString(req.URL), nil
	case "method":
		return values.NewString(req.Method), nil
	case "headers":
		return req.Headers, nil
	case "body":
		return values.NewBinary(req.Body), nil
	}

	return values.None, nil
}
