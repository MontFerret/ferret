package stdlib

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Set is an immutable selection of standard library capability groups.
type Set struct {
	groups  map[Group]struct{}
	invalid []Group
}

// Full returns a set containing every standard library group.
func Full() Set {
	set := Empty()

	for _, registration := range groupRegistrations {
		set.groups[registration.group] = struct{}{}
	}

	return set
}

// Safe returns the full standard library without filesystem or network groups.
func Safe() Set {
	return Full().Without(FS, NET)
}

// Empty returns a set containing no standard library groups.
func Empty() Set {
	return Set{
		groups: make(map[Group]struct{}),
	}
}

// Only returns a set containing only the provided standard library groups.
func Only(groups ...Group) Set {
	return Empty().With(groups...)
}

// With returns a new set that also contains the provided groups.
func (s Set) With(groups ...Group) Set {
	next := s.clone()

	for _, group := range groups {
		expanded, ok := expandGroup(group)
		if !ok {
			next.addInvalid(group)
			continue
		}

		for _, item := range expanded {
			next.groups[item] = struct{}{}
		}
	}

	return next
}

// Without returns a new set without the provided groups.
func (s Set) Without(groups ...Group) Set {
	next := s.clone()

	for _, group := range groups {
		expanded, ok := expandGroup(group)
		if !ok {
			next.addInvalid(group)
			continue
		}

		for _, item := range expanded {
			delete(next.groups, item)
		}
	}

	return next
}

// Register registers the selected standard library groups into the namespace.
func (s Set) Register(ns runtime.Namespace) error {
	if ns == nil {
		return fmt.Errorf("stdlib namespace cannot be nil")
	}

	if err := invalidGroupsError(s.invalid); err != nil {
		return err
	}

	for _, registration := range groupRegistrations {
		if _, ok := s.groups[registration.group]; ok {
			registration.register(ns)
		}
	}

	return nil
}

func (s Set) clone() Set {
	next := Set{
		groups:  make(map[Group]struct{}, len(s.groups)),
		invalid: append([]Group(nil), s.invalid...),
	}

	for group := range s.groups {
		next.groups[group] = struct{}{}
	}

	return next
}

func (s *Set) addInvalid(group Group) {
	for _, existing := range s.invalid {
		if existing == group {
			return
		}
	}

	s.invalid = append(s.invalid, group)
}
