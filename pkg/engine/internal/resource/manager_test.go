package resource_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
)

func TestOwnershipAndReplacement(t *testing.T) {
	for _, firstOwned := range []bool{false, true} {
		for _, nextOwned := range []bool{false, true} {
			name := map[bool]string{false: "borrowed", true: "owned"}
			t.Run(name[firstOwned]+" to "+name[nextOwned], func(t *testing.T) {
				manager := resource.NewManager()
				firstCloses, nextCloses := 0, 0
				var err error

				if firstOwned {
					err = manager.Own(resource.Network, func() error {
						firstCloses++

						return nil
					})
				} else {
					err = manager.Borrow(resource.Network)
				}

				if err != nil {
					t.Fatal(err)
				}

				if nextOwned {
					err = manager.Own(resource.Network, func() error {
						nextCloses++

						return nil
					})
				} else {
					err = manager.Borrow(resource.Network)
				}

				if err != nil {
					t.Fatal(err)
				}

				wantFirst := 0
				if firstOwned {
					wantFirst = 1
				}

				if firstCloses != wantFirst || nextCloses != 0 {
					t.Fatalf("replacement closes = %d/%d, want %d/0", firstCloses, nextCloses, wantFirst)
				}

				wantNext := 0
				if nextOwned {
					wantNext = 1
				}

				for range 2 {
					if err := manager.Close(); err != nil {
						t.Fatal(err)
					}

					if firstCloses != wantFirst || nextCloses != wantNext {
						t.Fatalf("shutdown closes = %d/%d, want %d/%d", firstCloses, nextCloses, wantFirst, wantNext)
					}
				}
			})
		}
	}
}

func TestCloseReverseOrderAndJoinErrors(t *testing.T) {
	manager := resource.NewManager()
	firstErr, lastErr := errors.New("first failure"), errors.New("last failure")
	var order []resource.Name

	for _, item := range []struct {
		err  error
		name resource.Name
	}{
		{name: resource.FileSystem, err: firstErr},
		{name: resource.Network},
		{name: "future", err: lastErr},
	} {
		if err := manager.Own(item.name, func() error {
			order = append(order, item.name)

			return item.err
		}); err != nil {
			t.Fatal(err)
		}
	}

	if err := manager.Borrow("borrowed"); err != nil {
		t.Fatal(err)
	}

	err := manager.Close()
	if !errors.Is(err, firstErr) || !errors.Is(err, lastErr) {
		t.Fatalf("cleanup lost error causes: %v", err)
	}

	if got, want := err.Error(), "close future: last failure\nclose filesystem: first failure"; got != want {
		t.Fatalf("cleanup labels/order = %q, want %q", got, want)
	}

	if again := manager.Close(); again != err {
		t.Fatalf("repeated close did not retain its result: %v", again)
	}

	if want := []resource.Name{"future", resource.Network, resource.FileSystem}; !reflect.DeepEqual(order, want) {
		t.Fatalf("close order = %v, want %v", order, want)
	}
}

func TestReplacementBecomesNewestAcquisition(t *testing.T) {
	manager := resource.NewManager()
	var order []string

	for _, item := range []struct {
		name  resource.Name
		label string
	}{
		{name: resource.Network, label: "old network"},
		{name: resource.FileSystem, label: "filesystem"},
		{name: resource.Network, label: "new network"},
	} {
		if err := manager.Own(item.name, func() error {
			order = append(order, item.label)

			return nil
		}); err != nil {
			t.Fatal(err)
		}
	}

	if want := []string{"old network"}; !reflect.DeepEqual(order, want) {
		t.Fatalf("replacement order = %v, want %v", order, want)
	}

	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}

	if want := []string{"old network", "new network", "filesystem"}; !reflect.DeepEqual(order, want) {
		t.Fatalf("close order = %v, want %v", order, want)
	}
}

func TestReplacementIsAcceptedAfterRetirementError(t *testing.T) {
	for _, borrowed := range []bool{false, true} {
		t.Run(map[bool]string{false: "owned replacement", true: "borrowed replacement"}[borrowed], func(t *testing.T) {
			manager := resource.NewManager()
			retirementErr := errors.New("retirement failed")
			oldCloses, newCloses := 0, 0

			if err := manager.Own(resource.Network, func() error {
				oldCloses++

				return retirementErr
			}); err != nil {
				t.Fatal(err)
			}

			var err error

			if borrowed {
				err = manager.Borrow(resource.Network)
			} else {
				err = manager.Own(resource.Network, func() error {
					newCloses++

					return nil
				})
			}

			if !errors.Is(err, retirementErr) || !strings.Contains(err.Error(), "close network") {
				t.Fatalf("retirement error = %v", err)
			}

			if oldCloses != 1 || newCloses != 0 {
				t.Fatalf("replacement closes = %d/%d, want 1/0", oldCloses, newCloses)
			}

			for range 2 {
				if err := manager.Close(); err != nil {
					t.Fatalf("shutdown repeated retirement error: %v", err)
				}
			}

			wantNew := 1
			if borrowed {
				wantNew = 0
			}

			if oldCloses != 1 || newCloses != wantNew {
				t.Fatalf("shutdown closes = %d/%d, want 1/%d", oldCloses, newCloses, wantNew)
			}
		})
	}
}

func TestInvalidRegistrationPreservesOwnership(t *testing.T) {
	var manager resource.Manager
	closes := 0

	if err := manager.Own(resource.FileSystem, func() error {
		closes++

		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err := manager.Own(resource.FileSystem, nil); err == nil {
		t.Fatal("expected nil cleanup to be rejected")
	}

	if closes != 0 {
		t.Fatal("invalid replacement closed the current resource")
	}

	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}

	if closes != 1 {
		t.Fatalf("current resource closed %d times, want 1", closes)
	}

	if err := manager.Own(resource.FileSystem, func() error {
		closes++

		return nil
	}); err == nil {
		t.Fatal("expected ownership after closure to be rejected")
	}

	if err := manager.Borrow(resource.Network); err == nil {
		t.Fatal("expected borrowing after closure to be rejected")
	}

	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}

	if closes != 1 {
		t.Fatal("manager closed a rejected resource")
	}
}

func TestCloseEmptyManager(t *testing.T) {
	var manager resource.Manager

	for range 2 {
		if err := manager.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
