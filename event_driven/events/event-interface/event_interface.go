package event_interface

import "context"

type Message interface {
	GetPayload() []byte
	Ack() bool
	GetTraceId() string
}

type Event interface {
	Subscribe(ctx context.Context, topic, groupId string) (<-chan Message, error)
	Publish(ctx context.Context, topic string, messages ...Message) error
}
