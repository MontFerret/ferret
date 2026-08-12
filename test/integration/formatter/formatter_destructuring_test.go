package formatter_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterDestructuringBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`LET{name,age:years,nested:[first,_],}=user
RETURN[name,years,first]`, `LET { name, age: years, nested: [first, _] } = user
RETURN [name, years, first]`),
		S(`VAR[first,{value:alias},_]=values
RETURN alias`, `VAR [first, { value: alias }, _] = values
RETURN alias`),
		S(`FOR{name,score:points}IN users RETURN[name,points]`, `FOR { name, score: points } IN users
    RETURN [name, points]`),
		S(`FOR[name,_]IN rows{RETURN name}`, `FOR [name, _] IN rows {
    RETURN name
}`),
	})
}

func TestFormatterMultilineDestructuringBindings(t *testing.T) {
	RunSpecsWith(t, formatter.New(formatter.WithPrintWidth(36)), []Spec{
		S(`LET {firstName, lastName, metadata: {createdAt, updatedAt}} = user
RETURN firstName`, `LET {
    firstName,
    lastName,
    metadata: {
        createdAt,
        updatedAt
    }
} = user
RETURN firstName`),
		S(`FOR [firstValue, secondValue, thirdValue] IN rows RETURN firstValue`, `FOR [
    firstValue,
    secondValue,
    thirdValue
] IN rows
    RETURN firstValue`),
	})
}

func TestFormatterDestructuringComments(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`LET {
// primary name
name,
metadata: [
created,
// deliberately ignored
_
]
} = user
RETURN [name, created]`, `LET {
    // primary name
    name,
    metadata: [
        created,
        // deliberately ignored
        _
    ]
} = user
RETURN [name, created]`),
		S(`LET {
// validate keyed access
} = user
LET [
// validate indexed access
] = values
RETURN 1`, `LET {
    // validate keyed access
} = user
LET [
    // validate indexed access
] = values
RETURN 1`),
	})
}

func TestFormatterDestructuringBracketSpacing(t *testing.T) {
	RunSpecsWith(t, formatter.New(formatter.WithBracketSpacing(false)), []Spec{
		S(`LET { name, nested: { value } } = input RETURN [name, value]`, `LET {name, nested: {value}} = input
RETURN [name, value]`),
	})
}
