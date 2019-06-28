package input

import (
	"context"
	"fmt"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/gofrs/uuid"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
)

type Manager struct {
	client   *cdp.Client
	exec     *eval.ExecutionContext
	keyboard *Keyboard
	mouse    *Mouse
}

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

func (m *Manager) Scroll(ctx context.Context, x, y values.Float) error {
	return m.exec.Eval(ctx, fmt.Sprintf(`
		window.scrollBy({
  			top: %s,
  			left: %s,
  			behavior: 'instant'
		});
	`,
		eval.ParamFloat(float64(x)),
		eval.ParamFloat(float64(y)),
	))
}

func (m *Manager) ScrollIntoViewBySelector(ctx context.Context, selector values.String) error {
	return m.exec.Eval(ctx, fmt.Sprintf(`
		var el = document.querySelector(%s);

		if (el == null) {
			throw new Error("element not found");
		}

		el.scrollIntoView({
    		behavior: 'instant'
  		});

		return true;
	`, eval.ParamString(selector.String()),
	))
}

func (m *Manager) ScrollIntoViewByNodeID(ctx context.Context, nodeID dom.NodeID) error {
	var attrID = "data-ferret-scroll"

	id, err := uuid.NewV4()

	if err != nil {
		return err
	}

	err = m.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(nodeID, attrID, id.String()))

	if err != nil {
		return err
	}

	err = m.exec.Eval(
		ctx,
		fmt.Sprintf(`
			var el = document.querySelector('[%s="%s"]');
			if (el == null) {
				throw new Error('element not found');
			}
			
			el.scrollIntoView({
    			behavior: 'instant',
				inline: 'center',
				block: 'center'
  			});
		`,
			attrID,
			id.String(),
		))

	if err != nil {
		return err
	}

	err = m.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(nodeID, attrID))

	return err
}

func (m *Manager) ScrollTop(ctx context.Context) error {
	return m.exec.Eval(ctx, `
		window.scrollTo({
			left: 0,
			top: 0,
    		behavior: 'instant'
  		});
	`)
}

func (m *Manager) ScrollBottom(ctx context.Context) error {
	return m.exec.Eval(ctx, `
		window.scrollTo({
			left: 0,
			top: window.document.body.scrollHeight,
    		behavior: 'instant'
  		});
	`)
}

func (m *Manager) MoveMouseBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	return m.MoveMouseByNodeID(ctx, found.NodeID)
}

func (m *Manager) MoveMouseByNodeID(ctx context.Context, nodeID dom.NodeID) error {
	err := m.ScrollIntoViewByNodeID(ctx, nodeID)

	if err != nil {
		return err
	}

	q, err := GetClickablePointByNodeID(ctx, m.client, nodeID)

	if err != nil {
		return err
	}

	return m.mouse.Move(ctx, q.X, q.Y)
}

func (m *Manager) MoveMouse(ctx context.Context, x, y values.Float) error {
	return m.mouse.Move(ctx, float64(x), float64(y))
}

func (m *Manager) ClickBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String) error {
	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	return m.ClickByNodeID(ctx, found.NodeID)
}

func (m *Manager) ClickByNodeID(ctx context.Context, nodeID dom.NodeID) error {
	if err := m.ScrollIntoViewByNodeID(ctx, nodeID); err != nil {
		return err
	}

	points, err := GetClickablePointByNodeID(ctx, m.client, nodeID)

	if err != nil {
		return err
	}

	if err := m.mouse.Click(ctx, points.X, points.Y, 50); err != nil {
		return nil
	}

	return nil
}

func (m *Manager) TypeBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String, text core.Value, delay values.Int) error {
	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return err
	}

	return m.TypeByNodeID(ctx, found.NodeID, text, delay)
}

func (m *Manager) TypeByNodeID(ctx context.Context, nodeID dom.NodeID, text core.Value, delay values.Int) error {
	if err := m.ScrollIntoViewByNodeID(ctx, nodeID); err != nil {
		return err
	}

	if err := m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(nodeID)); err != nil {
		return err
	}

	_, min := core.NumberBoundaries(float64(delay))
	beforeTypeDelay := time.Duration(min)

	time.Sleep(beforeTypeDelay * time.Millisecond)

	return m.keyboard.Type(ctx, text.String(), int(delay))
}

func (m *Manager) SelectBySelector(ctx context.Context, parentNodeID dom.NodeID, selector values.String, value *values.Array) (*values.Array, error) {
	found, err := m.client.DOM.QuerySelector(ctx, dom.NewQuerySelectorArgs(parentNodeID, selector.String()))

	if err != nil {
		return nil, err
	}

	return m.SelectByNodeID(ctx, found.NodeID, value)
}

func (m *Manager) SelectByNodeID(ctx context.Context, nodeID dom.NodeID, value *values.Array) (*values.Array, error) {
	if err := m.ScrollIntoViewByNodeID(ctx, nodeID); err != nil {
		return nil, err
	}

	if err := m.client.DOM.Focus(ctx, dom.NewFocusArgs().SetNodeID(nodeID)); err != nil {
		return nil, err
	}

	var attrID = "data-ferret-select"

	id, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	err = m.client.DOM.SetAttributeValue(ctx, dom.NewSetAttributeValueArgs(nodeID, attrID, id.String()))

	if err != nil {
		return nil, err
	}

	res, err := m.exec.EvalWithValue(
		ctx,
		fmt.Sprintf(`
			var el = document.querySelector('[%s="%s"]');
			if (el == null) {
				return [];
			}
			var values = %s;
			if (el.nodeName.toLowerCase() !== 'select') {
				throw new Error('element is not a <select> element.');
			}
			var options = Array.from(el.options);
      		el.value = undefined;
			for (var option of options) {
        		option.selected = values.includes(option.value);
        	
				if (option.selected && !el.multiple) {
          			break;
				}
      		}
      		el.dispatchEvent(new Event('input', { 'bubbles': true }));
      		el.dispatchEvent(new Event('change', { 'bubbles': true }));
      		
			return options.filter(option => option.selected).map(option => option.value);
		`,
			attrID,
			id.String(),
			value.String(),
		),
	)

	if err != nil {
		return nil, err
	}

	err = m.client.DOM.RemoveAttribute(ctx, dom.NewRemoveAttributeArgs(nodeID, attrID))

	if err != nil {
		return nil, err
	}

	arr, ok := res.(*values.Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(types.Array, res.Type())
}
