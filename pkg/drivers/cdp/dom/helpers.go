package dom

import (
	"bytes"
	"context"
	"github.com/rs/zerolog"
	"regexp"
	"strings"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
)

var camelMatcher = regexp.MustCompile("[A-Za-z0-9]+")

func resolveFrame(ctx context.Context, logger zerolog.Logger, client *cdp.Client, frameID page.FrameID) (dom.Node, *eval.Runtime, error) {
	exec, err := eval.New(ctx, logger, client, frameID)

	if err != nil {
		return dom.Node{}, nil, errors.Wrap(err, "create JS executor")
	}

	evalRes, err := exec.EvalRef(ctx, eval.F("return document"))

	if err != nil {
		return dom.Node{}, nil, err
	}

	if evalRes.ObjectID == nil {
		return dom.Node{}, nil, errors.New("failed to resolve frame document")
	}

	req, err := client.DOM.RequestNode(ctx, dom.NewRequestNodeArgs(*evalRes.ObjectID))

	if err != nil {
		return dom.Node{}, nil, err
	}

	if req.NodeID == 0 {
		return dom.Node{}, nil, errors.New("framed document is resolved with empty node id")
	}

	desc, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetNodeID(req.NodeID).
			SetDepth(1),
	)

	if err != nil {
		return dom.Node{}, nil, err
	}

	// Returned node, by some reason, does not contain the NodeID
	// So, we have to set it manually
	desc.Node.NodeID = req.NodeID

	return desc.Node, exec, nil
}

func toCamelCase(input string) string {
	var buf bytes.Buffer

	matched := camelMatcher.FindAllString(input, -1)

	if matched == nil {
		return ""
	}

	for i, match := range matched {
		res := match

		if i > 0 {
			if len(match) > 1 {
				res = strings.ToUpper(match[0:1]) + match[1:]
			} else {
				res = strings.ToUpper(match)
			}
		}

		buf.WriteString(res)
	}

	return buf.String()
}
