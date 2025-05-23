package events

import "errors"

var ErrorHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handlers, ok := ed.handlers[event.GetName()]; ok {
		for _, handler := range handlers {
			handler.Handle(event)
		}
	}
	return nil
}

func (ed *EventDispatcher) GetHandlers() map[string][]EventHandlerInterface {
	return ed.handlers
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if handlers, ok := ed.handlers[eventName]; ok {
		for _, h := range handlers {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrorHandlerAlreadyRegistered
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if handlers, ok := ed.handlers[eventName]; ok {
		for i, h := range handlers {
			if h == handler {
				ed.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
				return nil
			}
		}
	}
	return nil
}
