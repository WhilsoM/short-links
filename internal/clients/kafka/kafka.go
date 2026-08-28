package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafkaClient(addr, topic string) *kafka.Writer {

	writer := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Topic:        topic,
		BatchTimeout: 10 * time.Millisecond,
	}

	return writer
}
