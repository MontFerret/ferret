package html

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Download a resource from the given GetURL.
// @param GetURL (String) - GetURL to download.
// @returns data (Binary) - Returns a base64 encoded string in binary format.
func Download(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	err = core.ValidateType(arg1, types.String)

	if err != nil {
		return values.None, err
	}

	resp, err := http.Get(arg1.String())

	if err != nil {
		return values.None, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return values.None, err
	}

	return values.NewBinary(data), nil
}
