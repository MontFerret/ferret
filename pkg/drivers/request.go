package drivers

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/wI2L/jettison"
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
	if other.Type() != HTTPRequestType {
		return Compare(HTTPResponseType, other.Type())
	}

	// this is a safe cast. Only *HTTPRequest implements core.Value.
	// HTTPRequest does not.
	otherReq := other.(*HTTPRequest)

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
	return *(&req)
}

func (req *HTTPRequest) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	if len(path) == 0 {
		return req, nil
	}

	segmentIdx := 0
	segment := path[segmentIdx]

	if typ := segment.Type(); typ != types.String {
		return values.None, core.NewPathError(core.TypeError(typ, types.String), segmentIdx)
	}

	field := segment.String()

	switch field {
	case "url", "URL":
		return values.NewString(req.URL), nil
	case "method":
		return values.NewString(req.Method), nil
	case "headers":
		if len(path) == 1 {
			return req.Headers, nil
		}

		out, pathErr := req.Headers.GetIn(ctx, path[1:])

		if pathErr != nil {
			return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
		}

		return out, nil
	case "body":
		return values.NewBinary(req.Body), nil
	}

	return values.None, nil
}
