package html

import (
	"bytes"
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

func HtmlParse(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateType(args[0], types.String)

	if err != nil {
		return values.False, err
	}
	str := values.ToString(args[1])
	buf := bytes.NewBuffer([]byte(str))

	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a document")
	}

	return values.NewString(doc.Text()), nil
}
