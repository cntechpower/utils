package kafka

import (
	"context"

	"github.com/cntechpower/utils/tracing"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pkg/errors"
)

const (
	dbTypeKafka = "kafka"
)

type SyncProducer struct {
	sarama.SyncProducer
}

func NewSyncProducer(brokerList []string) (producer *SyncProducer, err error) {
	p, err := newSyncProducer(brokerList)
	if err != nil {
		return
	}
	return &SyncProducer{p}, nil
}

func (p *SyncProducer) SendMessage(ctx context.Context, msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	span, _ := tracing.New(ctx, "kafka.SendMessage")
	ext.DBType.Set(span, dbTypeKafka)
	span.SetTag("topic", msg.Topic)
	partition, offset, err = p.SyncProducer.SendMessage(msg)
	span.SetTag("partition", partition)
	span.SetTag("offset", offset)
	if err != nil {
		ext.LogError(span, err)
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func (p *SyncProducer) SendMessages(ctx context.Context, msgs []*sarama.ProducerMessage) (err error) {
	span, _ := tracing.New(ctx, "kafka.SendMessages")
	ext.DBType.Set(span, dbTypeKafka)
	span.SetTag("msg_count", len(msgs))
	err = p.SyncProducer.SendMessages(msgs)
	if err != nil {
		ext.LogError(span, err)
		ext.Error.Set(span, true)
	}
	span.Finish()
	return
}

func newSyncProducer(brokerList []string) (producer sarama.SyncProducer, err error) {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err = sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		err = errors.Wrapf(err, "Failed to start Sarama producer: %+v", err)
		return
	}

	return
}
