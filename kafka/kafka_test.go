package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"

	"github.com/stretchr/testify/assert"

	"github.com/cntechpower/utils/tracing"
)

func TestTracingKafka(t *testing.T) {
	tracing.Init("unit-test", "10.0.0.2:6831")
	defer tracing.Close()

	kafka, err := NewSyncProducer([]string{"10.0.0.2:9093"})
	if !assert.Equal(t, nil, err) {
		assert.FailNow(t, "connect to kafka error: %v", err)
	}
	defer func() {
		_ = kafka.Close()
	}()

	_, _, err = kafka.SendMessage(context.Background(), &sarama.ProducerMessage{
		Topic:     "test-topic",
		Value:     sarama.StringEncoder("hello-world"),
		Timestamp: time.Now(),
	})
	assert.Equal(t, nil, err)
}
