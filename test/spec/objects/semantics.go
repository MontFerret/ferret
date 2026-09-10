// Package objects provides fresh semantic fixtures for object transformations.
package objects

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type (
	MergeCase struct {
		Shallow runtime.Map
		Deep    runtime.Map
		Name    string
		Sources []runtime.Map
	}
	FilterCase struct {
		Source  runtime.Map
		Kept    runtime.Map
		Omitted runtime.Map
		Name    string
		Keys    []runtime.String
	}
)

// MergeCases returns independent inputs for destination, immutable, and mutable APIs.
func MergeCases() []MergeCase {
	return []MergeCase{
		{Name: "empty", Sources: []runtime.Map{runtime.NewObject()}, Shallow: runtime.NewObject(), Deep: runtime.NewObject()},
		{Name: "later wins", Sources: []runtime.Map{
			runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}),
			runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)}),
		}, Shallow: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)}), Deep: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)})},
		{Name: "three sources", Sources: []runtime.Map{
			runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}),
			runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.None}),
			runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(3)}),
		}, Shallow: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(3)}), Deep: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(3)})},
		{Name: "nested conflicts", Sources: []runtime.Map{
			runtime.NewObjectWith(map[string]runtime.Value{"map": runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1)}), "list": runtime.NewArrayWith(runtime.Int(1)), "gone": runtime.Int(1)}),
			runtime.NewObjectWith(map[string]runtime.Value{"map": runtime.NewObjectWith(map[string]runtime.Value{"right": runtime.Int(2)}), "list": runtime.NewArrayWith(runtime.Int(2)), "gone": runtime.None}),
		}, Shallow: runtime.NewObjectWith(map[string]runtime.Value{"map": runtime.NewObjectWith(map[string]runtime.Value{"right": runtime.Int(2)}), "list": runtime.NewArrayWith(runtime.Int(2)), "gone": runtime.None}),
			Deep: runtime.NewObjectWith(map[string]runtime.Value{"map": runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1), "right": runtime.Int(2)}), "list": runtime.NewArrayWith(runtime.Int(2)), "gone": runtime.None})},
	}
}

// FilterCases returns independent inputs with missing, repeated, and empty keys.
func FilterCases() []FilterCase {
	return []FilterCase{
		{Name: "missing repeated and present none", Source: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1), "b": runtime.None}), Keys: []runtime.String{"a", "missing", "a"},
			Kept: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), Omitted: runtime.NewObjectWith(map[string]runtime.Value{"b": runtime.None})},
		{Name: "empty keys", Source: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), Kept: runtime.NewObject(), Omitted: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)})},
		{Name: "empty map", Source: runtime.NewObject(), Keys: []runtime.String{"missing"}, Kept: runtime.NewObject(), Omitted: runtime.NewObject()},
	}
}
