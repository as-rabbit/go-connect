package water_mill

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	event_interface "github.com/as-rabbit/go-connect/event_driven/events/event-interface"
)

type WaterMillKafka struct {
	sub        *kafka.Subscriber
	pub        *kafka.Publisher
	chanBuffer int64
}

func NewWaterMillKafka(brokers []string, config *sarama.Config, consumerGroup string, adapter watermill.LoggerAdapter) event_interface.Event {
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: config,
			ConsumerGroup:         consumerGroup,
		},
		adapter,
	)

	if err != nil {
		panic(err)
	}

	pub, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   brokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}

	return &WaterMillKafka{
		sub: sub,
		pub: pub,
	}
}

func (w *WaterMillKafka) Subscribe(ctx context.Context, topic string) (<-chan event_interface.Message, error) {
	c, err := w.sub.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	messageChan := make(chan event_interface.Message, w.chanBuffer)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-c:
				messageChan <- &WaterMillMessage{msg: m}
			}
		}
	}()

	return messageChan, err
}

func (w *WaterMillKafka) Publish(ctx context.Context, topic string, messages ...event_interface.Message) error {
	msgList := make([]*message.Message, len(messages))
	for i := range messages {
		msgList[i] = message.NewMessage("", messages[i].GetPayload())
	}
	return w.pub.Publish(topic, msgList...)
}
