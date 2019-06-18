package cdp

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

const BlankPageURL = "about:blank"

type HTMLDocument struct {
	mu       sync.Mutex
	logger   *zerolog.Logger
	client   *cdp.Client
	events   *events.EventBroker
	exec     *eval.ExecutionContext
	frames   page.FrameTree
	element  *HTMLElement
	parent   *HTMLDocument
	children *common.LazyValue
}

func LoadRootHTMLDocument(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
) (*HTMLDocument, error) {
	gdRepl, err := client.DOM.GetDocument(ctx, dom.NewGetDocumentArgs().SetDepth(1))

	if err != nil {
		return nil, err
	}

	ftRepl, err := client.Page.GetFrameTree(ctx)

	if err != nil {
		return nil, err
	}

	return LoadHTMLDocument(
		ctx,
		logger,
		client,
		broker,
		gdRepl.Root,
		ftRepl.FrameTree,
		eval.EmptyExecutionContextID,
		nil,
	)
}

func LoadHTMLDocument(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	node dom.Node,
	tree page.FrameTree,
	execID runtime.ExecutionContextID,
	parent *HTMLDocument,
) (*HTMLDocument, error) {
	exec := eval.NewExecutionContext(client, tree.Frame, execID)

	rootElement, err := LoadHTMLElement(
		ctx,
		logger,
		client,
		broker,
		exec,
		node.NodeID,
		node.BackendNodeID,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load root element")
	}

	return NewHTMLDocument(
		logger,
		client,
		broker,
		exec,
		rootElement,
		tree,
		parent,
	), nil
}

func NewHTMLDocument(
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	exec *eval.ExecutionContext,
	rootElement *HTMLElement,
	frames page.FrameTree,
	parent *HTMLDocument,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logger
	doc.client = client
	doc.events = broker
	doc.exec = exec
	doc.element = rootElement
	doc.frames = frames
	doc.parent = parent
	doc.children = common.NewLazyValue(doc.loadChildren)

	return doc
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.frames.Frame.URL
}

func (doc *HTMLDocument) Unwrap() interface{} {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element
}

func (doc *HTMLDocument) Hash() uint64 {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	h := fnv.New64a()

	h.Write([]byte(doc.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(doc.frames.Frame.ID))
	h.Write([]byte(doc.frames.Frame.URL))

	return h.Sum64()
}

func (doc *HTMLDocument) Copy() core.Value {
	return values.None
}

func (doc *HTMLDocument) Compare(other core.Value) int64 {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	switch other.Type() {
	case drivers.HTMLDocumentType:
		other := other.(drivers.HTMLDocument)

		return values.NewString(doc.frames.Frame.URL).Compare(other.GetURL())
	default:
		return drivers.Compare(doc.Type(), other.Type())
	}
}

func (doc *HTMLDocument) Iterate(ctx context.Context) (core.Iterator, error) {
	return doc.element.Iterate(ctx)
}

func (doc *HTMLDocument) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInDocument(ctx, doc, path)
}

func (doc *HTMLDocument) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInDocument(ctx, doc, path, value)
}

func (doc *HTMLDocument) Close() error {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	errs := make([]error, 0, 5)

	if doc.children.Ready() {
		val, err := doc.children.Read(context.Background())

		if err == nil {
			arr := val.(*values.Array)

			arr.ForEach(func(value core.Value, _ int) bool {
				doc := value.(drivers.HTMLDocument)

				err := doc.Close()

				if err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to close nested document: %s", doc.GetURL()))
				}

				return true
			})
		} else {
			errs = append(errs, err)
		}
	}

	err := doc.element.Close()

	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	return core.Errors(errs...)
}

func (doc *HTMLDocument) IsDetached() values.Boolean {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.IsDetached()
}

func (doc *HTMLDocument) GetNodeType() values.Int {
	return 9
}

func (doc *HTMLDocument) GetNodeName() values.String {
	return "#document"
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) core.Value {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) core.Value {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector values.String) core.Value {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector values.String) core.Value {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector values.String) values.Int {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector values.String) values.Boolean {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) GetTitle() values.String {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	value, err := doc.exec.ReadProperty(context.Background(), doc.element.id.objectID, "title")

	if err != nil {
		doc.logError(errors.Wrap(err, "failed to read document title"))

		return values.EmptyString
	}

	return values.NewString(value.String())
}

func (doc *HTMLDocument) GetName() values.String {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	if doc.frames.Frame.Name != nil {
		return values.NewString(*doc.frames.Frame.Name)
	}

	return values.EmptyString
}

func (doc *HTMLDocument) GetParentDocument() drivers.HTMLDocument {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.parent
}

func (doc *HTMLDocument) GetChildDocuments(ctx context.Context) (*values.Array, error) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	children, err := doc.children.Read(ctx)

	if err != nil {
		return values.NewArray(0), errors.Wrap(err, "failed to load child documents")
	}

	return children.Copy().(*values.Array), nil
}

func (doc *HTMLDocument) Length() values.Int {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element.Length()
}

func (doc *HTMLDocument) GetElement() drivers.HTMLElement {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return doc.element
}

func (doc *HTMLDocument) GetURL() values.String {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	return values.NewString(doc.frames.Frame.URL)
}

func (doc *HTMLDocument) ClickBySelector(ctx context.Context, selector values.String) (values.Boolean, error) {
	res, err := doc.exec.EvalWithReturn(
		ctx,
		fmt.Sprintf(`
			var el = document.querySelector(%s);
			if (el == null) {
				return false;
			}
			var evt = new window.MouseEvent('click', { bubbles: true, cancelable: true });
			el.dispatchEvent(evt);
			return true;
		`, eval.ParamString(selector.String())),
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) ClickBySelectorAll(ctx context.Context, selector values.String) (values.Boolean, error) {
	res, err := doc.exec.EvalWithReturn(
		ctx,
		fmt.Sprintf(`
			var elements = document.querySelectorAll(%s);
			if (elements == null) {
				return false;
			}
			elements.forEach((el) => {
				var evt = new window.MouseEvent('click', { bubbles: true, cancelable: true });
				el.dispatchEvent(evt);	
			});
			return true;
		`, eval.ParamString(selector.String())),
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) (values.Boolean, error) {
	valStr := value.String()

	res, err := doc.exec.EvalWithReturn(
		ctx,
		fmt.Sprintf(`
			var el = document.querySelector(%s);
			if (el == null) {
				return false;
			}
			el.focus();
			return true;
		`, eval.ParamString(selector.String())),
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean && res.(values.Boolean) == values.False {
		return values.False, nil
	}

	// Initial delay after focusing but before typing
	time.Sleep(time.Duration(delay) * time.Millisecond)

	for _, ch := range valStr {
		for _, ev := range []string{"keyDown", "keyUp"} {
			ke := input.NewDispatchKeyEventArgs(ev).SetText(string(ch))

			if err := doc.client.Input.DispatchKeyEvent(ctx, ke); err != nil {
				return values.False, err
			}
		}

		time.Sleep(randomDuration(delay) * time.Millisecond)
	}

	return values.True, nil
}

func (doc *HTMLDocument) SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error) {
	res, err := doc.exec.EvalWithReturn(
		ctx,
		fmt.Sprintf(`
			var element = document.querySelector(%s);
			if (element == null) {
				return [];
			}
			var values = %s;
			if (element.nodeName.toLowerCase() !== 'select') {
				throw new Error('GetElement is not a <select> element.');
			}
			var options = Array.from(element.options);
      		element.value = undefined;
			for (var option of options) {
        		option.selected = values.includes(option.value);
        	
				if (option.selected && !element.multiple) {
          			break;
				}
      		}
      		element.dispatchEvent(new Event('input', { 'bubbles': true, cancelable: true }));
      		element.dispatchEvent(new Event('change', { 'bubbles': true, cancelable: true }));
      		
			return options.filter(option => option.selected).map(option => option.value);
		`,
			eval.ParamString(selector.String()),
			value.String(),
		),
	)

	if err != nil {
		return nil, err
	}

	arr, ok := res.(*values.Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(types.Array, res.Type())
}

func (doc *HTMLDocument) MoveMouseBySelector(ctx context.Context, selector values.String) error {
	err := doc.ScrollBySelector(ctx, selector)

	if err != nil {
		return err
	}

	selectorArgs := dom.NewQuerySelectorArgs(doc.element.id.nodeID, selector.String())
	found, err := doc.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		doc.element.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return err
	}

	if found.NodeID <= 0 {
		return errors.New("element not found")
	}

	q, err := getClickablePoint(ctx, doc.client, HTMLElementIdentity{
		nodeID: found.NodeID,
	})

	if err != nil {
		return err
	}

	return doc.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", q.X, q.Y),
	)
}

func (doc *HTMLDocument) MoveMouseByXY(ctx context.Context, x, y values.Float) error {
	return doc.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", float64(x), float64(y)),
	)
}

func (doc *HTMLDocument) WaitForElement(ctx context.Context, selector values.String, when drivers.WaitEvent) error {
	var operator string

	if when == drivers.WaitEventPresence {
		operator = "!="
	} else {
		operator = "=="
	}

	task := events.NewEvalWaitTask(
		doc.exec,
		fmt.Sprintf(
			`
				var el = document.querySelector(%s);
				
				if (el %s null) {
					return true;
				}
				
				// null means we need to repeat
				return null;
			`,
			eval.ParamString(selector.String()),
			operator,
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelector(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelector(
			selector,
			when,
			class,
			fmt.Sprintf("el.className.split(' ').find(i => i === %s)", eval.ParamString(class.String())),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelectorAll(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelectorAll(
			selector,
			when,
			class,
			fmt.Sprintf("el.className.split(' ').find(i => i === %s)", eval.ParamString(class.String())),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForAttributeBySelector(
	ctx context.Context,
	selector,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelector(
			selector,
			when,
			value,
			templates.AttributeRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForAttributeBySelectorAll(
	ctx context.Context,
	selector,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelectorAll(
			selector,
			when,
			value,
			templates.AttributeRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelector(ctx context.Context, selector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelector(
			selector,
			when,
			value,
			templates.StyleRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelectorAll(ctx context.Context, selector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.exec,
		templates.WaitBySelectorAll(
			selector,
			when,
			value,
			templates.StyleRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) ScrollTop(ctx context.Context) error {
	return doc.exec.Eval(ctx, `
		window.scrollTo({
			left: 0,
			top: 0,
    		behavior: 'instant'
  		});
	`)
}

func (doc *HTMLDocument) ScrollBottom(ctx context.Context) error {
	return doc.exec.Eval(ctx, `
		window.scrollTo({
			left: 0,
			top: window.document.body.scrollHeight,
    		behavior: 'instant'
  		});
	`)
}

func (doc *HTMLDocument) ScrollBySelector(ctx context.Context, selector values.String) error {
	return doc.exec.Eval(ctx, fmt.Sprintf(`
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

func (doc *HTMLDocument) ScrollByXY(ctx context.Context, x, y values.Float) error {
	return doc.exec.Eval(ctx, fmt.Sprintf(`
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

func (doc *HTMLDocument) loadChildren(ctx context.Context) (value core.Value, e error) {
	children := values.NewArray(len(doc.frames.ChildFrames))

	if len(doc.frames.ChildFrames) > 0 {
		for _, cf := range doc.frames.ChildFrames {
			cfNode, cfExecID, err := resolveFrame(ctx, doc.client, cf.Frame)

			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve frame node")
			}

			cfDocument, err := LoadHTMLDocument(ctx, doc.logger, doc.client, doc.events, cfNode, cf, cfExecID, doc)

			if err != nil {
				return nil, errors.Wrap(err, "failed to load frame document")
			}

			children.Push(cfDocument)
		}
	}

	return children, nil
}

func (doc *HTMLDocument) logError(err error) *zerolog.Event {
	return doc.logger.
		Error().
		Timestamp().
		Str("url", string(doc.frames.Frame.URL)).
		Str("securityOrigin", string(doc.frames.Frame.SecurityOrigin)).
		Str("mimeType", string(doc.frames.Frame.MimeType)).
		Str("frameID", string(doc.frames.Frame.ID)).
		Err(err)
}
