package input

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
	"time"
)

type Keyboard struct {
	client *cdp.Client
}

func NewKeyboard(client *cdp.Client) *Keyboard {
	return &Keyboard{client}
}

func (k *Keyboard) Down(ctx context.Context, char string) error {
	return k.client.Input.DispatchKeyEvent(
		ctx,
		input.NewDispatchKeyEventArgs("keyDown").
			SetText(char),
	)
}

func (k *Keyboard) Up(ctx context.Context, char string) error {
	return k.client.Input.DispatchKeyEvent(
		ctx,
		input.NewDispatchKeyEventArgs("keyUp").
			SetText(char),
	)
}

func (k *Keyboard) Type(ctx context.Context, text string, delay int) error {
	for _, ch := range text {
		ch := string(ch)

		if err := k.Down(ctx, ch); err != nil {
			return err
		}

		max, min := core.NumberBoundaries(float64(delay))
		releaseDelay := time.Duration(core.Random(max, min))
		time.Sleep(releaseDelay)

		if err := k.Up(ctx, ch); err != nil {
			return err
		}
	}

	return nil
}
