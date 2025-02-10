package utils

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// AppConfig holds all system-wide configurations
type AppConfig struct {
	Httpserver HttpServerConfig
	Kafka  KafkaConfig
	Logger LoggerConfig
}

// ServerConfig - HTTP server related configs
type HttpServerConfig struct {
	Host string `env:"HTTP_HOST" envDefault:"localhost"`
	Port int    `env:"HTTP_PORT" envDefault:"8080"`
}

// KafkaConfig - Holds Kafka settings for producer & consumer
type KafkaConfig struct {
	Brokers            string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	ClientID           string `env:"KAFKA_CLIENT_ID" envDefault:"default-client"`
	GroupID            string `env:"KAFKA_CONSUMER_GROUP_ID" envDefault:"default-group"`
	AutoOffsetReset    string `env:"KAFKA_CONSUMER_AUTO_OFFSET_RESET" envDefault:"earliest"`
	EnableAutoCommit   bool   `env:"KAFKA_CONSUMER_ENABLE_AUTO_COMMIT" envDefault:"false"`
	MaxPollIntervalMs  int    `env:"KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS" envDefault:"300000"`
	SessionTimeoutMs   int    `env:"KAFKA_CONSUMER_SESSION_TIMEOUT_MS" envDefault:"45000"`
	HeartbeatIntervalMs int   `env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS" envDefault:"3000"`
	RetryBackoffMs     int    `env:"KAFKA_CONSUMER_RETRY_BACKOFF_MS" envDefault:"100"`
	FetchMinBytes      int    `env:"KAFKA_CONSUMER_FETCH_MIN_BYTES" envDefault:"1"`
	FetchWaitMaxMs     int    `env:"KAFKA_CONSUMER_FETCH_WAIT_MAX_MS" envDefault:"500"`
	SchemaRegistryURL  string `env:"KAFKA_SCHEMA_REGISTRY_URL" envDefault:"http://localhost:8081"`
}

// LoggerConfig - Logger settings
type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

// LoadConfig loads the full app configuration from environment variables
func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}
	// Parse environment variables
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return nil, err
	}

	log.Println("Configuration successfully loaded")
	return cfg, nil
}
