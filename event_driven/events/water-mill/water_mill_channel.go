package water_mill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/as-rabbit/go-connect/event_driven/events/event-interface"
	"sync"
)

type WaterMillChannel struct {
	c          *gochannel.GoChannel
	chanBuffer int64
	groupMap   map[string]chan event_interface.Message
	subMux     sync.Mutex
}

func NewWaterMillChannel(config gochannel.Config, adapter watermill.LoggerAdapter) event_interface.Event {
	return &WaterMillChannel{
		c:          gochannel.NewGoChannel(config, adapter),
		chanBuffer: config.OutputChannelBuffer,
		groupMap:   make(map[string]chan event_interface.Message, 0),
	}
}

func (w *WaterMillChannel) Subscribe(ctx context.Context, topic, groupId string) (<-chan event_interface.Message, error) {
	w.subMux.Lock()
	defer w.subMux.Unlock()

	if msgChan, ok := w.groupMap[groupId]; ok {
		return msgChan, nil
	}

	c, err := w.c.Subscribe(ctx, topic)
	if err != nil {
		return nil, err
	}

	w.groupMap[groupId] = make(chan event_interface.Message, w.chanBuffer)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-c:
				w.groupMap[groupId] <- &WaterMillMessage{msg: m}
			}
		}

	}()

	return w.groupMap[groupId], err
}

func (w *WaterMillChannel) Publish(ctx context.Context, topic string, messages ...event_interface.Message) error {
	msgList := make([]*message.Message, len(messages))
	for i := range messages {
		msgList[i] = message.NewMessage(messages[i].GetTraceId(), messages[i].GetPayload())
	}
	return w.c.Publish(topic, msgList...)
}
