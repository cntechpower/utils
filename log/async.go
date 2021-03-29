package log

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func WithKafka(appId, kafkaAddr, topic string) Option {
	return newLogOption(func(option *logOptions) {
		l := newKafkaWriter(appId, kafkaAddr, topic)
		lc := &loggerWithConfig{
			typ:    OutputTypeJson,
			buffer: make(chan string, 1000),
			Logger: l,
		}
		loggers = append(loggers, lc)
		go lc.run()
	})
}

type AsyncLogging struct {
	App string
	Msg string
}

type kafkaWriter struct {
	appId     string
	kafkaAddr string
	topic     string
	writer    *kafka.Writer
}

func newKafkaWriter(appId, kafkaAddr, topic string) *kafkaWriter {
	w := &kafkaWriter{
		appId:     appId,
		kafkaAddr: kafkaAddr,
		topic:     topic,
	}
	w.writer = &kafka.Writer{
		Addr:        kafka.TCP(kafkaAddr),
		Topic:       topic,
		Async:       false,
		Logger:      nil,
		ErrorLogger: log.New(os.Stderr, "", log.LstdFlags),
		Compression: kafka.Lz4,
		Balancer:    &kafka.LeastBytes{},
	}

	return w
}
func (w *kafkaWriter) Println(v ...interface{}) {
	if len(v) != 1 {
		fmt.Println("kafkaWriter Println got multi v")
		return
	}
	s, ok := v[0].(string)
	if !ok {
		fmt.Println("kafkaWriter Println got non string")
		return
	}
	var err error
	msg := &AsyncLogging{
		App: w.appId,
		Msg: s,
	}
	msgBytes, _ := json.Marshal(msg)
	for i := 0; i < 3; i++ {
		err = w.writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(fmt.Sprintf("%v-%v", w.appId, time.Now().UnixNano())),
			Value: msgBytes,
		})
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Printf("send kafka error: %v\n", err)
	}
}
