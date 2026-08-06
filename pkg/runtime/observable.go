package runtime

import (
	"context"
	"io"
	"time"
)

// DefaultStreamTimeout is the exact timeout used by NewStreamIterator.
const DefaultStreamTimeout = 5 * time.Second

type (
	// Subscription represents an event subscription object that contains target event name
	// and optional event options.
	Subscription struct {
		Options   Map    `json:"options"`
		EventName String `json:"eventName"`
	}

	// Message represents an event message that an Observable can emit.
	Message interface {
		// Value returns a value of the message.
		Value() Value
		// Err returns an error of the message.
		Err() error
	}

	defaultMessage struct {
		value Value
		err   error
	}

	// Stream represents an event stream that produces target event objects.
	Stream interface {
		io.Closer
		// Read must observe ctx while waiting for messages.
		Read(ctx context.Context) <-chan Message
	}

	// Observable represents an interface of
	// complex types that returns stream of events. Subscribe implementations that
	// may block must observe ctx while establishing the subscription.
	Observable interface {
		Subscribe(ctx context.Context, subscription Subscription) (Stream, error)
	}
)

// NewErrorMessage creates a new Message with an error.
func NewErrorMessage(err error) Message {
	return &defaultMessage{err: err, value: None}
}

// NewValueMessage creates a new Message with a value.
func NewValueMessage(val Value) Message {
	return &defaultMessage{err: nil, value: val}
}

func (m *defaultMessage) Value() Value {
	return m.value
}

func (m *defaultMessage) Err() error {
	return m.err
}
