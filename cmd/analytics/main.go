package main

import (
	"context"
	"fmt"
	"log/slog"
	"short-links/internal/config"

	"github.com/segmentio/kafka-go"
)

func main() {
	cfg := config.MustConfigLoad()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaAddr},
		Topic:   cfg.KafkaAnalyticTopic,
		GroupID: "analytics-service",
	})

	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			slog.Info("failed to read message kafka", "error", err)
			continue
		}

		fmt.Println(msg.Topic, msg.Offset, string(msg.Value))
	}
}
