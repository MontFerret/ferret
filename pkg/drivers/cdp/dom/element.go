package dom

import (
	"context"
	runtime2 "github.com/MontFerret/ferret/pkg/runtime/internal"
	"hash/fnv"
	"strings"
	"time"

	"github.com/MontFerret/ferret/pkg/logging"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type HTMLElement struct {
	logger   zerolog.Logger
	client   *cdp.Client
	dom      *Manager
	input    *input.Manager
	eval     *eval.Runtime
	id       runtime.RemoteObjectID
	nodeType *common.LazyValue
	nodeName *common.LazyValue
}

func NewHTMLElement(
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	id runtime.RemoteObjectID,
) *HTMLElement {
	el := new(HTMLElement)
	el.logger = logging.WithName(logger.With(), "dom_element").
		Str("object_id", string(id)).
		Logger()
	el.client = client
	el.dom = domManager
	el.input = input
	el.eval = exec
	el.id = id
	el.nodeType = common.NewLazyValue(func(ctx context.Context) (core.Value, error) {
		return el.eval.EvalValue(ctx, templates.GetNodeType(el.id))
	})
	el.nodeName = common.NewLazyValue(func(ctx context.Context) (core.Value, error) {
		return el.eval.EvalValue(ctx, templates.GetNodeName(el.id))
	})

	return el
}

func (el *HTMLElement) RemoteID() runtime.RemoteObjectID {
	return el.id
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
	other, ok := other.(drivers.HTMLElement)

	if !ok {
		return drivers.CompareTypes(el.Type(), core.Reflect(other))
	}

	return int64(strings.Compare(el.String(), other.String()))
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
	return core.None
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) Get(ctx context.Context, key string) (core.Value, error) {
	return common.GetInElement(ctx, key, el)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInElement(ctx, path, el, value)
}

func (el *HTMLElement) GetValue(ctx context.Context) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetValue(el.id))
}

func (el *HTMLElement) SetValue(ctx context.Context, value core.Value) error {
	return el.eval.Eval(ctx, templates.SetValue(el.id, value))
}

func (el *HTMLElement) GetNodeType(ctx context.Context) (core.Int, error) {
	out, err := el.nodeType.Read(ctx)

	if err != nil {
		return core.ZeroInt, err
	}

	return internal.ToInt(out), nil
}

func (el *HTMLElement) GetNodeName(ctx context.Context) (core.String, error) {
	out, err := el.nodeName.Read(ctx)

	if err != nil {
		return core.EmptyString, err
	}

	return internal.ToString(out), nil
}

func (el *HTMLElement) Length() core.Int {
	value, err := el.eval.EvalValue(context.Background(), templates.GetChildrenCount(el.id))

	if err != nil {
		el.logError(err)

		return 0
	}

	return internal.ToInt(value)
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*runtime2.Object, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetStyles(el.id))

	if err != nil {
		return runtime2.NewObject(), err
	}

	return internal.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name core.String) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetStyle(el.id, name))
}

func (el *HTMLElement) SetStyles(ctx context.Context, styles *runtime2.Object) error {
	return el.eval.Eval(ctx, templates.SetStyles(el.id, styles))
}

func (el *HTMLElement) SetStyle(ctx context.Context, name, value core.String) error {
	return el.eval.Eval(ctx, templates.SetStyle(el.id, name, value))
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, names ...core.String) error {
	return el.eval.Eval(ctx, templates.RemoveStyles(el.id, names))
}

func (el *HTMLElement) GetAttributes(ctx context.Context) (*runtime2.Object, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetAttributes(el.id))

	if err != nil {
		return runtime2.NewObject(), err
	}

	return internal.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name core.String) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetAttribute(el.id, name))
}

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *runtime2.Object) error {
	return el.eval.Eval(ctx, templates.SetAttributes(el.id, attrs))
}

func (el *HTMLElement) SetAttribute(ctx context.Context, name, value core.String) error {
	return el.eval.Eval(ctx, templates.SetAttribute(el.id, name, value))
}

func (el *HTMLElement) RemoveAttribute(ctx context.Context, names ...core.String) error {
	return el.eval.Eval(ctx, templates.RemoveAttributes(el.id, names))
}

func (el *HTMLElement) GetChildNodes(ctx context.Context) (*runtime2.Array, error) {
	return el.eval.EvalElements(ctx, templates.GetChildren(el.id))
}

func (el *HTMLElement) GetChildNode(ctx context.Context, idx core.Int) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetChildByIndex(el.id, idx))
}

func (el *HTMLElement) GetParentElement(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetParent(el.id))
}

func (el *HTMLElement) GetPreviousElementSibling(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetPreviousElementSibling(el.id))
}

func (el *HTMLElement) GetNextElementSibling(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetNextElementSibling(el.id))
}

func (el *HTMLElement) QuerySelector(ctx context.Context, selector drivers.QuerySelector) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.QuerySelector(el.id, selector))
}

func (el *HTMLElement) QuerySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*runtime2.Array, error) {
	return el.eval.EvalElements(ctx, templates.QuerySelectorAll(el.id, selector))
}

func (el *HTMLElement) XPath(ctx context.Context, expression core.String) (result core.Value, err error) {
	return el.eval.EvalValue(ctx, templates.XPath(el.id, expression))
}

func (el *HTMLElement) GetInnerText(ctx context.Context) (core.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerText(el.id))

	if err != nil {
		return core.EmptyString, err
	}

	return internal.ToString(out), nil
}

func (el *HTMLElement) SetInnerText(ctx context.Context, innerText core.String) error {
	return el.eval.Eval(
		ctx,
		templates.SetInnerText(el.id, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector) (core.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerTextBySelector(el.id, selector))

	if err != nil {
		return core.EmptyString, err
	}

	return internal.ToString(out), nil
}

func (el *HTMLElement) SetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector, innerText core.String) error {
	return el.eval.Eval(
		ctx,
		templates.SetInnerTextBySelector(el.id, selector, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*runtime2.Array, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerTextBySelectorAll(el.id, selector))

	if err != nil {
		return runtime2.EmptyArray(), err
	}

	return internal.ToArray(ctx, out), nil
}

func (el *HTMLElement) GetInnerHTML(ctx context.Context) (core.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTML(el.id))

	if err != nil {
		return core.EmptyString, err
	}

	return internal.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTML(ctx context.Context, innerHTML core.String) error {
	return el.eval.Eval(ctx, templates.SetInnerHTML(el.id, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector) (core.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTMLBySelector(el.id, selector))

	if err != nil {
		return core.EmptyString, err
	}

	return internal.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector, innerHTML core.String) error {
	return el.eval.Eval(ctx, templates.SetInnerHTMLBySelector(el.id, selector, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*runtime2.Array, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTMLBySelectorAll(el.id, selector))

	if err != nil {
		return runtime2.EmptyArray(), err
	}

	return internal.ToArray(ctx, out), nil
}

func (el *HTMLElement) CountBySelector(ctx context.Context, selector drivers.QuerySelector) (core.Int, error) {
	out, err := el.eval.EvalValue(ctx, templates.CountBySelector(el.id, selector))

	if err != nil {
		return core.ZeroInt, err
	}

	return internal.ToInt(out), nil
}

func (el *HTMLElement) ExistsBySelector(ctx context.Context, selector drivers.QuerySelector) (core.Boolean, error) {
	out, err := el.eval.EvalValue(ctx, templates.ExistsBySelector(el.id, selector))

	if err != nil {
		return core.False, err
	}

	return internal.ToBoolean(out), nil
}

func (el *HTMLElement) WaitForElement(ctx context.Context, selector drivers.QuerySelector, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForElement(el.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForElementAll(ctx context.Context, selector drivers.QuerySelector, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForElementAll(el.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClass(ctx context.Context, class core.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForClass(el.id, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClassBySelector(ctx context.Context, selector drivers.QuerySelector, class core.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForClassBySelector(el.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForClassBySelectorAll(ctx context.Context, selector drivers.QuerySelector, class core.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForClassBySelectorAll(el.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttribute(
	ctx context.Context,
	name core.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForAttribute(el.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttributeBySelector(ctx context.Context, selector drivers.QuerySelector, name core.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForAttributeBySelector(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttributeBySelectorAll(ctx context.Context, selector drivers.QuerySelector, name core.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForAttributeBySelectorAll(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyle(ctx context.Context, name core.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForStyle(el.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyleBySelector(ctx context.Context, selector drivers.QuerySelector, name core.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForStyleBySelector(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyleBySelectorAll(ctx context.Context, selector drivers.QuerySelector, name core.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		el.eval,
		templates.WaitForStyleBySelectorAll(el.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) Click(ctx context.Context, count core.Int) error {
	return el.input.Click(ctx, el.id, int(count))
}

func (el *HTMLElement) ClickBySelector(ctx context.Context, selector drivers.QuerySelector, count core.Int) error {
	return el.input.ClickBySelector(ctx, el.id, selector, count)
}

func (el *HTMLElement) ClickBySelectorAll(ctx context.Context, selector drivers.QuerySelector, count core.Int) error {
	elements, err := el.QuerySelectorAll(ctx, selector)

	if err != nil {
		return err
	}

	elements.ForEach(func(value core.Value, idx int) bool {
		found := value.(*HTMLElement)

		if e := found.Click(ctx, count); e != nil {
			err = e
			return false
		}

		return true
	})

	return err
}

func (el *HTMLElement) Input(ctx context.Context, value core.Value, delay core.Int) error {
	name, err := el.GetNodeName(ctx)

	if err != nil {
		return err
	}

	if strings.ToLower(string(name)) != "input" {
		return core.Error(core.ErrInvalidOperation, "element is not an <input> element.")
	}

	return el.input.Type(ctx, el.id, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) InputBySelector(ctx context.Context, selector drivers.QuerySelector, value core.Value, delay core.Int) error {
	return el.input.TypeBySelector(ctx, el.id, selector, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) Press(ctx context.Context, keys []core.String, count core.Int) error {
	return el.input.Press(ctx, internal.UnwrapStrings(keys), int(count))
}

func (el *HTMLElement) PressBySelector(ctx context.Context, selector drivers.QuerySelector, keys []core.String, count core.Int) error {
	return el.input.PressBySelector(ctx, el.id, selector, internal.UnwrapStrings(keys), int(count))
}

func (el *HTMLElement) Clear(ctx context.Context) error {
	return el.input.Clear(ctx, el.id)
}

func (el *HTMLElement) ClearBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return el.input.ClearBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Select(ctx context.Context, value *runtime2.Array) (*runtime2.Array, error) {
	return el.input.Select(ctx, el.id, value)
}

func (el *HTMLElement) SelectBySelector(ctx context.Context, selector drivers.QuerySelector, value *runtime2.Array) (*runtime2.Array, error) {
	return el.input.SelectBySelector(ctx, el.id, selector, value)
}

func (el *HTMLElement) ScrollIntoView(ctx context.Context, options drivers.ScrollOptions) error {
	return el.input.ScrollIntoView(ctx, el.id, options)
}

func (el *HTMLElement) Focus(ctx context.Context) error {
	return el.input.Focus(ctx, el.id)
}

func (el *HTMLElement) FocusBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return el.input.FocusBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Blur(ctx context.Context) error {
	return el.input.Blur(ctx, el.id)
}

func (el *HTMLElement) BlurBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return el.input.BlurBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) Hover(ctx context.Context) error {
	return el.input.MoveMouse(ctx, el.id)
}

func (el *HTMLElement) HoverBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return el.input.MoveMouseBySelector(ctx, el.id, selector)
}

func (el *HTMLElement) logError(err error) *zerolog.Event {
	return el.logger.
		Error().
		Timestamp().
		Str("objectID", string(el.id)).
		Err(err)
}
