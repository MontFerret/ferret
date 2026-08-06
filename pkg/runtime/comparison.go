package runtime

import (
	"bytes"
	"context"
	"reflect"
	"strings"
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

// EqualValues applies canonical runtime equality. Apart from mixed numeric
// comparison and compatible host domains, values must share a comparison type.
// It forwards ctx to host capabilities without polling cancellation.
func EqualValues(ctx context.Context, left, right Value) (Boolean, error) {
	leftKind := builtinComparisonKindOf(left)
	rightKind := builtinComparisonKindOf(right)

	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return equalBuiltinValues(ctx, left, right, leftKind, rightKind)
	}

	return dispatchEquality(ctx, left, right, compatibleComparisonDomain(left, right))
}

// CompareValues applies canonical runtime relational comparison. Apart from
// mixed numeric comparison and compatible host domains, native Duration values
// can be compared only with other native Duration values. Incompatible values
// return ErrInvalidOperation. It forwards ctx to host capabilities without
// polling cancellation.
func CompareValues(ctx context.Context, left, right Value) (Ordering, error) {
	leftKind := builtinComparisonKindOf(left)
	rightKind := builtinComparisonKindOf(right)

	if leftKind != builtinComparisonUnknown && rightKind != builtinComparisonUnknown {
		return compareBuiltinValues(ctx, left, right, leftKind, rightKind)
	}

	return dispatchOrdering(ctx, left, right, compatibleComparisonDomain(left, right))
}

func equalBuiltinValues(
	ctx context.Context,
	left, right Value,
	leftKind, rightKind builtinComparison,
) (Boolean, error) {
	if isNumericComparison(leftKind, rightKind) {
		return compareNumericValues(left, right) == Equal, nil
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
		return dispatchEquality(ctx, left, right, true)
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
		if leftKind == builtinComparisonDuration || rightKind == builtinComparisonDuration {
			return Equal, incompatibleComparisonError(left, right)
		}

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
		return dispatchOrdering(ctx, left, right, true)
	default:
		return Equal, nil
	}
}

func dispatchEquality(ctx context.Context, left, right Value, compatible bool) (Boolean, error) {
	if !compatible {
		return False, nil
	}

	leftReceiver, leftOK := left.(Equatable)
	receiver := leftReceiver
	other := right

	if !leftOK || reflect.TypeOf(left) != reflect.TypeOf(right) {
		rightReceiver, rightOK := right.(Equatable)
		useRight, ok := selectComparisonReceiver(left, right, leftOK, rightOK)
		if !ok {
			return False, nil
		}

		if useRight {
			receiver = rightReceiver
			other = left
		}
	}

	result, err := receiver.Equal(ctx, other)
	if err != nil {
		return False, err
	}

	return Boolean(result), nil
}

func dispatchOrdering(ctx context.Context, left, right Value, compatible bool) (Ordering, error) {
	if !compatible {
		return Equal, incompatibleComparisonError(left, right)
	}

	leftReceiver, leftOK := left.(Comparable)
	receiver := leftReceiver
	other := right
	reverse := false

	if !leftOK || reflect.TypeOf(left) != reflect.TypeOf(right) {
		rightReceiver, rightOK := right.(Comparable)
		useRight, ok := selectComparisonReceiver(left, right, leftOK, rightOK)
		if !ok {
			return Equal, incompatibleComparisonError(left, right)
		}

		if useRight {
			receiver = rightReceiver
			other = left
			reverse = true
		}
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

func selectComparisonReceiver(left, right Value, leftOK, rightOK bool) (useRight, ok bool) {
	if !leftOK && !rightOK {
		return false, false
	}

	if !leftOK || (rightOK && canonicalReceiverIsRight(left, right)) {
		return true, true
	}

	return false, true
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

		return compareIntFloatValues(leftInt, right.(Float))
	}

	leftFloat := left.(Float)
	if rightFloat, ok := right.(Float); ok {
		return compareFloatValues(leftFloat, rightFloat)
	}

	return reverseOrdering(compareIntFloatValues(right.(Int), leftFloat))
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
