package dom

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/wI2L/jettison"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLElement struct {
	logger   zerolog.Logger
	client   *cdp.Client
	dom      *Manager
	input    *input.Manager
	exec     *eval.Runtime
	id       runtime.RemoteObjectID
	nodeType html.NodeType
	nodeName values.String
}

func LoadHTMLElement(
	ctx context.Context,
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	nodeID dom.NodeID,
) (*HTMLElement, error) {
	if client == nil {
		return nil, core.Error(core.ErrMissedArgument, "client")
	}

	// getting a remote object that represents the current DOM Node
	args := dom.NewResolveNodeArgs().SetNodeID(nodeID).SetExecutionContextID(exec.ContextID())

	obj, err := client.DOM.ResolveNode(ctx, args)

	if err != nil {
		return nil, err
	}

	return ResolveHTMLElement(
		ctx,
		logger,
		client,
		domManager,
		input,
		exec,
		obj.Object,
	)
}

func ResolveHTMLElement(
	ctx context.Context,
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	ref runtime.RemoteObject,
) (*HTMLElement, error) {
	if ref.ObjectID == nil {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("element %s", ref.Value))
	}

	id := *ref.ObjectID

	node, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetObjectID(id).
			SetDepth(0),
	)

	if err != nil {
		return nil, core.Error(err, string(id))
	}

	return NewHTMLElement(
		logger,
		client,
		domManager,
		input,
		exec,
		id,
		node.Node.NodeType,
		node.Node.NodeName,
	), nil
}

func NewHTMLElement(
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	id runtime.RemoteObjectID,
	nodeType int,
	nodeName string,
) *HTMLElement {
	el := new(HTMLElement)
	el.logger = logging.
		WithName(logger.With(), "dom_element").
		Str("object_id", string(id)).
		Str("node_name", nodeName).
		Logger()
	el.client = client
	el.dom = domManager
	el.input = input
	el.exec = exec
	el.id = id
	el.nodeType = common.ToHTMLType(nodeType)
	el.nodeName = values.NewString(nodeName)

	return el
}

func (el *HTMLElement) Close() error {
	return nil
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(el.String(), jettison.NoHTMLEscaping())
}

func (el *HTMLElement) String() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(drivers.DefaultWaitTimeout)*time.Millisecond)
	defer cancel()

	res, err := el.GetInnerHTML(ctx)

	if err != nil {
		el.logError(errors.Wrap(err, "HTMLElement.String"))

		return ""
	}

	return res.String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		return int64(strings.Compare(el.String(), other.String()))
	default:
		return drivers.Compare(el.Type(), other.Type())
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el
}

func (el *HTMLElement) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(el.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(el.id))

	return h.Sum64()
}

func (el *HTMLElement) Copy() core.Value {
	return values.None
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInElement(ctx, el, path)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInElement(ctx, el, path, value)
}

func (el *HTMLElement) GetValue(ctx context.Context) (core.Value, error) {
	return el.exec.EvalValue(ctx, templates.GetValue(el.id))
}

func (el *HTMLElement) SetValue(ctx context.Context, value core.Value) error {
	return el.exec.Eval(ctx, templates.SetValue(el.id, value))
}

func (el *HTMLElement) GetNodeType() values.Int {
	return values.NewInt(common.FromHTMLType(el.nodeType))
}

func (el *HTMLElement) GetNodeName() values.String {
	return el.nodeName
}

func (el *HTMLElement) Length() values.Int {
	value, err := el.exec.EvalValue(context.Background(), templates.GetChildrenCount(el.id))

	if err != nil {
		el.logError(err)

		return 0
	}

	return values.ToInt(value)
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*values.Object, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetStyles(el.id))

	if err != nil {
		return values.NewObject(), err
	}

	return values.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name values.String) (core.Value, error) {
	return el.exec.EvalValue(ctx, templates.GetStyle(el.id, name))
}

func (el *HTMLElement) SetStyles(ctx context.Context, styles *values.Object) error {
	return el.exec.Eval(ctx, templates.SetStyles(el.id, styles))
}

func (el *HTMLElement) SetStyle(ctx context.Context, name, value values.String) error {
	return el.exec.Eval(ctx, templates.SetStyle(el.id, name, value))
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, names ...values.String) error {
	return el.exec.Eval(ctx, templates.RemoveStyles(el.id, names))
}

func (el *HTMLElement) GetAttributes(ctx context.Context) (*values.Object, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetAttributes(el.id))

	if err != nil {
		return values.NewObject(), err
	}

	return values.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name values.String) (core.Value, error) {
	return el.exec.EvalValue(ctx, templates.GetAttribute(el.id, name))
}

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *values.Object) error {
	return el.exec.Eval(ctx, templates.SetAttributes(el.id, attrs))
}

func (el *HTMLElement) SetAttribute(ctx context.Context, name, value values.String) error {
	return el.exec.Eval(ctx, templates.SetAttribute(el.id, name, value))
}

func (el *HTMLElement) RemoveAttribute(ctx context.Context, names ...values.String) error {
	return el.exec.Eval(ctx, templates.RemoveAttributes(el.id, names))
}

func (el *HTMLElement) GetChildNodes(ctx context.Context) (*values.Array, error) {
	out, err := el.evalTo(ctx, templates.GetChildren(el.id))

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, out), nil
}

func (el *HTMLElement) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	return el.evalToElement(ctx, templates.GetChildByIndex(el.id, idx))
}

func (el *HTMLElement) GetParentElement(ctx context.Context) (core.Value, error) {
	return el.evalToElement(ctx, templates.GetParent(el.id))
}

func (el *HTMLElement) GetPreviousElementSibling(ctx context.Context) (core.Value, error) {
	return el.evalToElement(ctx, templates.GetPreviousElementSibling(el.id))
}

func (el *HTMLElement) GetNextElementSibling(ctx context.Context) (core.Value, error) {
	return el.evalToElement(ctx, templates.GetNextElementSibling(el.id))
}

func (el *HTMLElement) QuerySelector(ctx context.Context, selector values.String) (core.Value, error) {
	return el.evalToElement(ctx, templates.QuerySelector(el.id, selector))
}

func (el *HTMLElement) QuerySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	out, err := el.exec.EvalRef(ctx, templates.QuerySelectorAll(el.id, selector))

	if err != nil {
		return values.EmptyArray(), err
	}

	res, err := el.fromEvalRef(ctx, out)

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, res), nil
}

func (el *HTMLElement) XPath(ctx context.Context, expression values.String) (result core.Value, err error) {
	out, err := el.exec.EvalRef(ctx, templates.XPath(el.id, expression))

	if err != nil {
		return values.None, err
	}

	return el.fromEvalRef(ctx, out)
}

func (el *HTMLElement) GetInnerText(ctx context.Context) (values.String, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerText(el.id))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerText(ctx context.Context, innerText values.String) error {
	return el.exec.Eval(
		ctx,
		templates.SetInnerText(el.id, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelector(ctx context.Context, selector values.String) (values.String, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerTextBySelector(el.id, selector))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerTextBySelector(ctx context.Context, selector, innerText values.String) error {
	return el.exec.Eval(
		ctx,
		templates.SetInnerTextBySelector(el.id, selector, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerTextBySelectorAll(el.id, selector))

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, out), nil
}

func (el *HTMLElement) GetInnerHTML(ctx context.Context) (values.String, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerHTML(el.id))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTML(ctx context.Context, innerHTML values.String) error {
	return el.exec.Eval(ctx, templates.SetInnerHTML(el.id, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelector(ctx context.Context, selector values.String) (values.String, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerHTMLBySelector(el.id, selector))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTMLBySelector(ctx context.Context, selector, innerHTML values.String) error {
	return el.exec.Eval(ctx, templates.SetInnerHTMLBySelector(el.id, selector, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	out, err := el.exec.EvalValue(ctx, templates.GetInnerHTMLBySelectorAll(el.id, selector))

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, out), nil
}

func (el *HTMLElement) CountBySelector(ctx context.Context, selector values.String) (values.Int, error) {
	out, err := el.exec.EvalValue(ctx, templates.CountBySelector(el.id, selector))

	if err != nil {
		return values.ZeroInt, err
	}

	return values.ToInt(out), nil
}

func (el *HTMLElement) ExistsBySelector(ctx context.Context, selector values.String) (values.Boolean, error) {
	out, err := el.exec.EvalValue(ctx, templates.ExistsBySelector(el.id, selector))

	if err != nil {
		return values.False, err
	}

	return values.ToBoolean(out), nil
}

func (el *HTMLElement) WaitForElement(ctx context.Context, selector values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForElement(el.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForElementAll(ctx context.Context, selector values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForElementAll(el.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClass(ctx context.Context, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForClass(el.id, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClassBySelector(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForClassBySelector(el.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClassBySelectorAll(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForClassBySelectorAll(el.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttribute(
	ctx context.Context,
	name,
	value values.String,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForAttribute(el.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttributeBySelector(ctx context.Context, selector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForAttributeBySelector(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttributeBySelectorAll(ctx context.Context, selector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForAttributeBySelectorAll(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyle(ctx context.Context, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForStyle(el.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyleBySelector(ctx context.Context, selector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForStyleBySelector(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyleBySelectorAll(ctx context.Context, selector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.exec,
		templates.WaitForStyleBySelectorAll(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) Click(ctx context.Context, count values.Int) error {
	return el.input.Click(ctx, el.id, int(count))
}

func (el *HTMLElement) ClickBySelector(ctx context.Context, selector values.String, count values.Int) error {
	return el.input.ClickBySelector(ctx, el.id, selector, count)
}

func (el *HTMLElement) ClickBySelectorAll(ctx context.Context, selector values.String, count values.Int) error {
	return el.input.ClickBySelectorAll(ctx, el.id, selector, count)
}

func (el *HTMLElement) Input(ctx context.Context, value core.Value, delay values.Int) error {
	if el.GetNodeName() != "INPUT" {
		return core.Error(core.ErrInvalidOperation, "element is not an <input> element.")
	}

	return el.input.Type(ctx, el.id, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) error {
	return el.input.TypeBySelector(ctx, el.id, selector, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) Press(ctx context.Context, keys []values.String, count values.Int) error {
	return el.input.Press(ctx, values.UnwrapStrings(keys), int(count))
}

func (el *HTMLElement) PressBySelector(ctx context.Context, selector values.String, keys []values.String, count values.Int) error {
	return el.input.PressBySelector(ctx, el.id, selector, values.UnwrapStrings(keys), int(count))
}

func (el *HTMLElement) Clear(ctx context.Context) error {
	return el.input.Clear(ctx, el.id)
}

func (el *HTMLElement) ClearBySelector(ctx context.Context, selector values.String) error {
	return el.input.ClearBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Select(ctx context.Context, value *values.Array) (*values.Array, error) {
	return el.input.Select(ctx, el.id, value)
}

func (el *HTMLElement) SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error) {
	return el.input.SelectBySelector(ctx, el.id, selector, value)
}

func (el *HTMLElement) ScrollIntoView(ctx context.Context, options drivers.ScrollOptions) error {
	return el.input.ScrollIntoView(ctx, el.id, options)
}

func (el *HTMLElement) Focus(ctx context.Context) error {
	return el.input.Focus(ctx, el.id)
}

func (el *HTMLElement) FocusBySelector(ctx context.Context, selector values.String) error {
	return el.input.FocusBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Blur(ctx context.Context) error {
	return el.input.Blur(ctx, el.id)
}

func (el *HTMLElement) BlurBySelector(ctx context.Context, selector values.String) error {
	return el.input.BlurBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Hover(ctx context.Context) error {
	return el.input.MoveMouse(ctx, el.id)
}

func (el *HTMLElement) HoverBySelector(ctx context.Context, selector values.String) error {
	return el.input.MoveMouseBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) evalTo(ctx context.Context, fn *eval.Function) (core.Value, error) {
	out, err := el.exec.EvalRef(ctx, fn)

	if err != nil {
		return values.None, err
	}

	return el.fromEvalRef(ctx, out)
}

func (el *HTMLElement) evalToElement(ctx context.Context, fn *eval.Function) (core.Value, error) {
	obj, err := el.exec.EvalRef(ctx, fn)

	if err != nil {
		return values.None, err
	}

	if obj.Type != "object" || obj.ObjectID == nil {
		return values.None, nil
	}

	return ResolveHTMLElement(
		ctx,
		el.logger,
		el.client,
		el.dom,
		el.input,
		el.exec,
		obj,
	)
}

func (el *HTMLElement) fromEvalRef(ctx context.Context, out runtime.RemoteObject) (core.Value, error) {
	typeName := out.Type
	var className string
	var subtype string

	if out.ClassName != nil {
		className = *out.ClassName
	}

	if out.Subtype != nil {
		subtype = *out.Subtype
	}

	if subtype == "null" || subtype == "undefined" {
		return values.None, nil
	}

	// checking whether it's actually an array
	if typeName == "object" {
		switch className {
		case "Array":
			typeName = "array"
		case "HTMLCollection":
			typeName = "HTMLCollection"
		default:
			break
		}
	}

	switch typeName {
	case "array", "HTMLCollection":
		if out.ObjectID == nil {
			return values.None, nil
		}

		props, err := el.client.Runtime.GetProperties(ctx, runtime.NewGetPropertiesArgs(*out.ObjectID).SetOwnProperties(true))

		if err != nil {
			return values.None, err
		}

		if props.ExceptionDetails != nil {
			exception := *props.ExceptionDetails

			return values.None, errors.New(exception.Text)
		}

		result := values.NewArray(len(props.Result))

		for _, descr := range props.Result {
			if !descr.Enumerable {
				continue
			}

			if descr.Value == nil {
				continue
			}

			// it's not a Node, it's an attr value
			if descr.Value.ObjectID == nil {
				var value interface{}

				if err := json.Unmarshal(descr.Value.Value, &value); err != nil {
					return values.None, err
				}

				result.Push(values.Parse(value))

				continue
			}

			el, err := ResolveHTMLElement(
				ctx,
				el.logger,
				el.client,
				el.dom,
				el.input,
				el.exec,
				*descr.Value,
			)

			if err != nil {
				return values.None, err
			}

			result.Push(el)
		}

		return result, nil
	case "object":
		if out.ObjectID == nil {
			var value interface{}

			if err := json.Unmarshal(out.Value, &value); err != nil {
				return values.None, err
			}

			return values.Parse(value), nil
		}

		return ResolveHTMLElement(
			ctx,
			el.logger,
			el.client,
			el.dom,
			el.input,
			el.exec,
			out,
		)
	case "string", "number", "boolean":
		return eval.Unmarshal(out)
	default:
		return values.None, nil
	}
}

func (el *HTMLElement) logError(err error) *zerolog.Event {
	return el.logger.
		Error().
		Timestamp().
		Str("objectID", string(el.id)).
		Err(err)
}
