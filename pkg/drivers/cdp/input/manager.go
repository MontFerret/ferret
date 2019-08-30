package input

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	TypeParams struct {
		Text core.Value
		Clear values.Boolean
		Delay values.Int
	}

	Manager struct {
		client   *cdp.Client
		exec     *eval.ExecutionContext
		keyboard *Keyboard
		mouse    *Mouse
	}
)

func NewManager(
	client *cdp.Client,
	exec *eval.ExecutionContext,
	keyboard *Keyboard,
	mouse *Mouse,
) *Manager {
	return &Manager{
		client,
		exec,
		keyboard,
		mouse,
	}
}

func (m *Manager) Keyboard() *Keyboard {
	return m.keyboard
}

func (m *Manager) Mouse() *Mouse {
	return m.mouse
}

func (m *Manager) ScrollTop(ctx context.Context) error {
	return m.exec.Eval(ctx, templates.ScrollTop())
}

func (m *Manager) ScrollBottom(ctx context.Context) error {
	return m.exec.Eval(ctx, templates.ScrollBottom())
}

func (m *Manager) ScrollIntoView(ctx context.Context, objectID runtime.RemoteObjectID) error {
	return m.exec.EvalWithArguments(
		ctx,
		templates.ScrollIntoView(),
		runtime.CallArgument{
			ObjectID: &objectID,
		},
	)
}

func (m *Manager) ScrollIntoViewBySelector(ctx context.Context, selector values.String) error {
	return m.exec.Eval(ctx, templates.ScrollIntoViewBySelector(selector.String()))
}

func (m *Manager) ScrollByXY(ctx context.Context, x, y values.Float) error {
	return m.exec.Eval(
		ctx,
		templates.Scroll(eval.ParamFloat(float64(x)), eval.ParamFloat(float64(y))),
	)
}

func (m *Manager) Focus(ctx context.Context, objectID runtime.RemoteObjectID) error {
	err := m.ScrollIntoView(ctx, objectID)

	if err != nil {
		return err
	}

	return m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID))
}

func (m *Manager) FocusBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	err := m.ScrollIntoViewBySelector(ctx, selector)

	if err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return nil
	}

	return m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(found.NodeID))
}

func (m *Manager) MoveMouse(ctx context.Context, objectID runtime.RemoteObjectID) error {
	if err := m.ScrollIntoView(ctx, objectID); err != nil {
		return err
	}

	q, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		return err
	}

	return m.mouse.Move(ctx, q.X, q.Y)
}

func (m *Manager) MoveMouseBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	if err := m.ScrollIntoViewBySelector(ctx, selector); err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	q, err := GetClickablePointByNodeID(ctx, m.client, found.NodeID)

	if err != nil {
		return err
	}

	return m.mouse.Move(ctx, q.X, q.Y)
}

func (m *Manager) MoveMouseByXY(ctx context.Context, x, y values.Float) error {
	if err := m.ScrollByXY(ctx, x, y); err != nil {
		return err
	}

	return m.mouse.Move(ctx, float64(x), float64(y))
}

func (m *Manager) Click(ctx context.Context, objectID runtime.RemoteObjectID) error {
	if err := m.ScrollIntoView(ctx, objectID); err != nil {
		return err
	}

	points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		return err
	}

	if err := m.mouse.Click(ctx, points.X, points.Y, 50); err != nil {
		return nil
	}

	return nil
}

func (m *Manager) ClickBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	if err := m.ScrollIntoViewBySelector(ctx, selector); err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	points, err := GetClickablePointByNodeID(ctx, m.client, found.NodeID)

	if err != nil {
		return err
	}

	if err := m.mouse.Click(ctx, points.X, points.Y, 50); err != nil {
		return nil
	}

	return nil
}

func (m *Manager) ClickBySelectorAll(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	if err := m.ScrollIntoViewBySelector(ctx, selector); err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelectorAll(ctx, dom.NewQuerySelectorAllArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	for _, nodeID := range found.NodeIDs {
		_, min := core.NumberBoundaries(100)
		beforeTypeDelay := time.Duration(min)

		time.Sleep(beforeTypeDelay * time.Millisecond)

		points, err := GetClickablePointByNodeID(ctx, m.client, nodeID)

		if err != nil {
			return err
		}

		if err := m.mouse.Click(ctx, points.X, points.Y, 50); err != nil {
			return nil
		}
	}

	return nil
}

func (m *Manager) Type(ctx context.Context, objectID runtime.RemoteObjectID, params TypeParams) error {
	err := m.ScrollIntoView(ctx, objectID)

	if err != nil {
		return err
	}

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID))

	if err != nil {
		return err
	}

	if params.Clear {
		points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

		if err != nil {
			return err
		}

		if err := m.ClearByXY(ctx, points); err != nil {
			return err
		}
	}

	_, min := core.NumberBoundaries(float64(params.Delay))
	beforeTypeDelay := time.Duration(min)

	time.Sleep(beforeTypeDelay * time.Millisecond)

	return m.keyboard.Type(ctx, params.Text.String(), int(params.Delay))
}

func (m *Manager) TypeBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String, params TypeParams) error {
	err := m.ScrollIntoViewBySelector(ctx, selector)

	if err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(found.NodeID))

	if err != nil {
		return err
	}

	if params.Clear {
		points, err := GetClickablePointByNodeID(ctx, m.client, found.NodeID)

		if err != nil {
			return err
		}

		if err := m.ClearByXY(ctx, points); err != nil {
			return err
		}
	}

	_, min := core.NumberBoundaries(float64(params.Delay))
	beforeTypeDelay := time.Duration(min)

	time.Sleep(beforeTypeDelay * time.Millisecond)

	return m.keyboard.Type(ctx, params.Text.String(), int(params.Delay))
}

func (m *Manager) Clear(ctx context.Context, objectID runtime.RemoteObjectID) error {
	err := m.ScrollIntoView(ctx, objectID)

	if err != nil {
		return err
	}

	points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		return err
	}

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID))

	if err != nil {
		return err
	}

	return m.ClearByXY(ctx, points)
}

func (m *Manager) ClearBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	err := m.ScrollIntoViewBySelector(ctx, selector)

	if err != nil {
		return err
	}

	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	points, err := GetClickablePointByNodeID(ctx, m.client, found.NodeID)

	if err != nil {
		return err
	}

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(found.NodeID))

	if err != nil {
		return err
	}

	return m.ClearByXY(ctx, points)
}

func (m *Manager) ClearByXY(ctx context.Context, points Quad) error {
	err := m.mouse.ClickWithCount(ctx, points.X, points.Y, 3, 5)

	if err != nil {
		return err
	}

	return m.keyboard.Press(ctx, "Backspace")
}

func (m *Manager) Select(ctx context.Context, objectID runtime.RemoteObjectID, value *values.Array) (*values.Array, error) {
	if err := m.Focus(ctx, objectID); err != nil {
		return values.NewArray(0), err
	}

	val, err := m.exec.EvalWithArgumentsAndReturnValue(ctx, templates.Select(value.String()), runtime.CallArgument{
		ObjectID: &objectID,
	})

	if err != nil {
		return nil, err
	}

	arr, ok := val.(*values.Array)

	if !ok {
		return values.NewArray(0), core.ErrUnexpected
	}

	return arr, nil
}

func (m *Manager) SelectBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String, value *values.Array) (*values.Array, error) {
	if err := m.FocusBySelector(ctx, parentNodeID, selector); err != nil {
		return values.NewArray(0), err
	}

	res, err := m.exec.EvalWithReturnValue(ctx, templates.SelectBySelector(selector.String(), value.String()))

	if err != nil {
		return values.NewArray(0), err
	}

	arr, ok := res.(*values.Array)

	if !ok {
		return values.NewArray(0), core.ErrUnexpected
	}

	return arr, nil
}