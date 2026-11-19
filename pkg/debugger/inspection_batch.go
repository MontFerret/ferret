package debugger

import (
	"context"
	"errors"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ErrValueReferencesExhausted reports that this inspection cannot issue all its handles.
var ErrValueReferencesExhausted = errors.New("debugger value references exhausted")

const (
	evaluationResult uint8 = iota
	localBinding
	paramBinding
	indexedChild
	namedChild
)

type (
	occurrence struct {
		name     string
		parent   ValueReference
		position int64
		frame    int
		kind     uint8
	}

	stagedReference struct {
		value runtime.Value
		key   occurrence
	}

	inspectionBatch struct {
		owner  *inspector
		staged []stagedReference
	}
)

func (b *inspectionBatch) describe(ctx context.Context, key occurrence, raw runtime.Value) (Value, error) {
	if err := ctx.Err(); err != nil {
		return Value{}, err
	}

	d := b.owner.values.Describe(raw)
	result := formatValue(raw, d, b.owner.format.MaxBytes)
	named, indexed := d.Children.Named, d.Children.Indexed
	if !named.CanExpand() && !indexed.CanExpand() {
		return result, nil
	}

	if ref := b.owner.occurrences[key]; ref.Valid() {
		result.Reference = ref

		return result, nil
	}

	if len(b.staged) >= b.owner.handleLimit-len(b.owner.valueRefs) || len(b.staged) >= math.MaxInt-int(b.owner.lastValueRef) {
		return Value{}, ErrValueReferencesExhausted
	}

	b.staged = append(b.staged, stagedReference{key: key, value: raw})
	result.Reference = b.owner.lastValueRef + ValueReference(len(b.staged))

	return result, nil
}

// The command lock keeps stop identities stable. The lifecycle lock orders
// this final context check and commit with Close and native termination.
func (b *inspectionBatch) publish(ctx context.Context) error {
	life := b.owner.lifecycle
	if err := life.lockPublication(ctx, "inspect"); err != nil {
		return err
	}

	defer life.unlockPublication()

	for _, entry := range b.staged {
		b.owner.lastValueRef++
		ref := b.owner.lastValueRef
		b.owner.valueRefs[ref] = entry.value
		if entry.key.kind != evaluationResult {
			b.owner.occurrences[entry.key] = ref
		}
	}

	return nil
}
