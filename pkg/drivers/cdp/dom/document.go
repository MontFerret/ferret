package dom

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
	logger    *zerolog.Logger
	client    *cdp.Client
	dom       *Manager
	input     *input.Manager
	exec      *eval.ExecutionContext
	frameTree page.FrameTree
	element   *HTMLElement
}

func LoadRootHTMLDocument(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
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
		domManager,
		mouse,
		keyboard,
		gdRepl.Root,
		ftRepl.FrameTree,
		worldRepl.ExecutionContextID,
	)
}

func LoadHTMLDocument(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
	node dom.Node,
	frameTree page.FrameTree,
	execID runtime.ExecutionContextID,
) (*HTMLDocument, error) {
	exec := eval.NewExecutionContext(client, frameTree.Frame, execID)
	inputManager := input.NewManager(client, exec, keyboard, mouse)

	rootElement, err := LoadHTMLElement(
		ctx,
		logger,
		client,
		domManager,
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
		domManager,
		inputManager,
		exec,
		rootElement,
		frameTree,
	), nil
}

func NewHTMLDocument(
	logger *zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.ExecutionContext,
	rootElement *HTMLElement,
	frames page.FrameTree,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logger
	doc.client = client
	doc.dom = domManager
	doc.input = input
	doc.exec = exec
	doc.element = rootElement
	doc.frameTree = frames

	return doc
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	return doc.frameTree.Frame.URL
}

func (doc *HTMLDocument) Unwrap() interface{} {
	return doc.element
}

func (doc *HTMLDocument) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(doc.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(doc.frameTree.Frame.ID))
	h.Write([]byte(doc.frameTree.Frame.URL))

	return h.Sum64()
}

func (doc *HTMLDocument) Copy() core.Value {
	return values.None
}

func (doc *HTMLDocument) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLDocumentType:
		other := other.(drivers.HTMLDocument)

		return values.NewString(doc.frameTree.Frame.URL).Compare(other.GetURL())
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
	return doc.element.Close()
}

func (doc *HTMLDocument) Frame() page.FrameTree {
	return doc.frameTree
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
	if doc.frameTree.Frame.Name != nil {
		return values.NewString(*doc.frameTree.Frame.Name)
	}

	return values.EmptyString
}

func (doc *HTMLDocument) GetParentDocument(ctx context.Context) (drivers.HTMLDocument, error) {
	if doc.frameTree.Frame.ParentID == nil {
		return nil, nil
	}

	return doc.dom.GetFrameNode(ctx, *doc.frameTree.Frame.ParentID)
}

func (doc *HTMLDocument) GetChildDocuments(ctx context.Context) (*values.Array, error) {
	arr := values.NewArray(len(doc.frameTree.ChildFrames))

	for _, childFrame := range doc.frameTree.ChildFrames {
		frame, err := doc.dom.GetFrameNode(ctx, childFrame.Frame.ID)

		if err != nil {
			return nil, err
		}

		if frame != nil {
			arr.Push(frame)
		}
	}

	return arr, nil
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
	return values.NewString(doc.frameTree.Frame.URL)
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

func (doc *HTMLDocument) ScrollTop(ctx context.Context, options drivers.ScrollOptions) error {
	return doc.input.ScrollTop(ctx, options)
}

func (doc *HTMLDocument) ScrollBottom(ctx context.Context, options drivers.ScrollOptions) error {
	return doc.input.ScrollBottom(ctx, options)
}

func (doc *HTMLDocument) ScrollBySelector(ctx context.Context, selector values.String, options drivers.ScrollOptions) error {
	return doc.input.ScrollIntoViewBySelector(ctx, selector.String(), options)
}

func (doc *HTMLDocument) ScrollByXY(ctx context.Context, x, y values.Float, options drivers.ScrollOptions) error {
	return doc.input.ScrollByXY(ctx, float64(x), float64(y), options)
}

func (doc *HTMLDocument) logError(err error) *zerolog.Event {
	return doc.logger.
		Error().
		Timestamp().
		Str("url", doc.frameTree.Frame.URL).
		Str("securityOrigin", doc.frameTree.Frame.SecurityOrigin).
		Str("mimeType", doc.frameTree.Frame.MimeType).
		Str("frameID", string(doc.frameTree.Frame.ID)).
		Err(err)
}
