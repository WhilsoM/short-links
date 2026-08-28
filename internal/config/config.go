package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string `env:"PORT" env-default:":8080"`
	DatabaseURL        string `env:"DATABASE_URL" env-required:"true"`
	JWTSecret          string `env:"JWT_SECRET" env-required:"true"`
	RedisAddr          string `env:"REDIS_ADDR" env-required:"true"`
	KafkaAddr          string `env:"KAFKA_ADDR" env-required:"true"`
	KafkaAnalyticTopic string `env:"KAFKA_ANALYTIC_TOPIC" env-required:"true"`
}

func MustConfigLoad() *Config {
	_ = godotenv.Load()

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	return &cfg
}
