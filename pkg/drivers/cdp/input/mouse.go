package input

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
	"time"
)

type Mouse struct {
	client *cdp.Client
	x      float64
	y      float64
}

func NewMouse(client *cdp.Client) *Mouse {
	return &Mouse{client, 0, 0}
}

func (m *Mouse) Click(ctx context.Context, x, y float64, delay int) error {
	if err := m.Move(ctx, x, y); err != nil {
		return err
	}

	if err := m.Down(ctx, "left"); err != nil {
		return err
	}

	max, min := core.NumberBoundaries(float64(delay))
	releaseDelay := time.Duration(core.Random(max, min))

	time.Sleep(releaseDelay * time.Millisecond)

	return m.Up(ctx, "left")
}

func (m *Mouse) Down(ctx context.Context, button string) error {
	return m.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mousePressed", m.x, m.y).
			SetClickCount(1).
			SetButton(button),
	)
}

func (m *Mouse) Up(ctx context.Context, button string) error {
	return m.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseReleased", m.x, m.y).
			SetClickCount(1).
			SetButton(button),
	)
}

func (m *Mouse) Move(ctx context.Context, x, y float64) error {
	return m.MoveBySteps(ctx, x, y, 1)
}

func (m *Mouse) MoveBySteps(ctx context.Context, x, y float64, steps int) error {
	fromX := m.x
	fromY := m.y

	for i := 0; i <= steps; i++ {
		iFloat := float64(i)
		stepFloat := float64(steps)
		toX := fromX + (x-fromX)*(iFloat/stepFloat)
		toY := fromY + (y-fromY)*(iFloat/stepFloat)

		err := m.client.Input.DispatchMouseEvent(
			ctx,
			input.NewDispatchMouseEventArgs("mouseMoved", toX, toY),
		)

		if err != nil {
			return err
		}
	}

	m.x = x
	m.y = y

	return nil
}
