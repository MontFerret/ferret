package debugger

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSessionInspectionAndBreakpointContexts(t *testing.T) {
	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	expired, stop := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	defer stop()

	for name, operation := range map[string]func(*Session, context.Context) error{
		"pause": func(s *Session, ctx context.Context) error { return s.Pause(ctx) },
		"add": func(s *Session, ctx context.Context) error {
			_, err := s.SetBreakpoint(ctx, source.Location{Position: source.Position{Line: 1}})

			return err
		},
		"add_at": func(s *Session, ctx context.Context) error {
			_, err := s.SetBreakpointAt(ctx, source.Location{Position: source.Position{Line: 1}}, BreakpointOptions{})

			return err
		},
		"delete": func(s *Session, ctx context.Context) error { return s.DeleteBreakpoint(ctx, 1) },
		"list": func(s *Session, ctx context.Context) error {
			values, err := s.Breakpoints(ctx)
			if values != nil {
				t.Errorf("failed listing returned values: %+v", values)
			}

			return err
		},
		"frames": func(s *Session, ctx context.Context) error {
			_, err := s.Frames(ctx)

			return err
		},
		"locals": func(s *Session, ctx context.Context) error {
			_, err := s.Locals(ctx)

			return err
		},
		"frame_locals": func(s *Session, ctx context.Context) error {
			_, err := s.FrameLocals(ctx, 0)

			return err
		},
		"variables": func(s *Session, ctx context.Context) error {
			_, err := s.Variables(ctx, 1)

			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				ctx  context.Context
				want error
				name string
			}{
				{name: "nil", ctx: nil, want: runtime.ErrInvalidArgument},
				{name: "canceled", ctx: canceled, want: context.Canceled},
				{name: "expired", ctx: expired, want: context.DeadlineExceeded},
			} {
				t.Run(tc.name, func(t *testing.T) {
					session := replacementSession(t)
					for _, receiver := range []*Session{session, nil} {
						if err := operation(receiver, tc.ctx); !errors.Is(err, tc.want) {
							t.Fatalf("error = %v, want %v", err, tc.want)
						}
					}
				})
			}
		})
	}
}

func TestInspectionRechecksCancellationAfterCommandWait(t *testing.T) {
	for name, inspect := range map[string]func(*Session, context.Context) error{
		"frames": func(s *Session, ctx context.Context) error {
			_, err := s.Frames(ctx)

			return err
		},
		"locals": func(s *Session, ctx context.Context) error {
			_, err := s.Locals(ctx)

			return err
		},
		"frame_locals": func(s *Session, ctx context.Context) error {
			_, err := s.FrameLocals(ctx, 0)

			return err
		},
		"variables": func(s *Session, ctx context.Context) error {
			_, err := s.Variables(ctx, 1)

			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			session := replacementSession(t)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			checked := make(chan struct{}, 1)
			observed := &contextObserver{Context: ctx, check: func() {
				select {
				case checked <- struct{}{}:
				default:
				}
			}}

			done := make(chan error, 1)
			func() {
				session.lifecycle.commandMu.Lock()
				defer session.lifecycle.commandMu.Unlock()
				go func() { done <- inspect(session, observed) }()
				waitForSignal(t, checked, "initial context check")
				cancel()
				assertBlocked(t, done, "inspection command admission")
			}()

			if err := receiveLiveError(t, done); !errors.Is(err, context.Canceled) {
				t.Fatalf("inspection entered after cancellation: %v", err)
			}
		})
	}
}

func TestCanceledPauseDoesNotInterruptExecution(t *testing.T) {
	session, execution := newBlockingSession(t)
	defer session.Close()
	if _, err := session.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		_, err := session.Continue(t.Context())
		done <- err
	}()
	waitForSignal(t, execution.resumeStarted, "resume")

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := session.Pause(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}

	select {
	case <-execution.pauseCalled:
		t.Fatal("canceled request paused execution")
	default:
	}

	assertBlocked(t, done, "execution after canceled pause")
	if err := session.Pause(t.Context()); err != nil {
		t.Fatal(err)
	}

	waitForError(t, done, "valid pause")
}
