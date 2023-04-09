package water_mill

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"testing"
	"time"
)

func TestWaterMillKafa(t *testing.T) {
	config := kafka.DefaultSaramaSubscriberConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	k := NewWaterMillKafka([]string{"www.wxxserver.com:9092"}, config, watermill.NewStdLogger(false, false))
	ctx, cancel := context.WithCancel(context.Background())
	topic := "wxx_test"

	c2, err := k.Subscribe(ctx, topic, "t2")
	if err != nil {
		t.Fatal(err)
	}

	c, err := k.Subscribe(ctx, topic, "t1")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case m := <-c:
				t.Logf("recive msg: %s", m.GetPayload())
				m.Ack()
				//wait.Done()
			case m2 := <-c2:
				t.Logf("recive msg2: %s", m2.GetPayload())
				m2.Ack()
				//wait.Done()
			}
		}
	}()

	t.Log("start")
	go func() {
		i := 0
		for {
			select {
			case <-ctx.Done():
				break
			default:
				msg := NewWaterMillMessage(fmt.Sprintf("uuid_%d", i), []byte(fmt.Sprintf("Hello World! %d", i)))
				t.Log(k.Publish(ctx, topic, msg))
				i++
			}
		}
	}()

	time.Sleep(30 * time.Second)
	cancel()
	t.Log("end")
}
