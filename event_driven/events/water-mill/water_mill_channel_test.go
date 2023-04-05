package water_mill

import (
	"context"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"sync"
	"testing"
)

func TestWaterMillChannel(t *testing.T) {
	topic := "test_topic"
	c := NewWaterMillChannel(gochannel.Config{OutputChannelBuffer: 1}, watermill.NewStdLogger(false, false))

	ctx, cancel := context.WithCancel(context.Background())
	msgs, err := c.Subscribe(ctx, topic, "test_group")
	if err != nil {
		t.Fatal(err)
	}

	msgs2, err := c.Subscribe(ctx, topic, "test_group")
	if err != nil {
		t.Fatal(err)
	}

	wait := sync.WaitGroup{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case msg := <-msgs:
				t.Logf("recive msg: %s", msg.GetPayload())
				msg.Ack()
			case msg2 := <-msgs2:
				t.Logf("recive msg2: %s", msg2.GetPayload())
				msg2.Ack()
			}
			wait.Done()
		}
	}()

	wait.Add(1)
	for i := 0; i < 10; i++ {
		wait.Add(1)
		msg := NewWaterMillMessage(fmt.Sprintf("uuid_%d", i), []byte(fmt.Sprintf("Hello World! %d", i)))
		c.Publish(ctx, topic, msg)
		fmt.Println("i:", i)
	}
	wait.Done()

	wait.Wait()
	cancel()
}
