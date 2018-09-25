package browser

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/page"
)

type (
	EventHandler func(event, message string)

	EventBroker struct {
		client   page.LifecycleEventClient
		handlers map[string][]EventHandler
		cancel   context.CancelFunc
	}
)

func NewEventBroker(client page.LifecycleEventClient) *EventBroker {
	return &EventBroker{
		client,
		make(map[string][]EventHandler),
		nil,
	}
}

func (broker *EventBroker) Start() error {
	if broker.cancel != nil {
		return core.Error(core.ErrInvalidOperation, "broker is already started")
	}

	ctx, cancel := context.WithCancel(context.Background())

	broker.cancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-broker.client.Ready():
				reply, err := broker.client.Recv()

				if err != nil {
					fmt.Println("FAILED TO GET EVENT", err)
					broker.Emit("error", err.Error())

					return
				}

				fmt.Println("EVENT", reply.Name)

				broker.Emit(reply.Name, "")
			}
		}
	}()

	return nil
}

func (broker *EventBroker) Stop() error {
	if broker.cancel == nil {
		return core.Error(core.ErrInvalidOperation, "broker is already stopped")
	}

	broker.cancel()
	broker.client = nil

	return nil
}

func (broker *EventBroker) Close() error {
	if broker.cancel != nil {
		broker.Stop()
	}

	return broker.client.Close()
}

func (broker *EventBroker) AddListener(event string, handler EventHandler) {
	handlers, ok := broker.handlers[event]

	if !ok {
		handlers = make([]EventHandler, 0, 5)

		broker.handlers[event] = handlers
	}

	handlers = append(handlers, handler)
}

func (broker *EventBroker) Emit(name, message string) {
	handlers, ok := broker.handlers[name]

	if !ok {
		return
	}

	for _, handler := range handlers {
		handler(name, message)
	}
}
