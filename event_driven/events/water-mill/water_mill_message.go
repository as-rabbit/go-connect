package water_mill

import "github.com/ThreeDotsLabs/watermill/message"

type WaterMillMessage struct {
	msg *message.Message
}

func NewWaterMillMessage(uuid string, payLoad []byte) *WaterMillMessage {
	return &WaterMillMessage{
		msg: message.NewMessage(uuid, payLoad),
	}
}

func (w *WaterMillMessage) Ack() bool {
	return w.msg.Ack()
}

func (w *WaterMillMessage) GetPayload() []byte {
	return w.msg.Payload
}

func (w *WaterMillMessage) GetTraceId() string {
	return w.msg.UUID
}
