package cdp

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/net/html"
	"math"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"golang.org/x/sync/errgroup"
)

var emptyExpires = time.Time{}

type (
	batchFunc = func() error

	Quad struct {
		X float64
		Y float64
	}
)

func runBatch(funcs ...batchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}

func fromProtocolQuad(quad dom.Quad) []Quad {
	return []Quad{
		{
			X: quad[0],
			Y: quad[1],
		},
		{
			X: quad[2],
			Y: quad[3],
		},
		{
			X: quad[4],
			Y: quad[5],
		},
		{
			X: quad[6],
			Y: quad[7],
		},
	}
}

func computeQuadArea(quads []Quad) float64 {
	var area float64

	for i := range quads {
		p1 := quads[i]
		p2 := quads[(i+1)%len(quads)]
		area += (p1.X*p2.Y - p2.X*p1.Y) / 2
	}

	return math.Abs(area)
}

func intersectQuadWithViewport(quad []Quad, width, height float64) []Quad {
	quads := make([]Quad, 0, len(quad))

	for _, point := range quad {
		quads = append(quads, Quad{
			X: math.Min(math.Max(point.X, 0), width),
			Y: math.Min(math.Max(point.Y, 0), height),
		})
	}

	return quads
}

func getClickablePoint(ctx context.Context, client *cdp.Client, id HTMLElementIdentity) (Quad, error) {
	qargs := dom.NewGetContentQuadsArgs()

	switch {
	case id.objectID != "":
		qargs.SetObjectID(id.objectID)
	case id.backendID != 0:
		qargs.SetBackendNodeID(id.backendID)
	default:
		qargs.SetNodeID(id.nodeID)
	}

	contentQuadsReply, err := client.DOM.GetContentQuads(ctx, qargs)

	if err != nil {
		return Quad{}, err
	}

	if contentQuadsReply.Quads == nil || len(contentQuadsReply.Quads) == 0 {
		return Quad{}, errors.New("node is either not visible or not an HTMLElement")
	}

	layoutMetricsReply, err := client.Page.GetLayoutMetrics(ctx)

	if err != nil {
		return Quad{}, err
	}

	clientWidth := layoutMetricsReply.LayoutViewport.ClientWidth
	clientHeight := layoutMetricsReply.LayoutViewport.ClientHeight

	quads := make([][]Quad, 0, len(contentQuadsReply.Quads))

	for _, q := range contentQuadsReply.Quads {
		quad := intersectQuadWithViewport(fromProtocolQuad(q), float64(clientWidth), float64(clientHeight))

		if computeQuadArea(quad) > 1 {
			quads = append(quads, quad)
		}
	}

	if len(quads) == 0 {
		return Quad{}, errors.New("node is either not visible or not an HTMLElement")
	}

	// Return the middle point of the first quad.
	quad := quads[0]
	var x float64
	var y float64

	for _, q := range quad {
		x += q.X
		y += q.Y
	}

	return Quad{
		X: x / 4,
		Y: y / 4,
	}, nil
}

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

func loadInnerHTML(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
	// not a document
	if nodeType != html.DocumentNode {
		var objID runtime.RemoteObjectID

		switch {
		case id.objectID != "":
			objID = id.objectID
		case id.backendID > 0:
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetBackendNodeID(id.backendID))

			if err != nil {
				return "", err
			}

			if repl.Object.ObjectID == nil {
				return "", errors.New("unable to resolve node")
			}

			objID = *repl.Object.ObjectID
		default:
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.nodeID))

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

	repl, err := exec.EvalWithReturn(ctx, "return document.documentElement.innerHTML")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

func loadInnerHTMLByNodeID(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, nodeID dom.NodeID) (values.String, error) {
	node, err := client.DOM.DescribeNode(ctx, dom.NewDescribeNodeArgs().SetNodeID(nodeID))

	if err != nil {
		return values.EmptyString, err
	}

	return loadInnerHTML(ctx, client, exec, HTMLElementIdentity{
		nodeID: nodeID,
	}, common.ToHTMLType(node.Node.NodeType))
}

func loadInnerText(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, id HTMLElementIdentity, nodeType html.NodeType) (values.String, error) {
	// not a document
	if nodeType != html.DocumentNode {
		var objID runtime.RemoteObjectID

		switch {
		case id.objectID != "":
			objID = id.objectID
		case id.backendID > 0:
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetBackendNodeID(id.backendID))

			if err != nil {
				return "", err
			}

			if repl.Object.ObjectID == nil {
				return "", errors.New("unable to resolve node")
			}

			objID = *repl.Object.ObjectID
		default:
			repl, err := client.DOM.ResolveNode(ctx, dom.NewResolveNodeArgs().SetNodeID(id.nodeID))

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

	repl, err := exec.EvalWithReturn(ctx, "return document.documentElement.innerText")

	if err != nil {
		return "", err
	}

	return values.NewString(repl.String()), nil
}

//func loadInnerTextByNodeID(ctx context.Context, client *cdp.Client, exec *eval.ExecutionContext, nodeID dom.NodeID) (values.String, error) {
//	node, err := client.DOM.DescribeNode(ctx, dom.NewDescribeNodeArgs().SetNodeID(nodeID))
//
//	if err != nil {
//		return values.EmptyString, err
//	}
//
//	return loadInnerText(ctx, client, exec, HTMLElementIdentity{
//		nodeID: nodeID,
//	}, common.ToHTMLType(node.Node.NodeType))
//}

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
			nodeID:    child.NodeID,
			backendID: child.BackendNodeID,
		}
	}

	return children
}

func fromDriverCookie(url string, cookie drivers.HTTPCookie) network.CookieParam {
	sameSite := network.CookieSameSiteNotSet

	switch cookie.SameSite {
	case drivers.SameSiteLaxMode:
		sameSite = network.CookieSameSiteLax
	case drivers.SameSiteStrictMode:
		sameSite = network.CookieSameSiteStrict
	}

	if cookie.Expires == emptyExpires {
		cookie.Expires = time.Now().Add(time.Duration(24) + time.Hour)
	}

	normalizedURL := normalizeCookieURL(url)

	return network.CookieParam{
		URL:      &normalizedURL,
		Name:     cookie.Name,
		Value:    cookie.Value,
		Secure:   &cookie.Secure,
		Path:     &cookie.Path,
		Domain:   &cookie.Domain,
		HTTPOnly: &cookie.HTTPOnly,
		SameSite: sameSite,
		Expires:  network.TimeSinceEpoch(cookie.Expires.Unix()),
	}
}

func fromDriverCookieDelete(url string, cookie drivers.HTTPCookie) *network.DeleteCookiesArgs {
	normalizedURL := normalizeCookieURL(url)

	return &network.DeleteCookiesArgs{
		URL:    &normalizedURL,
		Name:   cookie.Name,
		Path:   &cookie.Path,
		Domain: &cookie.Domain,
	}
}

func toDriverCookie(c network.Cookie) drivers.HTTPCookie {
	sameSite := drivers.SameSiteDefaultMode

	switch c.SameSite {
	case network.CookieSameSiteLax:
		sameSite = drivers.SameSiteLaxMode
	case network.CookieSameSiteStrict:
		sameSite = drivers.SameSiteStrictMode
	}

	return drivers.HTTPCookie{
		Name:     c.Name,
		Value:    c.Value,
		Path:     c.Path,
		Domain:   c.Domain,
		Expires:  time.Unix(int64(c.Expires), 0),
		SameSite: sameSite,
		Secure:   c.Secure,
		HTTPOnly: c.HTTPOnly,
	}
}

func normalizeCookieURL(url string) string {
	const httpPrefix = "http://"
	const httpsPrefix = "https://"

	if strings.HasPrefix(url, httpPrefix) || strings.HasPrefix(url, httpsPrefix) {
		return url
	}

	return httpPrefix + url
}

func randomDuration(delay values.Int) time.Duration {
	max, min := core.NumberBoundaries(float64(int64(delay)))
	value := core.Random(max, min)

	return time.Duration(int64(value))
}

func resolveFrame(ctx context.Context, client *cdp.Client, frame page.Frame) (dom.Node, runtime.ExecutionContextID, error) {
	worldRepl, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(frame.ID))

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
