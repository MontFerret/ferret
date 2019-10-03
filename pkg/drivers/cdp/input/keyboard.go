package input

import (
	"context"
	"github.com/pkg/errors"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
)

const DefaultDelay = 25

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

func (k *Keyboard) Press(ctx context.Context, name string) error {
	key, found := usKeyboardLayout[name]

	if !found {
		return errors.New("invalid key")
	}

	err := k.client.Input.DispatchKeyEvent(
		ctx,
		input.NewDispatchKeyEventArgs("keyDown").
			SetCode(key.Code).
			SetKey(key.Key).
			SetWindowsVirtualKeyCode(key.KeyCode),
	)

	if err != nil {
		return err
	}

	releaseDelay := randomDuration(DefaultDelay)
	time.Sleep(releaseDelay)

	return k.client.Input.DispatchKeyEvent(
		ctx,
		input.NewDispatchKeyEventArgs("keyUp").
			SetCode(key.Code).
			SetKey(key.Key).
			SetWindowsVirtualKeyCode(key.KeyCode),
	)
}
