package dyffiBroker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type kafkaBroker struct {
	brokerConfig BrokerConfig
}

type KafkaParams struct {
	Key       *string
	Partition *int
	Offset    *int64
	Headers   []kafka.Header
	Timeout   *time.Time
	Params    map[string]interface{}
}

func (k KafkaParams) GetQueueType() QueueType {
	return Kafka
}

func (k *kafkaBroker) Publish(topic string, message []byte, params DefaultParams) error {
	if k == nil || k.brokerConfig.Host == "" || k.brokerConfig.Port == "" {
		return fmt.Errorf("Kafka broker is not properly initialized")
	}

	if params == nil {
		params = KafkaParams{}
	}

	if params.GetQueueType() != Kafka {
		return fmt.Errorf("invalid parameters for Kafka broker")
	}

	kafkaParams, ok := params.(KafkaParams)
	if !ok {
		return fmt.Errorf("invalid Kafka parameters provided")
	}

	dns := fmt.Sprintf("%s:%s", k.brokerConfig.Host, k.brokerConfig.Port)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{dns},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
		},
	})
	defer writer.Close()

	var key []byte
	if kafkaParams.Key != nil {
		key = []byte(*kafkaParams.Key)
	} else {
		key = nil
	}

	var partition int
	if kafkaParams.Partition != nil {
		partition = *kafkaParams.Partition
	} else {
		partition = 0
	}

	var offset int64
	if kafkaParams.Offset != nil {
		offset = *kafkaParams.Offset
	} else {
		offset = kafka.FirstOffset
	}

	var timestamp time.Time
	if kafkaParams.Timeout != nil {
		timestamp = *kafkaParams.Timeout
	} else {
		timestamp = time.Now()
	}

	headers := kafkaParams.Headers

	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:       key,
		Partition: partition,
		Offset:    offset,
		Headers:   headers,
		Time:      timestamp,
		Value:     message,
	})

	if err != nil {
		return fmt.Errorf("failed to write messages to Kafka: %w", err)
	}

	return nil
}
