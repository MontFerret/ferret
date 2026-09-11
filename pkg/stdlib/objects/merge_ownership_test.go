package objects_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func TestImmutableMergeDestinationOwnership(t *testing.T) {
	for name, merge := range map[string]runtime.Function{"shallow": objects.Merge, "deep": objects.MergeDeep} {
		for _, listForm := range []bool{false, true} {
			for _, cleanupFails := range []bool{false, true} {
				for _, stage := range []string{"success", "set", "traversal", "clone", "lookup", "nested", "factory", "pre-cancel", "new-cancel", "set-cancel", "traversal-cancel", "failure-cancel"} {
					if (stage == "lookup" || stage == "nested") && name != "deep" {
						continue
					}

					t.Run(fmt.Sprintf("%s/list=%v/cleanup=%v/%s", name, listForm, cleanupFails, stage), func(t *testing.T) {
						ctx, cancel := context.WithCancel(t.Context())
						defer cancel()
						primary, cleanup := errors.New("population failed"), errors.New("destination close failed")
						first := &configuredMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"existing": runtime.Int(1)}), backend: "configured storage"}
						second := &configuredMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.Int(2)}), backend: "other storage"}
						first.onNew = func(dst *configuredMap) {
							if cleanupFails {
								dst.closeErr = cleanup
							}

							if stage == "new-cancel" {
								cancel()
							}

							if stage == "lookup" || stage == "nested" {
								dst.onLookup = func(c context.Context, _ runtime.Value) (runtime.Value, bool, error) {
									if c != ctx {
										t.Fatal("lookup context changed")
									}

									if stage == "lookup" {
										return runtime.None, false, primary
									}

									return &hostMap{Map: runtime.NewObject(), setErr: primary}, true, nil
								}
							}

							dst.onSet = func(c context.Context, key, value runtime.Value) error {
								if c != ctx {
									t.Fatal("mutation context changed")
								}

								if key == runtime.String("key") {
									if stage == "set" || stage == "failure-cancel" {
										if stage == "failure-cancel" {
											cancel()
										}

										return primary
									}
								}

								err := dst.Object.Set(c, key, value)
								if key == runtime.String("key") && stage == "set-cancel" {
									cancel()
								}

								return err
							}
						}

						switch stage {
						case "factory":
							first.newErr = primary
						case "pre-cancel":
							cancel()
						case "clone":
							second.Object = runtime.NewObjectWith(map[string]runtime.Value{"key": &cloneValue{Value: runtime.Int(2), clones: new(int), cloneErr: primary}})
						case "lookup", "nested":
							second.Object = runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.NewObjectWith(map[string]runtime.Value{"nested": runtime.Int(2)})})
						case "traversal", "traversal-cancel":
							second.onForEach = func(c context.Context, predicate runtime.KeyReadablePredicate) error {
								if err := second.Object.ForEach(c, predicate); err != nil {
									return err
								}

								if stage == "traversal-cancel" {
									cancel()

									return nil
								}

								return primary
							}
						}

						before := second.String()
						args := []runtime.Value{first, second}
						if listForm {
							args = []runtime.Value{runtime.NewArrayWith(args...)}
						}

						result, err := merge(ctx, args...)
						if first.closes != 0 || second.closes != 0 || second.creations != 0 || first.String() != `{"existing":1}` || second.String() != before {
							t.Fatal("borrowed sources changed or closed")
						}

						if stage == "pre-cancel" {
							if result != runtime.None || !errors.Is(err, context.Canceled) || first.created != nil {
								t.Fatalf("pre-cancel: %v, %v", result, err)
							}

							return
						}

						dst := first.created
						if dst == nil || first.creations != 1 || dst.backend != first.backend {
							t.Fatal("factory result or configuration lost")
						}

						if stage == "new-cancel" && dst.String() != "{}" {
							t.Fatal("population continued after construction canceled")
						}

						if stage == "success" {
							if err != nil || result != dst || dst.closes != 0 {
								t.Fatalf("successful ownership transfer: %v, %v, closes=%d", result, err, dst.closes)
							}

							assertObjectValue(t, dst, runtime.NewObjectWith(map[string]runtime.Value{"existing": runtime.Int(1), "key": runtime.Int(2)}))

							closeErr := dst.Close()
							if dst.closes != 1 || errors.Is(closeErr, cleanup) != cleanupFails {
								t.Fatal("caller cleanup changed")
							}

							return
						}

						wantPrimary := !strings.HasSuffix(stage, "-cancel") || stage == "failure-cancel"
						wantCancel := strings.HasSuffix(stage, "-cancel")
						if result != runtime.None || err == nil || dst.closes != 1 || errors.Is(err, primary) != wantPrimary || errors.Is(err, context.Canceled) != wantCancel || errors.Is(err, cleanup) != cleanupFails {
							t.Fatalf("failure ownership: %v, %v, closes=%d", result, err, dst.closes)
						}

						index := 1
						if stage == "factory" || stage == "new-cancel" {
							index = 0
						}

						attribution := fmt.Sprintf("position %d", index+1)
						if listForm {
							attribution = fmt.Sprintf("item %d", index)
							if !strings.Contains(err.Error(), "position 1") {
								t.Fatal("list argument attribution lost")
							}
						}

						if !strings.Contains(err.Error(), attribution) {
							t.Fatalf("source attribution lost: %v", err)
						}
					})
				}
			}
		}
	}
}

func TestMutableMergeBorrowsClosableTarget(t *testing.T) {
	for name, merge := range map[string]runtime.Function{"shallow": objects.MergeMutable, "deep": objects.MergeDeepMutable} {
		for _, stage := range []string{"success", "failure", "cancel"} {
			t.Run(name+"/"+stage, func(t *testing.T) {
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				primary := errors.New("mutation failed")
				if stage == "cancel" {
					primary = context.Canceled
				}

				target := &configuredMap{Object: runtime.NewObject()}
				source := &configuredMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"key": runtime.Int(1)})}
				if stage != "success" {
					target.onSet = func(ctx context.Context, key, value runtime.Value) error {
						if err := target.Object.Set(ctx, key, value); err != nil {
							return err
						}

						if stage == "cancel" {
							cancel()
						}

						return primary
					}
				}

				result, err := merge(ctx, target, source)
				if target.closes != 0 || source.closes != 0 || target.creations != 0 || source.creations != 0 || target.String() != `{"key":1}` {
					t.Fatal("borrowed target ownership or partial mutation changed")
				}

				if stage != "success" && (result != runtime.None || !errors.Is(err, primary)) || stage == "success" && (result != target || err != nil) {
					t.Fatalf("result=%v, err=%v", result, err)
				}
			})
		}
	}
}
