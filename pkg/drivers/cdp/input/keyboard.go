package input

import (
	"context"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
)

type (
	KeyboardModifier int

	KeyboardLocation int

	KeyboardKey struct {
		KeyCode  int
		Key      string
		Code     string
		Modifier KeyboardModifier
		Location KeyboardLocation
	}

	Keyboard struct {
		client *cdp.Client
	}
)

const (
	KeyboardModifierNone  KeyboardModifier = 0
	KeyboardModifierAlt   KeyboardModifier = 1
	KeyboardModifierCtrl  KeyboardModifier = 2
	KeyboardModifierCmd   KeyboardModifier = 4
	KeyboardModifierShift KeyboardModifier = 8

	// 1=Left, 2=Right
	KeyboardLocationNone  KeyboardLocation = 0
	KeyboardLocationLeft  KeyboardLocation = 1
	KeyboardLocationRight KeyboardLocation = 2
)

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

func (k *Keyboard) Type(ctx context.Context, text string, delay time.Duration) error {
	for _, ch := range text {
		ch := string(ch)

		if err := k.Down(ctx, ch); err != nil {
			return err
		}

		releaseDelay := randomDuration(int(delay))
		time.Sleep(releaseDelay)

		if err := k.Up(ctx, ch); err != nil {
			return err
		}
	}

	return nil
}

func (k *Keyboard) Press(ctx context.Context, keys []string, count int, delay time.Duration) error {
	for i := 0; i < count; i++ {
		if i > 0 {
			downDelay := randomDuration(int(delay))
			time.Sleep(downDelay)
		}

		if err := k.press(ctx, keys, delay); err != nil {
			return err
		}
	}

	return nil
}

func (k *Keyboard) press(ctx context.Context, keys []string, delay time.Duration) error {
	for i, key := range keys {
		if i > 0 {
			downDelay := randomDuration(int(delay))
			time.Sleep(downDelay)
		}

		if err := k.client.Input.DispatchKeyEvent(
			ctx,
			k.createPressEvent("keyDown", key),
		); err != nil {
			return err
		}
	}

	for _, key := range keys {
		upDelay := randomDuration(int(delay))
		time.Sleep(upDelay)

		if err := k.client.Input.DispatchKeyEvent(
			ctx,
			k.createPressEvent("keyUp", key),
		); err != nil {
			return err
		}
	}

	return nil
}

func (k *Keyboard) createPressEvent(event string, chars string) *input.DispatchKeyEventArgs {
	args := input.NewDispatchKeyEventArgs(event)

	key, found := usKeyboardLayout[chars]

	if found {
		args.
			SetCode(key.Code).
			SetKey(key.Key).
			SetModifiers(int(key.Modifier)).
			SetWindowsVirtualKeyCode(key.KeyCode)
	}

	return args
}
