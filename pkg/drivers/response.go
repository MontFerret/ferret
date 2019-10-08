package drivers

import (
	"context"
	"encoding/json"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// HTTPResponse HTTP response object.
type HTTPResponse struct {
	StatusCode int
	Status     string
	Headers    HTTPHeaders
}

func (resp *HTTPResponse) Type() core.Type {
	return HTTPResponseType
}

func (resp *HTTPResponse) String() string {
	return resp.Status
}

func (resp *HTTPResponse) Compare(other core.Value) int64 {
	if other.Type() != HTTPResponseType {
		return Compare(HTTPResponseType, other.Type())
	}

	// this is a safe cast. Only *HTTPResponse implements core.Value.
	// HTTPResponse does not.
	otherResp := other.(*HTTPResponse)

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
	return *(&resp)
}

func (resp *HTTPResponse) Hash() uint64 {
	return values.Parse(resp).Hash()
}

// responseMarshal is a structure that repeats HTTPResponse. It allows
// easily Marshal the HTTPResponse object.
type responseMarshal struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Headers    HTTPHeaders `json:"headers"`
}

func (resp *HTTPResponse) MarshalJSON() ([]byte, error) {
	if resp == nil {
		return json.Marshal(values.None)
	}

	return json.Marshal(responseMarshal(*resp))
}

func (resp *HTTPResponse) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return resp, nil
	}

	if typ := path[0].Type(); typ != types.String {
		return values.None, core.TypeError(typ, types.String)
	}

	field := path[0].(values.String).String()

	switch field {
	case "status":
		return values.NewString(resp.Status), nil
	case "statusCode":
		return values.NewInt(resp.StatusCode), nil
	case "headers":
		if len(path) == 1 {
			return resp.Headers, nil
		}

		return resp.Headers.GetIn(ctx, path[1:])
	}

	return values.None, nil
}
