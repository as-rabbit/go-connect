package event_driven

import "github.com/as-rabbit/go-connect/event_driven/events/event-interface"

type EventType string

const (
	WaterMillGoChannel EventType = "water_mill_go_channel"
)

func NewEvent(t EventType) event_interface.Event {
	switch t {
	}

	return nil
}
