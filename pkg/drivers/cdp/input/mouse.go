package input

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
)

type Mouse struct {
	client *cdp.Client
	x      float64
	y      float64
}

func NewMouse(client *cdp.Client) *Mouse {
	return &Mouse{client, 0, 0}
}

func (m *Mouse) Click(ctx context.Context, x, y float64, delay time.Duration) error {
	return m.ClickWithCount(ctx, x, y, delay, 1)
}

func (m *Mouse) ClickWithCount(ctx context.Context, x, y float64, delay time.Duration, count int) error {
	if err := m.Move(ctx, x, y); err != nil {
		return err
	}

	if err := m.DownWithCount(ctx, input.MouseButtonLeft, count); err != nil {
		return err
	}

	time.Sleep(randomDuration(int(delay)))

	return m.UpWithCount(ctx, input.MouseButtonLeft, count)
}

func (m *Mouse) Down(ctx context.Context, button input.MouseButton) error {
	return m.DownWithCount(ctx, button, 1)
}

func (m *Mouse) DownWithCount(ctx context.Context, button input.MouseButton, count int) error {
	return m.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mousePressed", m.x, m.y).
			SetButton(button).
			SetClickCount(count),
	)
}

func (m *Mouse) Up(ctx context.Context, button input.MouseButton) error {
	return m.UpWithCount(ctx, button, 1)
}

func (m *Mouse) UpWithCount(ctx context.Context, button input.MouseButton, count int) error {
	return m.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseReleased", m.x, m.y).
			SetButton(button).
			SetClickCount(count),
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
