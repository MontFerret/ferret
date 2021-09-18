package input

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	TypeParams struct {
		Text  string
		Clear bool
		Delay time.Duration
	}

	Manager struct {
		logger   zerolog.Logger
		client   *cdp.Client
		exec     *eval.Runtime
		keyboard *Keyboard
		mouse    *Mouse
	}
)

func New(
	logger zerolog.Logger,
	client *cdp.Client,
	exec *eval.Runtime,
	keyboard *Keyboard,
	mouse *Mouse,
) *Manager {
	logger = logging.WithName(logger.With(), "input_manager").Logger()

	return &Manager{
		logger,
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

func (m *Manager) ScrollTop(ctx context.Context, options drivers.ScrollOptions) error {
	m.logger.Trace().
		Str("behavior", options.Behavior.String()).
		Str("block", options.Block.String()).
		Str("inline", options.Inline.String()).
		Msg("scrolling to the top")

	if err := m.exec.Eval(ctx, templates.ScrollTop(options)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to scroll to the top")

		return err
	}

	m.logger.Trace().Msg("scrolled to the top")

	return nil
}

func (m *Manager) ScrollBottom(ctx context.Context, options drivers.ScrollOptions) error {
	m.logger.Trace().
		Str("behavior", options.Behavior.String()).
		Str("block", options.Block.String()).
		Str("inline", options.Inline.String()).
		Msg("scrolling to the bottom")

	if err := m.exec.Eval(ctx, templates.ScrollBottom(options)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to scroll to the bottom")

		return err
	}

	m.logger.Trace().Msg("scrolled to the bottom")

	return nil
}

func (m *Manager) ScrollIntoView(ctx context.Context, id runtime.RemoteObjectID, options drivers.ScrollOptions) error {
	m.logger.Trace().
		Str("object_id", string(id)).
		Str("behavior", options.Behavior.String()).
		Str("block", options.Block.String()).
		Str("inline", options.Inline.String()).
		Msg("scrolling to an element")

	if err := m.exec.Eval(ctx, templates.ScrollIntoView(id, options)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to scroll to an element")

		return err
	}

	m.logger.Trace().Msg("scrolled to an element")

	return nil
}

func (m *Manager) ScrollIntoViewBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector, options drivers.ScrollOptions) error {
	m.logger.Trace().
		Str("selector", selector.String()).
		Str("behavior", options.Behavior.String()).
		Str("block", options.Block.String()).
		Str("inline", options.Inline.String()).
		Msg("scrolling to an element by selector")

	if err := m.exec.Eval(ctx, templates.ScrollIntoViewBySelector(id, selector, options)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to scroll to an element by selector")

		return err
	}

	m.logger.Trace().Msg("scrolled to an element by selector")

	return nil
}

func (m *Manager) ScrollByXY(ctx context.Context, options drivers.ScrollOptions) error {
	m.logger.Trace().
		Float64("x", float64(options.Top)).
		Float64("y", float64(options.Left)).
		Str("behavior", options.Behavior.String()).
		Str("block", options.Block.String()).
		Str("inline", options.Inline.String()).
		Msg("scrolling to an element by given coordinates")

	if err := m.exec.Eval(ctx, templates.Scroll(options)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to scroll to an element by coordinates")

		return err
	}

	m.logger.Trace().Msg("scrolled to an element by given coordinates")

	return nil
}

func (m *Manager) Focus(ctx context.Context, objectID runtime.RemoteObjectID) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("focusing on an element")

	err := m.ScrollIntoView(ctx, objectID, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	if err := m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID)); err != nil {
		m.logger.Trace().Err(err).Msg("failed focusing on an element")

		return err
	}

	m.logger.Trace().Msg("focused on an element")

	return nil
}

func (m *Manager) FocusBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("focusing on an element by selector")

	err := m.ScrollIntoViewBySelector(ctx, id, selector, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	m.logger.Trace().Msg("resolving an element by selector")

	found, err := m.exec.EvalRef(ctx, templates.QuerySelector(id, selector))

	if err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed resolving an element by selector")

		return err
	}

	if found.ObjectID == nil {
		m.logger.Trace().
			Err(core.ErrNotFound).
			Msg("element not found by selector")

		return core.ErrNotFound
	}

	if err := m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(*found.ObjectID)); err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed focusing on an element by selector")

		return err
	}

	m.logger.Trace().Msg("focused on an element")

	return nil
}

func (m *Manager) Blur(ctx context.Context, objectID runtime.RemoteObjectID) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("removing focus from an element")

	if err := m.exec.Eval(ctx, templates.Blur(objectID)); err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed removing focus from an element")

		return err
	}

	m.logger.Trace().Msg("removed focus from an element")

	return nil
}

func (m *Manager) BlurBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("removing focus from an element by selector")

	if err := m.exec.Eval(ctx, templates.BlurBySelector(id, selector)); err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed removing focus from an element by selector")

		return err
	}

	m.logger.Trace().Msg("removed focus from an element by selector")

	return nil
}

func (m *Manager) MoveMouse(ctx context.Context, objectID runtime.RemoteObjectID) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("starting to move the mouse towards an element")

	if err := m.ScrollIntoView(ctx, objectID, drivers.ScrollOptions{}); err != nil {
		m.logger.Trace().Err(err).Msg("could not scroll into the object. failed to move the mouse")

		return err
	}

	m.logger.Trace().Msg("calculating clickable element points")

	q, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", q.X).Float64("y", q.Y).Msg("calculated clickable element points")

	if err := m.mouse.Move(ctx, q.X, q.Y); err != nil {
		m.logger.Trace().Err(err).Msg("failed to move the mouse")

		return err
	}

	m.logger.Trace().Msg("moved the mouse")

	return nil
}

func (m *Manager) MoveMouseBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("starting to move the mouse towards an element by selector")

	if err := m.ScrollIntoViewBySelector(ctx, id, selector, drivers.ScrollOptions{}); err != nil {
		return err
	}

	m.logger.Trace().Msg("looking up for an element by selector")

	found, err := m.exec.EvalRef(ctx, templates.QuerySelector(id, selector))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to find an element by selector")

		return err
	}

	if found.ObjectID == nil {
		m.logger.Trace().
			Err(core.ErrNotFound).
			Msg("element not found by selector")

		return core.ErrNotFound
	}

	m.logger.Trace().Str("object_id", string(*found.ObjectID)).Msg("calculating clickable element points")

	points, err := GetClickablePointByObjectID(ctx, m.client, *found.ObjectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

	if err := m.mouse.Move(ctx, points.X, points.Y); err != nil {
		m.logger.Trace().Err(err).Msg("failed to move the mouse")

		return err
	}

	m.logger.Trace().Msg("moved the mouse")

	return nil
}

func (m *Manager) MoveMouseByXY(ctx context.Context, xv, yv values.Float) error {
	x := float64(xv)
	y := float64(yv)

	m.logger.Trace().
		Float64("x", x).
		Float64("y", y).
		Msg("starting to move the mouse towards an element by given coordinates")

	if err := m.ScrollByXY(ctx, drivers.ScrollOptions{
		Top:  xv,
		Left: yv,
	}); err != nil {
		return err
	}

	if err := m.mouse.Move(ctx, x, y); err != nil {
		m.logger.Trace().Err(err).Msg("failed to move the mouse towards an element by given coordinates")

		return err
	}

	m.logger.Trace().Msg("moved the mouse")

	return nil
}

func (m *Manager) Click(ctx context.Context, objectID runtime.RemoteObjectID, count int) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("starting to click on an element")

	if err := m.ScrollIntoView(ctx, objectID, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	}); err != nil {
		return err
	}

	m.logger.Trace().Msg("calculating clickable element points")

	points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

	delay := time.Duration(drivers.DefaultMouseDelay) * time.Millisecond

	if err := m.mouse.ClickWithCount(ctx, points.X, points.Y, delay, count); err != nil {
		m.logger.Trace().
			Err(err).
			Msg("failed to click on an element")

		return err
	}

	m.logger.Trace().
		Err(err).
		Msg("clicked on an element")

	return nil
}

func (m *Manager) ClickBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector, count values.Int) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Int64("count", int64(count)).
		Msg("clicking on an element by selector")

	if err := m.ScrollIntoViewBySelector(ctx, id, selector, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	}); err != nil {
		return err
	}

	m.logger.Trace().Msg("looking up for an element by selector")

	found, err := m.exec.EvalRef(ctx, templates.QuerySelector(id, selector))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to find an element by selector")

		return err
	}

	if found.ObjectID == nil {
		m.logger.Trace().
			Err(core.ErrNotFound).
			Msg("element not found by selector")

		return core.ErrNotFound
	}

	m.logger.Trace().Str("object_id", string(*found.ObjectID)).Msg("calculating clickable element points")

	points, err := GetClickablePointByObjectID(ctx, m.client, *found.ObjectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

	delay := time.Duration(drivers.DefaultMouseDelay) * time.Millisecond

	if err := m.mouse.ClickWithCount(ctx, points.X, points.Y, delay, int(count)); err != nil {
		m.logger.Trace().Err(err).Msg("failed to click on an element")
		return nil
	}

	m.logger.Trace().Msg("clicked on an element")

	return nil
}

func (m *Manager) Type(ctx context.Context, objectID runtime.RemoteObjectID, params TypeParams) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("starting to type text")

	err := m.ScrollIntoView(ctx, objectID, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	m.logger.Trace().Msg("focusing on an element")

	if err := m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID)); err != nil {
		m.logger.Trace().Msg("failed to focus on an element")

		return err
	}

	m.logger.Trace().Bool("clear", params.Clear).Msg("is clearing text required?")

	if params.Clear {
		m.logger.Trace().Msg("calculating clickable element points")

		points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

		if err != nil {
			m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

			return err
		}

		m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

		if err := m.ClearByXY(ctx, points); err != nil {
			return err
		}
	}

	d := core.NumberLowerBoundary(float64(params.Delay))
	beforeTypeDelay := time.Duration(d)

	m.logger.Trace().Float64("delay", d).Msg("calculated pause delay")

	time.Sleep(beforeTypeDelay)

	m.logger.Trace().Msg("starting to type text")

	if err := m.keyboard.Type(ctx, params.Text, params.Delay); err != nil {
		m.logger.Trace().Err(err).Msg("failed to type text")

		return err
	}

	m.logger.Trace().Msg("typed text")

	return nil
}

func (m *Manager) TypeBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector, params TypeParams) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("starting to type text by selector")

	err := m.ScrollIntoViewBySelector(ctx, id, selector, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	m.logger.Trace().Msg("looking up for an element by selector")

	found, err := m.exec.EvalRef(ctx, templates.QuerySelector(id, selector))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to find an element by selector")

		return err
	}

	if found.ObjectID == nil {
		m.logger.Trace().
			Err(core.ErrNotFound).
			Msg("element not found by selector")

		return core.ErrNotFound
	}

	m.logger.Trace().Str("object_id", string(*found.ObjectID)).Msg("focusing on an element")

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(*found.ObjectID))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to focus on an element")

		return err
	}

	m.logger.Trace().Bool("clear", params.Clear).Msg("is clearing text required?")

	if params.Clear {
		m.logger.Trace().Msg("calculating clickable element points")

		points, err := GetClickablePointByObjectID(ctx, m.client, *found.ObjectID)

		if err != nil {
			m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

			return err
		}

		m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

		if err := m.ClearByXY(ctx, points); err != nil {
			return err
		}
	}

	d := core.NumberLowerBoundary(float64(params.Delay))
	beforeTypeDelay := time.Duration(d)

	m.logger.Trace().Float64("delay", d).Msg("calculated pause delay")

	time.Sleep(beforeTypeDelay)

	m.logger.Trace().Msg("starting to type text")

	if err := m.keyboard.Type(ctx, params.Text, params.Delay); err != nil {
		m.logger.Trace().Err(err).Msg("failed to type text")

		return err
	}

	m.logger.Trace().Msg("typed text")

	return nil
}

func (m *Manager) Clear(ctx context.Context, objectID runtime.RemoteObjectID) error {
	m.logger.Trace().
		Str("object_id", string(objectID)).
		Msg("starting to clear element")

	err := m.ScrollIntoView(ctx, objectID, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	m.logger.Trace().Msg("calculating clickable element points")

	points, err := GetClickablePointByObjectID(ctx, m.client, objectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")
	m.logger.Trace().Msg("focusing on an element")

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(objectID))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to focus on an element")

		return err
	}

	m.logger.Trace().Msg("clearing element")

	if err := m.ClearByXY(ctx, points); err != nil {
		m.logger.Trace().Err(err).Msg("failed to clear element")

		return err
	}

	m.logger.Trace().Msg("cleared element")

	return nil
}

func (m *Manager) ClearBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("starting to clear element by selector")

	err := m.ScrollIntoViewBySelector(ctx, id, selector, drivers.ScrollOptions{
		Behavior: drivers.ScrollBehaviorAuto,
		Block:    drivers.ScrollVerticalAlignmentCenter,
		Inline:   drivers.ScrollHorizontalAlignmentCenter,
	})

	if err != nil {
		return err
	}

	m.logger.Trace().Msg("looking up for an element by selector")

	found, err := m.exec.EvalRef(ctx, templates.QuerySelector(id, selector))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to find an element by selector")

		return err
	}

	if found.ObjectID == nil {
		m.logger.Trace().
			Err(core.ErrNotFound).
			Msg("element not found by selector")

		return core.ErrNotFound
	}

	m.logger.Trace().Str("object_id", string(*found.ObjectID)).Msg("calculating clickable element points")

	points, err := GetClickablePointByObjectID(ctx, m.client, *found.ObjectID)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed calculating clickable element points")

		return err
	}

	m.logger.Trace().Float64("x", points.X).Float64("y", points.Y).Msg("calculated clickable element points")

	m.logger.Trace().Msg("focusing on an element")

	err = m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetObjectID(*found.ObjectID))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to focus on an element")

		return err
	}

	m.logger.Trace().Msg("clearing element")

	if err := m.ClearByXY(ctx, points); err != nil {
		m.logger.Trace().Err(err).Msg("failed to clear element")

		return err
	}

	m.logger.Trace().Msg("cleared element")

	return nil
}

func (m *Manager) ClearByXY(ctx context.Context, points Quad) error {
	m.logger.Trace().
		Float64("x", points.X).
		Float64("y", points.Y).
		Msg("starting to clear element by coordinates")

	delay := time.Duration(drivers.DefaultMouseDelay) * time.Millisecond

	m.logger.Trace().Dur("delay", delay).Msg("clicking mouse button to select text")

	err := m.mouse.ClickWithCount(ctx, points.X, points.Y, delay, 3)

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to click mouse button")

		return err
	}

	delay = time.Duration(drivers.DefaultKeyboardDelay) * time.Millisecond

	m.logger.Trace().Dur("delay", delay).Msg("pressing 'Backspace'")

	if err := m.keyboard.Press(ctx, []string{"Backspace"}, 1, delay); err != nil {
		m.logger.Trace().Err(err).Msg("failed to press 'Backspace'")

		return err
	}

	return err
}

func (m *Manager) Press(ctx context.Context, keys []string, count int) error {
	delay := time.Duration(drivers.DefaultKeyboardDelay) * time.Millisecond

	m.logger.Trace().
		Strs("keys", keys).
		Int("count", count).
		Dur("delay", delay).
		Msg("pressing keyboard keys")

	if err := m.keyboard.Press(ctx, keys, count, delay); err != nil {
		m.logger.Trace().Err(err).Msg("failed to press keyboard keys")

		return err
	}

	return nil
}

func (m *Manager) PressBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector, keys []string, count int) error {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Strs("keys", keys).
		Int("count", count).
		Msg("starting to press keyboard keys by selector")

	if err := m.FocusBySelector(ctx, id, selector); err != nil {
		return err
	}

	return m.Press(ctx, keys, count)
}

func (m *Manager) Select(ctx context.Context, id runtime.RemoteObjectID, value *values.Array) (*values.Array, error) {
	m.logger.Trace().
		Str("object_id", string(id)).
		Msg("starting to select values")

	if err := m.Focus(ctx, id); err != nil {
		return values.NewArray(0), err
	}

	m.logger.Trace().Msg("selecting values")
	m.logger.Trace().Msg("evaluating a JS function")

	val, err := m.exec.EvalValue(ctx, templates.Select(id, value))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to evaluate a JS function")

		return nil, err
	}

	m.logger.Trace().Msg("validating JS result")

	arr, ok := val.(*values.Array)

	if !ok {
		m.logger.Trace().Err(err).Msg("JS result validation failed")

		return values.NewArray(0), core.ErrUnexpected
	}

	m.logger.Trace().Msg("selected values")

	return arr, nil
}

func (m *Manager) SelectBySelector(ctx context.Context, id runtime.RemoteObjectID, selector drivers.QuerySelector, value *values.Array) (*values.Array, error) {
	m.logger.Trace().
		Str("parent_object_id", string(id)).
		Str("selector", selector.String()).
		Msg("starting to select values by selector")

	if err := m.FocusBySelector(ctx, id, selector); err != nil {
		return values.NewArray(0), err
	}

	m.logger.Trace().Msg("selecting values")
	m.logger.Trace().Msg("evaluating a JS function")

	res, err := m.exec.EvalValue(ctx, templates.SelectBySelector(id, selector, value))

	if err != nil {
		m.logger.Trace().Err(err).Msg("failed to evaluate a JS function")

		return values.NewArray(0), err
	}

	m.logger.Trace().Msg("validating JS result")

	arr, ok := res.(*values.Array)

	if !ok {
		m.logger.Trace().Err(err).Msg("JS result validation failed")

		return values.NewArray(0), core.ErrUnexpected
	}

	m.logger.Trace().Msg("selected values")

	return arr, nil
}
