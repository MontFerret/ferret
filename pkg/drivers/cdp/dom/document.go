package dom

import (
	"context"
	"hash/fnv"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

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

type HTMLDocument struct {
	logger    zerolog.Logger
	client    *cdp.Client
	dom       *Manager
	input     *input.Manager
	eval      *eval.Runtime
	frameTree page.FrameTree
	element   *HTMLElement
}

func NewHTMLDocument(
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	rootElement *HTMLElement,
	frames page.FrameTree,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logging.WithName(logger.With(), "html_document").Logger()
	doc.client = client
	doc.dom = domManager
	doc.input = input
	doc.eval = exec
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
		cdpDoc, ok := other.(*HTMLDocument)

		if ok {
			thisID := values.NewString(string(doc.Frame().Frame.ID))
			otherID := values.NewString(string(cdpDoc.Frame().Frame.ID))

			return thisID.Compare(otherID)
		}

		other := other.(drivers.HTMLDocument)

		return values.NewString(doc.frameTree.Frame.URL).Compare(other.GetURL())
	case FrameIDType:
		return values.NewString(string(doc.frameTree.Frame.ID)).Compare(values.NewString(other.String()))
	default:
		return drivers.Compare(doc.Type(), other.Type())
	}
}

func (doc *HTMLDocument) Iterate(ctx context.Context) (core.Iterator, error) {
	return doc.element.Iterate(ctx)
}

func (doc *HTMLDocument) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInDocument(ctx, path, doc)
}

func (doc *HTMLDocument) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInDocument(ctx, path, doc, value)
}

func (doc *HTMLDocument) Close() error {
	return doc.element.Close()
}

func (doc *HTMLDocument) Frame() page.FrameTree {
	return doc.frameTree
}

func (doc *HTMLDocument) GetNodeType(_ context.Context) (values.Int, error) {
	return 9, nil
}

func (doc *HTMLDocument) GetNodeName(_ context.Context) (values.String, error) {
	return "#document", nil
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) (*values.Array, error) {
	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector drivers.QuerySelector) (core.Value, error) {
	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*values.Array, error) {
	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Int, error) {
	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Boolean, error) {
	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) GetTitle() values.String {
	value, err := doc.eval.EvalValue(context.Background(), templates.GetTitle())

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
	return doc.input.MoveMouseByXY(ctx, x, y)
}

func (doc *HTMLDocument) WaitForElement(ctx context.Context, selector drivers.QuerySelector, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForElement(doc.element.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelector(ctx context.Context, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForClassBySelector(doc.element.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelectorAll(ctx context.Context, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForClassBySelectorAll(doc.element.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForAttributeBySelector(
	ctx context.Context,
	selector drivers.QuerySelector,
	name,
	value values.String,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForAttributeBySelector(doc.element.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForAttributeBySelectorAll(
	ctx context.Context,
	selector drivers.QuerySelector,
	name,
	value values.String,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForAttributeBySelectorAll(doc.element.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelector(ctx context.Context, selector drivers.QuerySelector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForStyleBySelector(doc.element.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelectorAll(ctx context.Context, selector drivers.QuerySelector, name, value values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.eval,
		templates.WaitForStyleBySelectorAll(doc.element.id, selector, name, value, when),
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

func (doc *HTMLDocument) ScrollBySelector(ctx context.Context, selector drivers.QuerySelector, options drivers.ScrollOptions) error {
	return doc.input.ScrollIntoViewBySelector(ctx, doc.element.id, selector, options)
}

func (doc *HTMLDocument) Scroll(ctx context.Context, options drivers.ScrollOptions) error {
	return doc.input.ScrollByXY(ctx, options)
}

func (doc *HTMLDocument) Eval() *eval.Runtime {
	return doc.eval
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
