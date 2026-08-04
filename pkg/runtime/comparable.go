package runtime

import (
	"bytes"
	"context"
	"reflect"
	"strings"
)

// Ordering is the normalized result of a strict runtime comparison.
type Ordering int8

const (
	Less    Ordering = -1
	Equal   Ordering = 0
	Greater Ordering = 1
)

type (
	// Equatable gives a host value an independent, fallible equality contract.
	//
	// Equal must be reflexive, symmetric, and transitive within the value's
	// comparison domain. It returns false without an error for an incompatible
	// value. If the value also implements Comparable, Equal must return true
	// exactly when Compare returns Equal. Semantically equal values must have equal
	// hashes; hashes are only candidate selectors and never prove equality.
	Equatable interface {
		Equal(ctx context.Context, other Value) (bool, error)
	}

	// Comparable gives a host value a fallible total ordering within its
	// comparison domain.
	//
	// Compare must be reflexive, antisymmetric, and transitive. It returns an
	// ErrInvalidOperation-compatible error for an incompatible value and returns
	// operational failures, including context cancellation, unchanged. Any
	// negative or positive result is normalized by runtime dispatch to Less or
	// Greater. If the value also implements Equatable, Equal must agree with
	// Equal returning true.
	Comparable interface {
		Compare(ctx context.Context, other Value) (Ordering, error)
	}

	// Comparator compares two values for a fallible, context-aware sort.
	Comparator = func(ctx context.Context, first, second Value) (Ordering, error)
)

type builtinComparison uint8

const (
	builtinComparisonUnknown builtinComparison = iota
	builtinComparisonNone
	builtinComparisonBoolean
	builtinComparisonInt
	builtinComparisonFloat
	builtinComparisonDuration
	builtinComparisonString
	builtinComparisonDateTime
	builtinComparisonBinary
	builtinComparisonList
	builtinComparisonMap
)

// EqualValues applies strict runtime equality. It does not apply language-level
// coercions. Incompatible values compare unequal.
func EqualValues(ctx context.Context, left, right Value) (Boolean, error) {
	leftKind := builtinComparisonKindOf(left)
	rightKind := builtinComparisonKindOf(right)

	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return equalBuiltinValues(ctx, left, right, leftKind, rightKind)
	}

	return dispatchEquality(ctx, left, right)
}

// CompareValues applies strict runtime ordering. It does not apply
// language-level coercions. Incompatible host values return ErrInvalidOperation.
func CompareValues(ctx context.Context, left, right Value) (Ordering, error) {
	leftKind := builtinComparisonKindOf(left)
	rightKind := builtinComparisonKindOf(right)

	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return compareBuiltinValues(ctx, left, right, leftKind, rightKind)
	}

	return dispatchOrdering(ctx, left, right)
}

func equalBuiltinValues(
	ctx context.Context,
	left, right Value,
	leftKind, rightKind builtinComparison,
) (Boolean, error) {
	if isNumericComparison(leftKind, rightKind) {
		return Boolean(compareNumericValues(left, right) == Equal), nil
	}

	if leftKind != rightKind {
		return False, nil
	}

	switch leftKind {
	case builtinComparisonNone:
		return True, nil
	case builtinComparisonBoolean:
		return left.(Boolean) == right.(Boolean), nil
	case builtinComparisonDuration:
		return left.(Duration) == right.(Duration), nil
	case builtinComparisonString:
		return left.(String) == right.(String), nil
	case builtinComparisonDateTime:
		return compareDateTimeValues(left.(DateTime), right.(DateTime)) == Equal, nil
	case builtinComparisonBinary:
		return Boolean(bytes.Equal(left.(Binary), right.(Binary))), nil
	case builtinComparisonList, builtinComparisonMap:
		return dispatchBuiltinCollectionEquality(ctx, left, right)
	default:
		return False, nil
	}
}

func compareBuiltinValues(
	ctx context.Context,
	left, right Value,
	leftKind, rightKind builtinComparison,
) (Ordering, error) {
	if isNumericComparison(leftKind, rightKind) {
		return compareNumericValues(left, right), nil
	}

	if leftKind != rightKind {
		return CompareTypes(left, right), nil
	}

	switch leftKind {
	case builtinComparisonNone:
		return Equal, nil
	case builtinComparisonBoolean:
		leftBoolean := left.(Boolean)
		rightBoolean := right.(Boolean)
		if leftBoolean == rightBoolean {
			return Equal, nil
		}
		if !leftBoolean {
			return Less, nil
		}
		return Greater, nil
	case builtinComparisonDuration:
		return compareOrdered(left.(Duration), right.(Duration)), nil
	case builtinComparisonString:
		return normalizeOrdering(Ordering(strings.Compare(string(left.(String)), string(right.(String))))), nil
	case builtinComparisonDateTime:
		return compareDateTimeValues(left.(DateTime), right.(DateTime)), nil
	case builtinComparisonBinary:
		return compareBinaryValues(left.(Binary), right.(Binary)), nil
	case builtinComparisonList, builtinComparisonMap:
		return dispatchBuiltinCollectionOrdering(ctx, left, right)
	default:
		return Equal, nil
	}
}

func dispatchBuiltinCollectionEquality(ctx context.Context, left, right Value) (Boolean, error) {
	receiver := left.(Equatable)
	other := right
	if canonicalReceiverIsRight(left, right) {
		receiver = right.(Equatable)
		other = left
	}

	result, err := receiver.Equal(ctx, other)
	return Boolean(result), err
}

func dispatchBuiltinCollectionOrdering(ctx context.Context, left, right Value) (Ordering, error) {
	receiver := left.(Comparable)
	other := right
	reverse := false
	if canonicalReceiverIsRight(left, right) {
		receiver = right.(Comparable)
		other = left
		reverse = true
	}

	result, err := receiver.Compare(ctx, other)
	if err != nil {
		return Equal, err
	}

	result = normalizeOrdering(result)
	if reverse {
		result = reverseOrdering(result)
	}

	return result, nil
}

func dispatchEquality(ctx context.Context, left, right Value) (Boolean, error) {
	result, handled, err := dispatchCompatibleEquality(ctx, left, right)
	if handled || err != nil {
		return result, err
	}

	return False, nil
}

func dispatchCompatibleEquality(
	ctx context.Context,
	left, right Value,
) (Boolean, bool, error) {
	if !compatibleComparisonDomain(left, right) {
		return False, false, nil
	}

	leftEquatable, leftOK := left.(Equatable)
	rightEquatable, rightOK := right.(Equatable)

	if !leftOK && !rightOK {
		return False, false, nil
	}

	receiver := leftEquatable
	other := right

	if !leftOK || (rightOK && canonicalReceiverIsRight(left, right)) {
		receiver = rightEquatable
		other = left
	}

	result, err := receiver.Equal(ctx, other)
	if err != nil {
		return False, true, err
	}

	return Boolean(result), true, nil
}

func dispatchOrdering(ctx context.Context, left, right Value) (Ordering, error) {
	result, handled, err := dispatchCompatibleOrdering(ctx, left, right)
	if handled || err != nil {
		return result, err
	}

	return Equal, incompatibleComparisonError(left, right)
}

func dispatchCompatibleOrdering(
	ctx context.Context,
	left, right Value,
) (Ordering, bool, error) {
	if !compatibleComparisonDomain(left, right) {
		return Equal, false, nil
	}

	leftComparable, leftOK := left.(Comparable)
	rightComparable, rightOK := right.(Comparable)

	if !leftOK && !rightOK {
		return Equal, false, nil
	}

	receiver := leftComparable
	other := right
	reverse := false

	if !leftOK || (rightOK && canonicalReceiverIsRight(left, right)) {
		receiver = rightComparable
		other = left
		reverse = true
	}

	result, err := receiver.Compare(ctx, other)
	if err != nil {
		return Equal, true, err
	}

	result = normalizeOrdering(result)
	if reverse {
		result = reverseOrdering(result)
	}

	return result, true, nil
}

func compatibleComparisonDomain(left, right Value) bool {
	_, leftList := left.(List)
	_, rightList := right.(List)
	if leftList || rightList {
		return leftList && rightList
	}

	_, leftMap := left.(Map)
	_, rightMap := right.(Map)
	if leftMap || rightMap {
		return leftMap && rightMap
	}

	leftDomain, leftStable := stableTypedComparisonDomain(left)
	rightDomain, rightStable := stableTypedComparisonDomain(right)
	if leftStable || rightStable {
		return leftStable && rightStable && leftDomain == rightDomain
	}

	return reflect.TypeOf(left) == reflect.TypeOf(right)
}

func stableTypedComparisonDomain(value Value) (string, bool) {
	typed, ok := value.(Typed)
	if !ok {
		return "", false
	}

	typ := typed.Type()
	name := TypeName(typ)
	return name, typ != nil && name != "" && name != UnknownTypeName
}

func canonicalReceiverIsRight(left, right Value) bool {
	leftType := reflect.TypeOf(left)
	rightType := reflect.TypeOf(right)

	if leftType == rightType {
		return false
	}

	return fullyQualifiedGoType(rightType) < fullyQualifiedGoType(leftType)
}

func builtinComparisonKindOf(value Value) builtinComparison {
	if value == nil || value == None {
		return builtinComparisonNone
	}

	switch value.(type) {
	case Boolean:
		return builtinComparisonBoolean
	case Int:
		return builtinComparisonInt
	case Float:
		return builtinComparisonFloat
	case Duration:
		return builtinComparisonDuration
	case String:
		return builtinComparisonString
	case DateTime:
		return builtinComparisonDateTime
	case Binary:
		return builtinComparisonBinary
	case List:
		return builtinComparisonList
	case Map:
		return builtinComparisonMap
	default:
		return builtinComparisonUnknown
	}
}

func isNumericComparison(left, right builtinComparison) bool {
	leftNumeric := left == builtinComparisonInt || left == builtinComparisonFloat
	rightNumeric := right == builtinComparisonInt || right == builtinComparisonFloat
	return leftNumeric && rightNumeric
}

func compareNumericValues(left, right Value) Ordering {
	if leftInt, ok := left.(Int); ok {
		if rightInt, ok := right.(Int); ok {
			return compareOrdered(leftInt, rightInt)
		}

		return compareOrdered(Float(leftInt), right.(Float))
	}

	leftFloat := left.(Float)
	if rightFloat, ok := right.(Float); ok {
		return compareOrdered(leftFloat, rightFloat)
	}

	return compareOrdered(leftFloat, Float(right.(Int)))
}

func compareDateTimeValues(left, right DateTime) Ordering {
	if left.After(right.Time) {
		return Greater
	}
	if left.Before(right.Time) {
		return Less
	}
	return Equal
}

func compareBinaryValues(left, right Binary) Ordering {
	if len(left) < len(right) {
		return Less
	}
	if len(left) > len(right) {
		return Greater
	}

	return normalizeOrdering(Ordering(bytes.Compare(left, right)))
}

func compareOrdered[T ~int64 | ~float64](left, right T) Ordering {
	if left == right {
		return Equal
	}
	if left < right {
		return Less
	}
	return Greater
}

func normalizeOrdering(result Ordering) Ordering {
	if result < 0 {
		return Less
	}
	if result > 0 {
		return Greater
	}
	return Equal
}

func reverseOrdering(result Ordering) Ordering {
	switch result {
	case Less:
		return Greater
	case Greater:
		return Less
	default:
		return Equal
	}
}

func incompatibleComparisonError(left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"cannot order %s and %s",
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}
