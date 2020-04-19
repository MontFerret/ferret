package dom

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

var emptyNodeID = dom.NodeID(0)

type (
	HTMLElementIdentity struct {
		NodeID   dom.NodeID
		ObjectID runtime.RemoteObjectID
	}

	elementListeners struct {
		pageReload            events.ListenerID
		attrModified          events.ListenerID
		attrRemoved           events.ListenerID
		childNodeCountUpdated events.ListenerID
		childNodeInserted     events.ListenerID
		childNodeRemoved      events.ListenerID
	}

	HTMLElement struct {
		mu             sync.Mutex
		logger         *zerolog.Logger
		client         *cdp.Client
		dom            *Manager
		input          *input.Manager
		exec           *eval.ExecutionContext
		connected      values.Boolean
		id             HTMLElementIdentity
		nodeType       html.NodeType
		nodeName       values.String
		innerHTML      *common.LazyValue
		innerText      *common.LazyValue
		attributes     *common.LazyValue
		style          *common.LazyValue
		children       []HTMLElementIdentity
		loadedChildren *common.LazyValue
		listeners      *elementListeners
	}
)

func LoadHTMLElement(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.ExecutionContext,
	nodeID dom.NodeID,
) (*HTMLElement, error) {
	if client == nil {
		return nil, core.Error(core.ErrMissedArgument, "client")
	}

	// getting a remote object that represents the current DOM Node
	args := dom.NewResolveNodeArgs().SetNodeID(nodeID).SetExecutionContextID(exec.ID())

	obj, err := client.DOM.ResolveNode(ctx, args)

	if err != nil {
		return nil, err
	}

	if obj.Object.ObjectID == nil {
		return nil, core.Error(core.ErrNotFound, fmt.Sprintf("element %d", nodeID))
	}

	return LoadHTMLElementWithID(
		ctx,
		logger,
		client,
		domManager,
		input,
		exec,
		HTMLElementIdentity{
			NodeID:   nodeID,
			ObjectID: *obj.Object.ObjectID,
		},
	)
}

func LoadHTMLElementWithID(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.ExecutionContext,
	id HTMLElementIdentity,
) (*HTMLElement, error) {
	node, err := client.DOM.DescribeNode(
		ctx,
		dom.
			NewDescribeNodeArgs().
			SetObjectID(id.ObjectID).
			SetDepth(1),
	)

	if err != nil {
		return nil, core.Error(err, strconv.Itoa(int(id.NodeID)))
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
		createChildrenArray(node.Node.Children),
	), nil
}

func NewHTMLElement(
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.ExecutionContext,
	id HTMLElementIdentity,
	nodeType int,
	nodeName string,
	children []HTMLElementIdentity,
) *HTMLElement {
	el := new(HTMLElement)
	el.logger = logger
	el.client = client
	el.dom = domManager
	el.input = input
	el.exec = exec
	el.connected = values.True
	el.id = id
	el.nodeType = common.ToHTMLType(nodeType)
	el.nodeName = values.NewString(nodeName)
	el.innerHTML = common.NewLazyValue(el.loadInnerHTML)
	el.innerText = common.NewLazyValue(el.loadInnerText)
	el.attributes = common.NewLazyValue(el.loadAttrs)
	el.style = common.NewLazyValue(el.parseStyle)
	el.loadedChildren = common.NewLazyValue(el.loadChildren)
	el.children = children
	el.listeners = &elementListeners{
		pageReload:            domManager.AddDocumentUpdatedListener(el.handlePageReload),
		attrModified:          domManager.AddAttrModifiedListener(el.handleAttrModified),
		attrRemoved:           domManager.AddAttrRemovedListener(el.handleAttrRemoved),
		childNodeCountUpdated: domManager.AddChildNodeCountUpdatedListener(el.handleChildrenCountChanged),
		childNodeInserted:     domManager.AddChildNodeInsertedListener(el.handleChildInserted),
		childNodeRemoved:      domManager.AddChildNodeRemovedListener(el.handleChildRemoved),
	}

	return el
}

func (el *HTMLElement) Close() error {
	el.mu.Lock()
	defer el.mu.Unlock()

	// already closed
	if !el.connected {
		return nil
	}

	el.connected = values.False

	el.dom.RemoveReloadListener(el.listeners.pageReload)
	el.dom.RemoveAttrModifiedListener(el.listeners.attrModified)
	el.dom.RemoveAttrRemovedListener(el.listeners.attrRemoved)
	el.dom.RemoveChildNodeCountUpdatedListener(el.listeners.childNodeCountUpdated)
	el.dom.RemoveChildNodeInsertedListener(el.listeners.childNodeInserted)
	el.dom.RemoveChildNodeRemovedListener(el.listeners.childNodeRemoved)

	return nil
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(el.String())
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
	h.Write([]byte(strconv.Itoa(int(el.id.NodeID))))

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
	if el.IsDetached() {
		return values.None, drivers.ErrDetached
	}

	return el.exec.ReadProperty(ctx, el.id.ObjectID, "value")
}

func (el *HTMLElement) SetValue(ctx context.Context, value core.Value) error {
	if el.IsDetached() {
		return drivers.ErrDetached
	}

	return el.client.DOM.SetNodeValue(ctx, dom.NewSetNodeValueArgs(el.id.NodeID, value.String()))
}

func (el *HTMLElement) GetNodeType() values.Int {
	return values.NewInt(common.FromHTMLType(el.nodeType))
}

func (el *HTMLElement) GetNodeName() values.String {
	return el.nodeName
}

func (el *HTMLElement) Length() values.Int {
	return values.NewInt(len(el.children))
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*values.Object, error) {
	val, err := el.style.Read(ctx)

	if err != nil {
		return values.NewObject(), err
	}

	if val == values.None {
		return values.NewObject(), nil
	}

	styles := val.(*values.Object)

	return styles.Copy().(*values.Object), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name values.String) (core.Value, error) {
	styles, err := el.style.Read(ctx)

	if err != nil {
		return values.None, err
	}

	val, found := styles.(*values.Object).Get(name)

	if !found {
		return values.None, nil
	}

	return val, nil
}

func (el *HTMLElement) SetStyles(ctx context.Context, styles *values.Object) error {
	if styles == nil {
		return nil
	}

	val, err := el.style.Read(ctx)

	if err != nil {
		return err
	}

	currentStyles := val.(*values.Object)

	styles.ForEach(func(value core.Value, key string) bool {
		currentStyles.Set(values.NewString(key), value)

		return true
	})

	str := common.SerializeStyles(ctx, currentStyles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) SetStyle(ctx context.Context, name values.String, value core.Value) error {
	val, err := el.style.Read(ctx)

	if err != nil {
		return err
	}

	styles := val.(*values.Object)
	styles.Set(name, value)

	str := common.SerializeStyles(ctx, styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, names ...values.String) error {
	if len(names) == 0 {
		return nil
	}

	val, err := el.style.Read(ctx)

	if err != nil {
		return err
	}

	styles := val.(*values.Object)

	for _, name := range names {
		styles.Remove(name)
	}

	str := common.SerializeStyles(ctx, styles)

	return el.SetAttribute(ctx, "style", str)
}

func (el *HTMLElement) GetAttributes(ctx context.Context) (*values.Object, error) {
	val, err := el.attributes.Read(ctx)

	if err != nil {
		return values.NewObject(), err
	}

	attrs := val.(*values.Object)

	// returning shallow copy
	return attrs.Copy().(*values.Object), nil
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name values.String) (core.Value, error) {
	attrs, err := el.attributes.Read(ctx)

	if err != nil {
		return values.None, err
	}

	val, found := attrs.(*values.Object).Get(name)

	if !found {
		return values.None, nil
	}

	return val, nil
}

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *values.Object) error {
	var err error

	attrs.ForEach(func(value core.Value, key string) bool {
		err = el.SetAttribute(ctx, values.NewString(key), values.NewString(value.String()))

		return err == nil
	})

	return err
}

func (el *HTMLElement) SetAttribute(ctx context.Context, name, value values.String) error {
	return el.client.DOM.SetAttributeValue(
		ctx,
		dom.NewSetAttributeValueArgs(el.id.NodeID, string(name), string(value)),
	)
}

func (el *HTMLElement) RemoveAttribute(ctx context.Context, names ...values.String) error {
	for _, name := range names {
		err := el.client.DOM.RemoveAttribute(
			ctx,
			dom.NewRemoveAttributeArgs(el.id.NodeID, name.String()),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (el *HTMLElement) GetChildNodes(ctx context.Context) (*values.Array, error) {
	val, err := el.loadedChildren.Read(ctx)

	if err != nil {
		return values.NewArray(0), err
	}

	return val.Copy().(*values.Array), nil
}

func (el *HTMLElement) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	val, err := el.loadedChildren.Read(ctx)

	if err != nil {
		return values.None, err
	}

	return val.(*values.Array).Get(idx), nil
}

func (el *HTMLElement) QuerySelector(ctx context.Context, selector values.String) (core.Value, error) {
	if el.IsDetached() {
		return values.None, drivers.ErrDetached
	}

	selectorArgs := dom.NewQuerySelectorArgs(el.id.NodeID, selector.String())
	found, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		return values.None, err
	}

	if found.NodeID == emptyNodeID {
		return values.None, nil
	}

	res, err := LoadHTMLElement(
		ctx,
		el.logger,
		el.client,
		el.dom,
		el.input,
		el.exec,
		found.NodeID,
	)

	if err != nil {
		return values.None, nil
	}

	return res, nil
}

func (el *HTMLElement) QuerySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	if el.IsDetached() {
		return values.NewArray(0), drivers.ErrDetached
	}

	selectorArgs := dom.NewQuerySelectorAllArgs(el.id.NodeID, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		return values.NewArray(0), err
	}

	arr := values.NewArray(len(res.NodeIDs))

	for _, id := range res.NodeIDs {
		if id == emptyNodeID {
			el.logError(err).
				Str("selector", selector.String()).
				Msg("failed to find a node by selector. returned 0 NodeID")

			continue
		}

		childEl, err := LoadHTMLElement(
			ctx,
			el.logger,
			el.client,
			el.dom,
			el.input,
			el.exec,
			id,
		)

		if err != nil {
			// close elements that are already loaded, but won't be used because of the error
			if arr.Length() > 0 {
				arr.ForEach(func(e core.Value, _ int) bool {
					e.(*HTMLElement).Close()

					return true
				})
			}

			return values.NewArray(0), err
		}

		arr.Push(childEl)
	}

	return arr, nil
}

func (el *HTMLElement) XPath(ctx context.Context, expression values.String) (result core.Value, err error) {
	exp, err := expression.MarshalJSON()

	if err != nil {
		return values.None, err
	}

	out, err := el.exec.EvalWithArgumentsAndReturnReference(ctx, templates.XPath(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: json.RawMessage(exp),
		},
	)

	if err != nil {
		return values.None, err
	}

	typeName := out.Type

	// checking whether it's actually an array
	if typeName == "object" {
		if out.ClassName != nil && *out.ClassName == "Array" {
			typeName = "array"
		}
	}

	switch typeName {
	case "string", "number", "boolean":
		return eval.Unmarshal(&out)
	case "array":
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

		defer func() {
			if err != nil {
				result.ForEach(func(value core.Value, idx int) bool {
					el, ok := value.(*HTMLElement)

					if ok {
						el.Close()
					}

					return true
				})
			}
		}()

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

			repl, err := el.client.DOM.RequestNode(ctx, dom.NewRequestNodeArgs(*descr.Value.ObjectID))

			if err != nil {
				return values.None, err
			}

			el, err := LoadHTMLElementWithID(
				ctx,
				el.logger,
				el.client,
				el.dom,
				el.input,
				el.exec,
				HTMLElementIdentity{
					NodeID:   repl.NodeID,
					ObjectID: *descr.Value.ObjectID,
				},
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

		repl, err := el.client.DOM.RequestNode(ctx, dom.NewRequestNodeArgs(*out.ObjectID))

		if err != nil {
			return values.None, err
		}

		return LoadHTMLElementWithID(
			ctx,
			el.logger,
			el.client,
			el.dom,
			el.input,
			el.exec,
			HTMLElementIdentity{
				NodeID:   repl.NodeID,
				ObjectID: *out.ObjectID,
			},
		)
	default:
		return values.None, nil
	}
}

func (el *HTMLElement) GetInnerText(ctx context.Context) (values.String, error) {
	val, err := el.innerText.Read(ctx)

	if err != nil {
		return values.EmptyString, err
	}

	if val == values.None {
		return values.EmptyString, nil
	}

	return val.(values.String), nil
}

func (el *HTMLElement) SetInnerText(ctx context.Context, innerText values.String) error {
	if el.IsDetached() {
		return drivers.ErrDetached
	}

	el.innerText.Reset()

	return setInnerText(ctx, el.client, el.exec, el.id, innerText)
}

func (el *HTMLElement) GetInnerTextBySelector(ctx context.Context, selector values.String) (values.String, error) {
	if el.IsDetached() {
		return values.EmptyString, drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return values.EmptyString, err
	}

	out, err := el.exec.EvalWithArgumentsAndReturnValue(
		ctx,
		templates.GetInnerTextBySelector(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
	)

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(out.String()), nil
}

func (el *HTMLElement) SetInnerTextBySelector(ctx context.Context, selector, innerText values.String) error {
	if el.IsDetached() {
		return drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return err
	}

	val, err := innerText.MarshalJSON()

	if err != nil {
		return err
	}

	return el.exec.EvalWithArguments(
		ctx,
		templates.SetInnerTextBySelector(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
		runtime.CallArgument{
			Value: val,
		},
	)
}

func (el *HTMLElement) GetInnerTextBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	if el.IsDetached() {
		return values.NewArray(0), drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return nil, err
	}

	out, err := el.exec.EvalWithArgumentsAndReturnValue(
		ctx,
		templates.GetInnerTextBySelectorAll(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
	)

	if err != nil {
		return values.NewArray(0), err
	}

	arr, ok := out.(*values.Array)

	if !ok {
		return values.NewArray(0), errors.New("unexpected output")
	}

	return arr, nil
}

func (el *HTMLElement) GetInnerHTML(ctx context.Context) (values.String, error) {
	val, err := el.innerHTML.Read(ctx)

	if err != nil {
		return values.EmptyString, err
	}

	if val == values.None {
		return values.EmptyString, nil
	}

	return val.(values.String), nil
}

func (el *HTMLElement) SetInnerHTML(ctx context.Context, innerHTML values.String) error {
	if el.IsDetached() {
		return drivers.ErrDetached
	}

	el.innerHTML.Reset()

	return setInnerHTML(ctx, el.client, el.exec, el.id, innerHTML)
}

func (el *HTMLElement) GetInnerHTMLBySelector(ctx context.Context, selector values.String) (values.String, error) {
	if el.IsDetached() {
		return values.EmptyString, drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return values.EmptyString, err
	}

	out, err := el.exec.EvalWithArgumentsAndReturnValue(
		ctx,
		templates.GetInnerHTMLBySelector(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
	)

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(out.String()), nil
}

func (el *HTMLElement) SetInnerHTMLBySelector(ctx context.Context, selector, innerHTML values.String) error {
	if el.IsDetached() {
		return drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return err
	}

	val, err := innerHTML.MarshalJSON()

	if err != nil {
		return err
	}

	return el.exec.EvalWithArguments(
		ctx,
		templates.SetInnerHTMLBySelector(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
		runtime.CallArgument{
			Value: val,
		},
	)
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	if el.IsDetached() {
		return values.NewArray(0), drivers.ErrDetached
	}

	sel, err := selector.MarshalJSON()

	if err != nil {
		return values.NewArray(0), err
	}

	out, err := el.exec.EvalWithArgumentsAndReturnValue(
		ctx,
		templates.GetInnerHTMLBySelectorAll(),
		runtime.CallArgument{
			ObjectID: &el.id.ObjectID,
		},
		runtime.CallArgument{
			Value: sel,
		},
	)

	if err != nil {
		return values.NewArray(0), err
	}

	arr, ok := out.(*values.Array)

	if !ok {
		return values.NewArray(0), errors.New("unexpected output")
	}

	return arr, nil
}

func (el *HTMLElement) CountBySelector(ctx context.Context, selector values.String) (values.Int, error) {
	if el.IsDetached() {
		return values.ZeroInt, drivers.ErrDetached
	}

	selectorArgs := dom.NewQuerySelectorAllArgs(el.id.NodeID, selector.String())
	res, err := el.client.DOM.QuerySelectorAll(ctx, selectorArgs)

	if err != nil {
		return values.ZeroInt, err
	}

	return values.NewInt(len(res.NodeIDs)), nil
}

func (el *HTMLElement) ExistsBySelector(ctx context.Context, selector values.String) (values.Boolean, error) {
	if el.IsDetached() {
		return values.False, drivers.ErrDetached
	}

	// TODO: Can we use RemoteObjectID or BackendID instead of NodeId?
	selectorArgs := dom.NewQuerySelectorArgs(el.id.NodeID, selector.String())
	res, err := el.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		return values.False, err
	}

	if res.NodeID == 0 {
		return values.False, nil
	}

	return values.True, nil
}

func (el *HTMLElement) WaitForClass(ctx context.Context, class values.String, when drivers.WaitEvent) error {
	task := events.NewWaitTask(
		func(ctx2 context.Context) (core.Value, error) {
			current, err := el.GetAttribute(ctx2, "class")

			if err != nil {
				return values.None, nil
			}

			if current.Type() != types.String {
				return values.None, nil
			}

			str := current.(values.String)
			classStr := string(class)
			classes := strings.Split(string(str), " ")

			if when != drivers.WaitEventAbsence {
				for _, c := range classes {
					if c == classStr {
						// The value does not really matter if it's not None
						// None indicates that operation needs to be repeated
						return values.True, nil
					}
				}
			} else {
				var found values.Boolean

				for _, c := range classes {
					if c == classStr {
						found = values.True
						break
					}
				}

				if found == values.False {
					// The value does not really matter if it's not None
					// None indicates that operation needs to be repeated
					return values.False, nil
				}
			}

			return values.None, nil
		},
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForAttribute(
	ctx context.Context,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewValueWaitTask(when, value, func(ctx context.Context) (core.Value, error) {
		return el.GetAttribute(ctx, name)
	}, events.DefaultPolling)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) WaitForStyle(ctx context.Context, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewValueWaitTask(when, value, func(ctx context.Context) (core.Value, error) {
		return el.GetStyle(ctx, name)
	}, events.DefaultPolling)

	_, err := task.Run(ctx)

	return err
}

func (el *HTMLElement) Click(ctx context.Context, count values.Int) error {
	return el.input.Click(ctx, el.id.ObjectID, int(count))
}

func (el *HTMLElement) ClickBySelector(ctx context.Context, selector values.String, count values.Int) error {
	return el.input.ClickBySelector(ctx, el.id.NodeID, selector.String(), int(count))
}

func (el *HTMLElement) ClickBySelectorAll(ctx context.Context, selector values.String, count values.Int) error {
	return el.input.ClickBySelectorAll(ctx, el.id.NodeID, selector.String(), int(count))
}

func (el *HTMLElement) Input(ctx context.Context, value core.Value, delay values.Int) error {
	if el.GetNodeName() != "INPUT" {
		return core.Error(core.ErrInvalidOperation, "element is not an <input> element.")
	}

	return el.input.Type(ctx, el.id.ObjectID, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) error {
	return el.input.TypeBySelector(ctx, el.id.NodeID, selector.String(), input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (el *HTMLElement) Clear(ctx context.Context) error {
	return el.input.Clear(ctx, el.id.ObjectID)
}

func (el *HTMLElement) ClearBySelector(ctx context.Context, selector values.String) error {
	return el.input.ClearBySelector(ctx, el.id.NodeID, selector.String())
}

func (el *HTMLElement) Select(ctx context.Context, value *values.Array) (*values.Array, error) {
	return el.input.Select(ctx, el.id.ObjectID, value)
}

func (el *HTMLElement) SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error) {
	return el.input.SelectBySelector(ctx, el.id.NodeID, selector.String(), value)
}

func (el *HTMLElement) ScrollIntoView(ctx context.Context, options drivers.ScrollOptions) error {
	return el.input.ScrollIntoView(ctx, el.id.ObjectID, options)
}

func (el *HTMLElement) Focus(ctx context.Context) error {
	return el.input.Focus(ctx, el.id.ObjectID)
}

func (el *HTMLElement) FocusBySelector(ctx context.Context, selector values.String) error {
	return el.input.FocusBySelector(ctx, el.id.NodeID, selector.String())
}

func (el *HTMLElement) Blur(ctx context.Context) error {
	return el.input.Blur(ctx, el.id.ObjectID)
}

func (el *HTMLElement) BlurBySelector(ctx context.Context, selector values.String) error {
	return el.input.BlurBySelector(ctx, el.id.ObjectID, selector.String())
}

func (el *HTMLElement) Hover(ctx context.Context) error {
	return el.input.MoveMouse(ctx, el.id.ObjectID)
}

func (el *HTMLElement) HoverBySelector(ctx context.Context, selector values.String) error {
	return el.input.MoveMouseBySelector(ctx, el.id.NodeID, selector.String())
}

func (el *HTMLElement) IsDetached() values.Boolean {
	el.mu.Lock()
	defer el.mu.Unlock()

	return !el.connected
}

func (el *HTMLElement) loadInnerHTML(ctx context.Context) (core.Value, error) {
	if el.IsDetached() {
		return values.None, drivers.ErrDetached
	}

	return getInnerHTML(ctx, el.client, el.exec, el.id, el.nodeType)
}

func (el *HTMLElement) loadInnerText(ctx context.Context) (core.Value, error) {
	if !el.IsDetached() {
		text, err := getInnerText(ctx, el.client, el.exec, el.id, el.nodeType)

		if err == nil {
			return text, nil
		}

		el.logError(err).Msg("failed to get inner text from remote object")

		// and just parse cached innerHTML
	}

	h, err := el.GetInnerHTML(ctx)

	if err != nil {
		el.logError(err).Msg("failed to get inner html from remote object")

		return values.None, err
	}

	if h == values.EmptyString {
		return h, nil
	}

	parsed, err := parseInnerText(h.String())

	if err != nil {
		el.logError(err).Msg("failed to parse inner html")

		return values.EmptyString, err
	}

	return parsed, nil
}

func (el *HTMLElement) loadAttrs(ctx context.Context) (core.Value, error) {
	repl, err := el.client.DOM.GetAttributes(ctx, dom.NewGetAttributesArgs(el.id.NodeID))

	if err != nil {
		return values.None, err
	}

	return parseAttrs(repl.Attributes), nil
}

func (el *HTMLElement) loadChildren(ctx context.Context) (core.Value, error) {
	if el.IsDetached() {
		return values.NewArray(0), nil
	}

	loaded := values.NewArray(len(el.children))

	for _, childID := range el.children {
		child, err := LoadHTMLElement(
			ctx,
			el.logger,
			el.client,
			el.dom,
			el.input,
			el.exec,
			childID.NodeID,
		)

		if err != nil {
			el.logError(err).Msg("failed to load child elements")

			continue
		}

		loaded.Push(child)
	}

	return loaded, nil
}

func (el *HTMLElement) parseStyle(ctx context.Context) (core.Value, error) {
	value, err := el.GetAttribute(ctx, "style")

	if err != nil {
		return values.None, err
	}

	if value == values.None {
		return values.NewObject(), nil
	}

	if value.Type() != types.String {
		return values.NewObject(), nil
	}

	return common.DeserializeStyles(value.(values.String))
}

func (el *HTMLElement) handlePageReload(_ context.Context) {
	el.Close()
}

func (el *HTMLElement) handleAttrModified(ctx context.Context, nodeID dom.NodeID, name, value string) {
	// it's not for this el
	if nodeID != el.id.NodeID {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	if el.IsDetached() {
		return
	}

	el.attributes.Mutate(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update element")

			return
		}

		if name == "style" {
			el.style.Reset()
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Set(values.NewString(name), values.NewString(value))
	})
}

func (el *HTMLElement) handleAttrRemoved(ctx context.Context, nodeID dom.NodeID, name string) {
	// it's not for this el
	if nodeID != el.id.NodeID {
		return
	}

	// they are not event loaded
	// just ignore the event
	if !el.attributes.Ready() {
		return
	}

	if el.IsDetached() {
		return
	}

	el.attributes.Mutate(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logError(err).Msg("failed to update element")

			return
		}

		if name == "style" {
			el.style.Reset()
		}

		attrs, ok := v.(*values.Object)

		if !ok {
			return
		}

		attrs.Remove(values.NewString(name))
	})
}

func (el *HTMLElement) handleChildrenCountChanged(ctx context.Context, nodeID dom.NodeID, _ int) {
	if nodeID != el.id.NodeID {
		return
	}

	if el.IsDetached() {
		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	node, err := el.client.DOM.DescribeNode(
		ctx,
		dom.NewDescribeNodeArgs().SetObjectID(el.id.ObjectID),
	)

	if err != nil {
		el.logError(err).Msg("failed to update element")

		return
	}

	el.children = createChildrenArray(node.Node.Children)
}

func (el *HTMLElement) handleChildInserted(ctx context.Context, parentNodeID, prevNodeID dom.NodeID, node dom.Node) {
	if parentNodeID != el.id.NodeID {
		return
	}

	targetIDx := -1
	prevID := prevNodeID
	nextID := node.NodeID

	if el.IsDetached() {
		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	for idx, id := range el.children {
		if id.NodeID == prevID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	nextIdentity := HTMLElementIdentity{
		NodeID: nextID,
	}

	arr := el.children
	el.children = append(arr[:targetIDx], append([]HTMLElementIdentity{nextIdentity}, arr[targetIDx:]...)...)

	if !el.loadedChildren.Ready() {
		return
	}

	el.loadedChildren.Mutate(ctx, func(v core.Value, _ error) {
		loadedArr := v.(*values.Array)
		loadedEl, err := LoadHTMLElement(
			ctx,
			el.logger,
			el.client,
			el.dom,
			el.input,
			el.exec,
			nextID,
		)

		if err != nil {
			el.logError(err).Msg("failed to load an inserted element")

			return
		}

		loadedArr.Insert(values.NewInt(targetIDx), loadedEl)

		el.innerHTML.Reset()
		el.innerText.Reset()
	})
}

func (el *HTMLElement) handleChildRemoved(ctx context.Context, nodeID, prevNodeID dom.NodeID) {
	if nodeID != el.id.NodeID {
		return
	}

	targetIDx := -1
	targetID := prevNodeID

	if el.IsDetached() {
		return
	}

	el.mu.Lock()
	defer el.mu.Unlock()

	for idx, id := range el.children {
		if id.NodeID == targetID {
			targetIDx = idx
			break
		}
	}

	if targetIDx == -1 {
		return
	}

	arr := el.children
	el.children = append(arr[:targetIDx], arr[targetIDx+1:]...)

	if !el.loadedChildren.Ready() {
		return
	}

	el.loadedChildren.Mutate(ctx, func(v core.Value, err error) {
		if err != nil {
			el.logger.Error().
				Timestamp().
				Err(err).
				Int("nodeID", int(el.id.NodeID)).
				Msg("failed to update element")

			return
		}

		loadedArr := v.(*values.Array)
		loadedArr.RemoveAt(values.NewInt(targetIDx))

		el.innerHTML.Reset()
		el.innerText.Reset()
	})
}

func (el *HTMLElement) logError(err error) *zerolog.Event {
	return el.logger.
		Error().
		Timestamp().
		Int("nodeID", int(el.id.NodeID)).
		Str("objectID", string(el.id.ObjectID)).
		Err(err)
}
