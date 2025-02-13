package kafka

import (
	"context"
	"embed"
	"fmt"
	"log"
	"proposal-template/pkg/utils"
	config "proposal-template/pkg/utils/config"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
)

// region:      ======= producer implement =======
type kafkaPublisher struct {
	producer *kafka.Producer
	serde    serde.Serializer
	topic    string
}


var _ Publisher = (*kafkaPublisher)(nil)

func NewKafkaPublisher(producer *kafka.Producer, sr *SchemaRegistry, schemaID int, topic string) (*kafkaPublisher, error) {
	serde, err := avrov2.NewSerializer(sr.client, serde.ValueSerde, &avrov2.SerializerConfig{
		SerializerConfig: serde.SerializerConfig{
			AutoRegisterSchemas: false,
			UseSchemaID:         schemaID,
			UseLatestVersion:    true,
			NormalizeSchemas:    true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create avro serializer: %s", err)
	}
	return &kafkaPublisher{producer: producer, serde: serde, topic: topic}, nil
}

func (s *kafkaPublisher) SendMessage(ctx context.Context, value interface{}) error {
	deliveryChan := make(chan kafka.Event)

	payload, err := s.serde.Serialize(s.topic, &value)
	if err != nil {
		return fmt.Errorf("failed to serialize: %s", err)
	}

	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &s.topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("produce failed: %v", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %v", m.TopicPartition.Error)
	}

	fmt.Printf("Delivered message to topic: %s [%d] at offset: %v\n",
		*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)

	return nil
}

// endregion:      ======= producer implement =======

// region:      ======= consumer implement =======

type kafkaSubscriber struct {
	consumer *kafka.Consumer
	serde    serde.Deserializer
	topic    string
}

var _ Subscriber = (*kafkaSubscriber)(nil)

func NewKafkaSubscriber(consumer *kafka.Consumer, sr *SchemaRegistry, topic string) (*kafkaSubscriber, error) {
	serde, err := avrov2.NewDeserializer(sr.client, serde.ValueSerde, &avrov2.DeserializerConfig{
		DeserializerConfig: serde.DeserializerConfig{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create avro serializer: %s", err)
	}
	return &kafkaSubscriber{consumer: consumer, serde: serde, topic: topic}, nil
}

func (s *kafkaSubscriber) SubscribeToTopic(ctx context.Context) error {
	err := s.consumer.SubscribeTopics([]string{s.topic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %s", err)
	}
	return nil
}

func (s *kafkaSubscriber) ConsumeMessages(
	ctx context.Context,

	msgTypeConstructor func() ConsumerMessage,
) (<-chan ConsumerMessage, <-chan error, chan<- bool) {
	chMsg := make(chan ConsumerMessage)
	chCommitRequest := make(chan bool)
	chErr := make(chan error)
	go func() {
		defer close(chMsg)
		defer close(chErr)
		defer close(chCommitRequest)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := s.consumer.ReadMessage(100 * time.Millisecond)
				if err != nil {
					if err.(kafka.Error).Code() == kafka.ErrTimedOut {
						continue // Normal timeout, just retry
					}

					// Log and potentially retry or send to error channel
					chErr <- fmt.Errorf("consumer read error: %v", err)
					continue
				}

				// Deserialize and process message
				msgObj := msgTypeConstructor()
				err = s.serde.DeserializeInto(s.topic, msg.Value, &msgObj)
				if err != nil {
					chErr <- fmt.Errorf("deserialization error: %v", err)
					continue
				}
				log.Printf("Message on Topic: %s, Offset: %+v\n", *msg.TopicPartition.Topic, msg.TopicPartition.Offset)
				chMsg <- msgObj

				// Manual offset commit
				if <-chCommitRequest {
					_, err := s.consumer.CommitMessage(msg)
					if err != nil {
						chErr <- fmt.Errorf("offset commit error: %v", err)
					}
				}
			}
		}
	}()
	return chMsg, chErr, chCommitRequest
}

// endregion:      ======= consumer implement =======

// region:          ======= schema registry =======

const (
	retryCount    = 10
	retryInterval = 1 * time.Second
)


type SchemaRegistry struct {
	client schemaregistry.Client
}

func NewSchemaRegistry(opts ...Option) (*SchemaRegistry, error) {
	
	producerConfig := DefaultConfig.Producer
	consumerConfig := DefaultConfig.Consumer
	schemaConfig := DefaultConfig.Schema

	// Apply functional options
	for _, opt := range opts {
		opt(&producerConfig, &consumerConfig, &schemaConfig)
	}
	sr, err := utils.Retry(retryCount, retryInterval, func() (schemaregistry.Client, error) {
		return schemaregistry.NewClient(schemaregistry.NewConfig(schemaConfig.URL))
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create schema registry client: %s", err)
	}
	return &SchemaRegistry{
		client: sr,
	}, nil
}

func (s *SchemaRegistry) Close() {
	s.client.Close()
}

func (s *SchemaRegistry) FindOrCreateArvoSchema(topic string, baseFS embed.FS, fullFileName string) (int, error) {
	schema, err := baseFS.ReadFile(fullFileName)
	if err != nil {
		return 0, fmt.Errorf("failed to read avro schema file: %s", err)
	}
	id, err := s.createAvroSchema(topic, schema)
	if err != nil {
		return 0, fmt.Errorf("failed to create avro schema: %s", err)
	}
	return id, nil

}

func (s *SchemaRegistry) createAvroSchema(topic string, schema []byte) (int, error) {
	id, err := utils.Retry(retryCount, retryInterval, func() (int, error) {
		return s.client.Register(topic+"-value", schemaregistry.SchemaInfo{
			Schema: string(schema),
		}, false)
	})
	if err != nil {
		return 0, fmt.Errorf("failed to register schema: %s", err)
	}
	return id, nil
}

func (s *SchemaRegistry) CreateAvroSerializer(schemaConfig avrov2.SerializerConfig) (*avrov2.Serializer, error) {
	serde, err := utils.Retry(retryCount, retryInterval, func() (*avrov2.Serializer, error) {
		return avrov2.NewSerializer(s.client, serde.ValueSerde, &schemaConfig)
	})
	if err != nil {
		log.Printf("Failed to create avro serializer: %s", err)
		return nil, err
	}
	return serde, nil
}

func (s *SchemaRegistry) CreateAvroDeserializer(schemaConfig avrov2.DeserializerConfig) (*avrov2.Deserializer, error) {
	serde, err := utils.Retry(retryCount, retryInterval, func() (*avrov2.Deserializer, error) {
		return avrov2.NewDeserializer(s.client, serde.ValueSerde, &schemaConfig)
	})
	if err != nil {
		log.Printf("Failed to create avro deserializer: %s", err)
		return nil, err
	}
	return serde, nil
}

func CreateTopicIfNotExist(adminClient *kafka.AdminClient, topicName string, numPartitions int, replicationFactor int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		panic("ParseDuration(3s)")
	}
	// Create the topic if it does not exist
	results, err := utils.Retry(retryCount, retryInterval, func() ([]kafka.TopicResult, error) {
		return adminClient.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{{
				Topic:             topicName,
				NumPartitions:     numPartitions,
				ReplicationFactor: replicationFactor,
			}},
			kafka.SetAdminOperationTimeout(maxDur))
	})
	if err != nil {
		fmt.Printf("Failed to create topic: %v\n", err)
		os.Exit(1)
	}
	// Print results
	for _, result := range results {
		fmt.Printf("Topic created: %v\n", result)
	}
	return nil
}

// endregion:       ======= schema registry =======
// region:          ======= kafka utils =======
// type KafkaProducerConfig struct {
// 	// Kafka broker servers. Example :  "localhost:9092,localhost:9093"
// 	Brokers string `env:"KAFKA_BROKERS,required"`
// 	// An id string to pass to the server when making requests. The purpose of this is to be able to track the source of requests beyond just ip/port by allowing a logical application name to be included in server-side request logging.
// 	ClientID string `env:"KAFKA_CLIENT_ID" envDefault:""`
// }

// KAFKA_BROKERS=localhost:9092
// KAFKA_TOPIC=default
// KAFKA_CONSUMER_GROUP_ID=mygroup
// KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081

var _ config.IConfig = (*KafkaProducerConfig)(nil)


func NewKafkaProducer(opts ...Option) (*kafka.Producer, error) {
	producerConfig := DefaultConfig.Producer
	consumerConfig := DefaultConfig.Consumer
	schemaConfig := DefaultConfig.Schema

	// Apply functional options
	for _, opt := range opts {
		opt(&producerConfig, &consumerConfig, &schemaConfig)
	}
	return kafka.NewProducer(&kafka.ConfigMap{
		// "bootstrap.servers": cfg.Brokers,
		// "client.id":         cfg.ClientID,
		"bootstrap.servers": producerConfig.Brokers,
		"client.id":         producerConfig.ClientID,

	})
	// return p, nil
}

func NewAdminClientFromProducer(producer *kafka.Producer) (*kafka.AdminClient, error) {
	return kafka.NewAdminClientFromProducer(producer)
}


var _ config.IConfig = (*KafkaConsumerConfig)(nil)

/*
USAGE EXAMPLE:
- DEFAULT DECLARATION

consumer, err := kafka.NewKafkaConsumer() 
if err != nil {
	log.Fatalf("Failed to create Kafka producer: %v", err)
}

- CUSTOM DECLARATION:

consumer, err := service.NewKafkaConsumer(
	kafka.WithBrokers("custom-addr-of-broker"), //default addr == 'localhost:9092'
	kafka.WithClientID("custom-client-id"),
)

*/
func NewKafkaConsumer(opts ...Option) (*kafka.Consumer, error) {
	producerConfig := DefaultConfig.Producer
	consumerConfig := DefaultConfig.Consumer
	schemaConfig := DefaultConfig.Schema

	// Apply functional options
	for _, opt := range opts {
		opt(&producerConfig, &consumerConfig, &schemaConfig)
	}

	c, err := utils.Retry(10, 1*time.Second, func() (*kafka.Consumer, error) {
		return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     consumerConfig.Brokers,
		"group.id":              consumerConfig.GroupID,
		"auto.offset.reset":     consumerConfig.AutoOffsetReset,
		"enable.auto.commit":    consumerConfig.EnableAutoCommit,
		"max.poll.interval.ms":  consumerConfig.MaxPollIntervalMs,
		"session.timeout.ms":    consumerConfig.SessionTimeoutMs,
		"heartbeat.interval.ms": consumerConfig.HeartbeatIntervalMs,
		"retry.backoff.ms":      consumerConfig.RetryBackoffMs,
		"fetch.min.bytes":       consumerConfig.FetchMinBytes,
		"fetch.wait.max.ms":     consumerConfig.FetchWaitMaxMs,
	})
	})
	if err != nil {
		log.Printf("Failed to create kafka consumer: %s", err)
		return nil, err
	}
	return c, nil
}

func NewAdminClientFromConsumer(consumer *kafka.Consumer) (*kafka.AdminClient, error) {
	return kafka.NewAdminClientFromConsumer(consumer)
}

// endregion:       ======= kafka utils =======
