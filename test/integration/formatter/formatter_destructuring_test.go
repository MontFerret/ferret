package formatter_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterDestructuringBindings(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`let{name,age:years,nested:[first,_],}=user
return[name,years,first]`, `let { name, age: years, nested: [first, _] } = user
return [name, years, first]`),
		S(`var[first,{value:alias},_]=values
return alias`, `var [first, { value: alias }, _] = values
return alias`),
		S(`for{name,score:points}in users return[name,points]`, `for { name, score: points } in users
    return [name, points]`),
		S(`for[name,_]in rows{return name}`, `for [name, _] in rows {
    return name
}`),
	})
}

func TestFormatterMultilineDestructuringBindings(t *testing.T) {
	RunSpecsWith(t, mustNewFormatter(t, formatter.WithPrintWidth(36)), []Spec{
		S(`let {firstName, lastName, metadata: {createdAt, updatedAt}} = user
return firstName`, `let {
    firstName,
    lastName,
    metadata: {
        createdAt,
        updatedAt
    }
} = user
return firstName`),
		S(`for [firstValue, secondValue, thirdValue] in rows return firstValue`, `for [
    firstValue,
    secondValue,
    thirdValue
] in rows
    return firstValue`),
	})
}

func TestFormatterDestructuringComments(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`let {
// primary name
name,
metadata: [
created,
// deliberately ignored
_
]
} = user
return [name, created]`, `let {
    // primary name
    name,
    metadata: [
        created,
        // deliberately ignored
        _
    ]
} = user
return [name, created]`),
		S(`let {
// validate keyed access
} = user
let [
// validate indexed access
] = values
return 1`, `let {
    // validate keyed access
} = user
let [
    // validate indexed access
] = values
return 1`),
	})
}

func TestFormatterDestructuringBracketSpacing(t *testing.T) {
	RunSpecsWith(t, mustNewFormatter(t, formatter.WithBracketSpacing(false)), []Spec{
		S(`let { name, nested: { value } } = input return [name, value]`, `let {name, nested: {value}} = input
return [name, value]`),
	})
}
