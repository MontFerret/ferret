package dom

import (
	"bytes"
	"context"
	"regexp"
	"strings"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var camelMatcher = regexp.MustCompile("[A-Za-z0-9]+")

// traverseAttrs is a helper function that parses a given interleaved array of node attribute names and values,
// and calls a given attribute on each key-value pair
func traverseAttrs(attrs []string, predicate func(name, value string) bool) {
	count := len(attrs)

	for i := 0; i < count; i++ {
		if i%2 != 0 {
			if predicate(attrs[i-1], attrs[i]) == false {
				break
			}
		}
	}
}

func setInnerHTML(ctx context.Context, client *cdp.Client, exec *eval.Runtime, id HTMLElementIdentity, innerHTML values.String) error {
	var objID runtime.RemoteObjectID

	if id.ObjectID != "" {
		objID = id.ObjectID
	} else {
		repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

		if err != nil {
			return err
		}

		if repl.Object.ObjectID == nil {
			return errors.New("unable to resolve node")
		}

		objID = *repl.Object.ObjectID
	}

	return exec.Eval(
		ctx,
		templates.SetInnerHTML(),
		eval.WithArgRef(objID),
		eval.WithArgValue(innerHTML),
	)

}

func getInnerHTML(ctx context.Context, client *cdp.Client, exec *eval.Runtime, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
	// not a document
	if nodeType != html.DocumentNode {
		var objID runtime.RemoteObjectID

		if id.ObjectID != "" {
			objID = id.ObjectID
		} else {
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

			if err != nil {
				return "", err
			}

			if repl.Object.ObjectID == nil {
				return "", errors.New("unable to resolve node")
			}

			objID = *repl.Object.ObjectID
		}

		res, err := exec.ReadProperty(ctx, objID, "innerHTML")

		if err != nil {
			return "", err
		}

		return values.NewString(res.String()), nil
	}

	repl, err := exec.EvalValue(ctx, "return document.documentElement.innerHTML")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

func setInnerText(ctx context.Context, client *cdp.Client, exec *eval.Runtime, id HTMLElementIdentity, innerText values.String) error {
	var objID runtime.RemoteObjectID

	if id.ObjectID != "" {
		objID = id.ObjectID
	} else {
		repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

		if err != nil {
			return err
		}

		if repl.Object.ObjectID == nil {
			return errors.New("unable to resolve node")
		}

		objID = *repl.Object.ObjectID
	}

	return exec.Eval(
		ctx,
		templates.SetInnerText(),
		eval.WithArgRef(objID),
		eval.WithArgValue(innerText),
	)
}

func getInnerText(ctx context.Context, client *cdp.Client, exec *eval.Runtime, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
	// not a document
	if nodeType != html.DocumentNode {
		var objID runtime.RemoteObjectID

		if id.ObjectID != "" {
			objID = id.ObjectID
		} else {
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

			if err != nil {
				return "", err
			}

			if repl.Object.ObjectID == nil {
				return "", errors.New("unable to resolve node")
			}

			objID = *repl.Object.ObjectID
		}

		res, err := exec.ReadProperty(ctx, objID, "innerText")

		if err != nil {
			return "", err
		}

		return values.NewString(res.String()), err
	}

	repl, err := exec.EvalValue(ctx, "return document.documentElement.innerText")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

func resolveFrame(ctx context.Context, client *cdp.Client, frameID page.FrameID) (dom.Node, *eval.Runtime, error) {
	exec, err := eval.New(ctx, client, frameID)

	if err != nil {
		return dom.Node{}, nil, errors.Wrap(err, "create JS executor")
	}

	evalRes, err := exec.EvalRef(
		ctx,
		"return document",
	)

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
