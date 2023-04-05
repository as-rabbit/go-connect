package water_mill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/as-rabbit/go-connect/event_driven/events/event-interface"
)

type WaterMillChannel struct {
	c          *gochannel.GoChannel
	chanBuffer int64
}

func NewWaterMillChannel(config gochannel.Config, adapter watermill.LoggerAdapter) event_interface.Event {
	return &WaterMillChannel{
		c:          gochannel.NewGoChannel(config, adapter),
		chanBuffer: config.OutputChannelBuffer,
	}
}

func (w *WaterMillChannel) Subscribe(ctx context.Context, topic string) (<-chan event_interface.Message, error) {
	c, err := w.c.Subscribe(ctx, topic)
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

func (w *WaterMillChannel) Publish(ctx context.Context, topic string, messages ...event_interface.Message) error {
	msgList := make([]*message.Message, len(messages))
	for i := range messages {
		msgList[i] = message.NewMessage(messages[i].GetTraceId(), messages[i].GetPayload())
	}
	return w.c.Publish(topic, msgList...)
}
