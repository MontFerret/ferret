package dom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/net/html"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
)

var emptyExpires = time.Time{}

// parseAttrs is a helper function that parses a given interleaved array of node attribute names and values,
// and returns an object that represents attribute keys and values.
func parseAttrs(attrs []string) *values.Object {
	var attr values.String

	res := values.NewObject()

	for _, el := range attrs {
		el = strings.TrimSpace(el)
		str := values.NewString(el)

		if common.IsAttribute(el) {
			attr = str
			res.Set(str, values.EmptyString)
		} else {
			current, ok := res.Get(attr)

			if ok {
				if current.String() != "" {
					res.Set(attr, current.(values.String).Concat(values.SpaceString).Concat(str))
				} else {
					res.Set(attr, str)
				}
			}
		}
	}

	return res
}

func setInnerHTML(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, innerHTML values.String) error {
	var objID *runtime.RemoteObjectID

	if id.ObjectID != "" {
		objID = &id.ObjectID
	} else {
		repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

		if err != nil {
			return err
		}

		if repl.Object.ObjectID == nil {
			return errors.New("unable to resolve node")
		}

		objID = repl.Object.ObjectID
	}

	b, err := json.Marshal(innerHTML.String())

	if err != nil {
		return err
	}

	err = exec.EvalWithArguments(ctx, templates.SetInnerHTML(),
		runtime.CallArgument{
			ObjectID: objID,
		},
		runtime.CallArgument{
			Value: json.RawMessage(b),
		},
	)

	return err
}

func getInnerHTML(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
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

	repl, err := exec.EvalWithReturnValue(ctx, "return document.documentElement.innerHTML")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

func setInnerText(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, innerText values.String) error {
	var objID *runtime.RemoteObjectID

	if id.ObjectID != "" {
		objID = &id.ObjectID
	} else {
		repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.NodeID))

		if err != nil {
			return err
		}

		if repl.Object.ObjectID == nil {
			return errors.New("unable to resolve node")
		}

		objID = repl.Object.ObjectID
	}

	b, err := json.Marshal(innerText.String())

	if err != nil {
		return err
	}

	err = exec.EvalWithArguments(ctx, templates.SetInnerText(),
		runtime.CallArgument{
			ObjectID: objID,
		},
		runtime.CallArgument{
			Value: json.RawMessage(b),
		},
	)

	return err
}

func getInnerText(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
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

	repl, err := exec.EvalWithReturnValue(ctx, "return document.documentElement.innerText")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

func parseInnerText(innerHTML string) (values.String, error) {
	buff := bytes.NewBuffer([]byte(innerHTML))

	parsed, err := goquery.NewDocumentFromReader(buff)

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(parsed.Text()), nil
}

func createChildrenArray(nodes []dom.Node) []HTMLElementIdentity {
	children := make([]HTMLElementIdentity, len(nodes))

	for idx, child := range nodes {
		child := child
		children[idx] = HTMLElementIdentity{
			NodeID: child.NodeID,
		}
	}

	return children
}

func resolveFrame(ctx context.Context, client *cdp.Client, frameID page.FrameID) (dom.Node, runtime.ExecutionContextID, error) {
	worldRepl, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(frameID))

	if err != nil {
		return dom.Node{}, -1, err
	}

	evalRes, err := client.Runtime.Evaluate(
		ctx,
		runtime.NewEvaluateArgs(eval.PrepareEval("return document")).
			SetContextID(worldRepl.ExecutionContextID),
	)

	if err != nil {
		return dom.Node{}, -1, err
	}

	if evalRes.ExceptionDetails != nil {
		exception := *evalRes.ExceptionDetails

		return dom.Node{}, -1, errors.New(exception.Text)
	}

	if evalRes.Result.ObjectID == nil {
		return dom.Node{}, -1, errors.New("failed to resolve frame document")
	}

	req, err := client.DOM.RequestNode(ctx, dom.NewRequestNodeArgs(*evalRes.Result.ObjectID))

	if err != nil {
		return dom.Node{}, -1, err
	}

	if req.NodeID == 0 {
		return dom.Node{}, -1, errors.New("framed document is resolved with empty node id")
	}

	desc, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetNodeID(req.NodeID).
			SetDepth(1),
	)

	if err != nil {
		return dom.Node{}, -1, err
	}

	// Returned node, by some reason, does not contain the NodeID
	// So, we have to set it manually
	desc.Node.NodeID = req.NodeID

	return desc.Node, worldRepl.ExecutionContextID, nil
}
