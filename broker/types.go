package dyffiBroker

import (
	"fmt"
	"time"
)

type QueueType string

const (
	Kafka    QueueType = "kafka"
	RabbitMQ QueueType = "rabbitmq"
	Nats     QueueType = "nats"
)

type BrokerConfig struct {
	BrokerType QueueType `json:"queue_type"`
	Host       string    `json:"host"`
	Port       string    `json:"port"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
}

type Queue interface {
	Publish(topic string, message []byte, params DefaultParams) error
}

type DefaultParams interface {
	GetQueueType() QueueType
}

func String(s string) *string {
	return &s
}

func Int(i int) *int {
	return &i
}

func Int64(i int64) *int64 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

func Time(t time.Time) *time.Time {
	return &t
}

func Duration(d time.Duration) *time.Duration {
	return &d
}

func CreateNewBroker(conf BrokerConfig) (Queue, error) {
	switch conf.BrokerType {
	case Kafka:
		return &kafkaBroker{brokerConfig: conf}, nil
	case RabbitMQ:
		return &rabbitMQBroker{brokerConfig: conf}, nil
	case Nats:
		return &natsBroker{brokerConfig: conf}, nil
	default:
		return nil, fmt.Errorf("unsupported broker type: %s", conf.BrokerType)
	}
}
