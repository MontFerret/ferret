package cdp

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLDocument struct {
	logger   *zerolog.Logger
	client   *cdp.Client
	events   *events.EventBroker
	input    *input.Manager
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
	events *events.EventBroker,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
) (*HTMLDocument, error) {
	gdRepl, err := client.DOM.GetDocument(ctx, dom.NewGetDocumentArgs().SetDepth(1))

	if err != nil {
		return nil, err
	}

	ftRepl, err := client.Page.GetFrameTree(ctx)

	if err != nil {
		return nil, err
	}

	worldRepl, err := client.Page.CreateIsolatedWorld(ctx, page.NewCreateIsolatedWorldArgs(ftRepl.FrameTree.Frame.ID))

	if err != nil {
		return nil, err
	}

	return LoadHTMLDocument(
		ctx,
		logger,
		client,
		events,
		mouse,
		keyboard,
		gdRepl.Root,
		ftRepl.FrameTree,
		worldRepl.ExecutionContextID,
		nil,
	)
}

func LoadHTMLDocument(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	events *events.EventBroker,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
	node dom.Node,
	tree page.FrameTree,
	execID runtime.ExecutionContextID,
	parent *HTMLDocument,
) (*HTMLDocument, error) {
	exec := eval.NewExecutionContext(client, tree.Frame, execID)
	inputManager := input.NewManager(client, exec, keyboard, mouse)

	rootElement, err := LoadHTMLElement(
		ctx,
		logger,
		client,
		events,
		inputManager,
		exec,
		node.NodeID,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load root element")
	}

	return NewHTMLDocument(
		logger,
		client,
		events,
		inputManager,
		exec,
		rootElement,
		tree,
		parent,
	), nil
}

func NewHTMLDocument(
	logger *zerolog.Logger,
	client *cdp.Client,
	events *events.EventBroker,
	input *input.Manager,
	exec *eval.ExecutionContext,
	rootElement *HTMLElement,
	frames page.FrameTree,
	parent *HTMLDocument,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logger
	doc.client = client
	doc.events = events
	doc.input = input
	doc.exec = exec
	doc.element = rootElement
	doc.frames = frames
	doc.parent = parent
	doc.children = common.NewLazyValue(doc.loadChildren)

	return doc
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	return doc.frames.Frame.URL
}

func (doc *HTMLDocument) Unwrap() interface{} {
	return doc.element
}

func (doc *HTMLDocument) Hash() uint64 {
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
	return doc.element.IsDetached()
}

func (doc *HTMLDocument) GetNodeType() values.Int {
	return 9
}

func (doc *HTMLDocument) GetNodeName() values.String {
	return "#document"
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) (*values.Array, error) {
	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector values.String) (core.Value, error) {
	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector values.String) (*values.Array, error) {
	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector values.String) (values.Int, error) {
	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector values.String) (values.Boolean, error) {
	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) GetTitle() values.String {
	value, err := doc.exec.ReadProperty(context.Background(), doc.element.id.ObjectID, "title")

	if err != nil {
		doc.logError(errors.Wrap(err, "failed to read document title"))

		return values.EmptyString
	}

	return values.NewString(value.String())
}

func (doc *HTMLDocument) GetName() values.String {
	if doc.frames.Frame.Name != nil {
		return values.NewString(*doc.frames.Frame.Name)
	}

	return values.EmptyString
}

func (doc *HTMLDocument) GetParentDocument() drivers.HTMLDocument {
	return doc.parent
}

func (doc *HTMLDocument) GetChildDocuments(ctx context.Context) (*values.Array, error) {
	children, err := doc.children.Read(ctx)

	if err != nil {
		return values.NewArray(0), errors.Wrap(err, "failed to load child documents")
	}

	return children.Copy().(*values.Array), nil
}

func (doc *HTMLDocument) XPath(ctx context.Context, expression values.String) (core.Value, error) {
	return doc.element.XPath(ctx, expression)
}

func (doc *HTMLDocument) Length() values.Int {
	return doc.element.Length()
}

func (doc *HTMLDocument) GetElement() drivers.HTMLElement {
	return doc.element
}

func (doc *HTMLDocument) GetURL() values.String {
	return values.NewString(doc.frames.Frame.URL)
}

func (doc *HTMLDocument) MoveMouseByXY(ctx context.Context, x, y values.Float) error {
	return doc.input.MoveMouseByXY(ctx, float64(x), float64(y))
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
	return doc.input.ScrollTop(ctx)
}

func (doc *HTMLDocument) ScrollBottom(ctx context.Context) error {
	return doc.input.ScrollBottom(ctx)
}

func (doc *HTMLDocument) ScrollBySelector(ctx context.Context, selector values.String) error {
	return doc.input.ScrollIntoViewBySelector(ctx, selector.String())
}

func (doc *HTMLDocument) ScrollByXY(ctx context.Context, x, y values.Float) error {
	return doc.input.ScrollByXY(ctx, float64(x), float64(y))
}

func (doc *HTMLDocument) loadChildren(ctx context.Context) (value core.Value, e error) {
	children := values.NewArray(len(doc.frames.ChildFrames))

	if len(doc.frames.ChildFrames) > 0 {
		for _, cf := range doc.frames.ChildFrames {
			cfNode, cfExecID, err := resolveFrame(ctx, doc.client, cf.Frame)

			if err != nil {
				return nil, errors.Wrap(err, "failed to resolve frame node")
			}

			cfDocument, err := LoadHTMLDocument(
				ctx,
				doc.logger,
				doc.client,
				doc.events,
				doc.input.Mouse(),
				doc.input.Keyboard(),
				cfNode,
				cf,
				cfExecID,
				doc,
			)

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
