package kafka

import (
	"log"
	"proposal-template/pkg/utils/config"
)

// KafkaProducerConfig holds Kafka producer settings
type KafkaProducerConfig struct {
	Brokers  string `env:"KAFKA_BROKERS,required"`
	ClientID string `env:"KAFKA_CLIENT_ID" envDefault:"default-client"`
}

var _ utils.Config = (*KafkaProducerConfig)(nil)

func (c *KafkaProducerConfig) Load() error {
	log.Printf("Loading KafkaProducerConfig")
	return utils.ParseConfig(c)
}

// KafkaConsumerConfig holds Kafka consumer settings
type KafkaConsumerConfig struct {
	Brokers            string `env:"KAFKA_BROKERS,required"`
	GroupID            string `env:"KAFKA_CONSUMER_GROUP_ID" envDefault:"default-group"`
	AutoOffsetReset    string `env:"KAFKA_CONSUMER_AUTO_OFFSET_RESET" envDefault:"earliest"`
	EnableAutoCommit   bool   `env:"KAFKA_CONSUMER_ENABLE_AUTO_COMMIT" envDefault:"false"`
	MaxPollIntervalMs  int    `env:"KAFKA_CONSUMER_MAX_POLL_INTERVAL_MS" envDefault:"300000"`
	SessionTimeoutMs   int    `env:"KAFKA_CONSUMER_SESSION_TIMEOUT_MS" envDefault:"45000"`
	HeartbeatIntervalMs int   `env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL_MS" envDefault:"3000"`
	RetryBackoffMs     int    `env:"KAFKA_CONSUMER_RETRY_BACKOFF_MS" envDefault:"100"`
	FetchMinBytes      int    `env:"KAFKA_CONSUMER_FETCH_MIN_BYTES" envDefault:"1"`
	FetchWaitMaxMs     int    `env:"KAFKA_CONSUMER_FETCH_WAIT_MAX_MS" envDefault:"500"`
}

var _ utils.Config = (*KafkaConsumerConfig)(nil)

func (c *KafkaConsumerConfig) Load() error {
	log.Printf("Loading KafkaConsumerConfig")
	return utils.ParseConfig(c)
}

// SchemaRegistryConfig holds Schema Registry settings
type SchemaRegistryConfig struct {
	URL string `env:"KAFKA_SCHEMA_REGISTRY_URL,required"`
}

var _ utils.Config = (*SchemaRegistryConfig)(nil)

func (c *SchemaRegistryConfig) Load() error {
	log.Printf("Loading SchemaRegistryConfig")
	return utils.ParseConfig(c)
}

// DefaultConfig holds the default Kafka settings
var DefaultConfig = struct {
	Producer KafkaProducerConfig
	Consumer KafkaConsumerConfig
	Schema   SchemaRegistryConfig
}{
	Producer: KafkaProducerConfig{
		Brokers:  "localhost:9092",
		ClientID: "default-client",
	},
	Consumer: KafkaConsumerConfig{
		Brokers:            "localhost:9092",
		GroupID:            "default-group",
		AutoOffsetReset:    "earliest",
		EnableAutoCommit:   false,
		MaxPollIntervalMs:  300000,
		SessionTimeoutMs:   45000,
		HeartbeatIntervalMs: 3000,
		RetryBackoffMs:     100,
		FetchMinBytes:      1,
		FetchWaitMaxMs:     500,
	},
	Schema: SchemaRegistryConfig{
		URL: "http://localhost:8081",
	},
}

// Option is a functional option for configuring Kafka
type Option func(*KafkaProducerConfig, *KafkaConsumerConfig, *SchemaRegistryConfig)

// WithBrokers sets Kafka brokers
func WithBrokers(brokers string) Option {
	return func(p *KafkaProducerConfig, c *KafkaConsumerConfig, _ *SchemaRegistryConfig) {
		p.Brokers = brokers
		c.Brokers = brokers
	}
}

// WithClientID sets Kafka client ID for producer
func WithClientID(clientID string) Option {
	return func(p *KafkaProducerConfig, _ *KafkaConsumerConfig, _ *SchemaRegistryConfig) {
		p.ClientID = clientID
	}
}

// WithConsumerGroupID sets Kafka consumer group ID
func WithConsumerGroupID(groupID string) Option {
	return func(_ *KafkaProducerConfig, c *KafkaConsumerConfig, _ *SchemaRegistryConfig) {
		c.GroupID = groupID
	}
}

// WithSchemaRegistryURL sets the schema registry URL
func WithSchemaRegistryURL(url string) Option {
	return func(_ *KafkaProducerConfig, _ *KafkaConsumerConfig, s *SchemaRegistryConfig) {
		s.URL = url
	}
}

// WithAutoOffsetReset sets the auto offset reset policy
func WithAutoOffsetReset(offset string) Option {
	return func(_ *KafkaProducerConfig, c *KafkaConsumerConfig, _ *SchemaRegistryConfig) {
		c.AutoOffsetReset = offset
	}
}
