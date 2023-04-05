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
	brokers      []string
	pub          *kafka.Publisher
	chanBuffer   int64
	saramaConfig *sarama.Config
	logAdapter   watermill.LoggerAdapter
}

func NewWaterMillKafka(brokers []string, config *sarama.Config, adapter watermill.LoggerAdapter) event_interface.Event {

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
		brokers:      brokers,
		pub:          pub,
		saramaConfig: config,
		logAdapter:   adapter,
	}
}

func (w *WaterMillKafka) Subscribe(ctx context.Context, topic, groupId string) (<-chan event_interface.Message, error) {
	sub, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               w.brokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: w.saramaConfig,
			ConsumerGroup:         groupId,
		},
		w.logAdapter,
	)

	c, err := sub.Subscribe(ctx, topic)
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
