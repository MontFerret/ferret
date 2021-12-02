package input

import (
	"context"
	"math"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/utils"
)

type Quad struct {
	X float64
	Y float64
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

func getClickablePoint(ctx context.Context, client *cdp.Client, qargs *dom.GetContentQuadsArgs) (Quad, error) {
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

	clientWidth, clientHeight := utils.GetLayoutViewportWH(layoutMetricsReply)
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

func GetClickablePointByObjectID(ctx context.Context, client *cdp.Client, objectID runtime.RemoteObjectID) (Quad, error) {
	return getClickablePoint(ctx, client, dom.NewGetContentQuadsArgs().SetObjectID(objectID))
}
