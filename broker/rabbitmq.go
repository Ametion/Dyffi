package dyffiBroker

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type rabbitMQBroker struct {
	brokerConfig BrokerConfig
}

type RabbitMQParams struct {
	Exchange    *string
	RoutingKey  *string
	Mandatory   bool
	Immediate   bool
	Expiration  *string
	ContentType *string
	Timeout     *time.Time
	Params      map[string]interface{}
}

func (r RabbitMQParams) GetQueueType() QueueType {
	return RabbitMQ
}

func (r RabbitMQParams) GetParams() map[string]interface{} {
	return r.Params
}

func (r *rabbitMQBroker) Publish(topic string, message []byte, params DefaultParams) error {
	dns := fmt.Sprintf("amqp://%s:%s@%s:%s/", r.brokerConfig.Username, r.brokerConfig.Password, r.brokerConfig.Host, r.brokerConfig.Port)

	if params.GetQueueType() != RabbitMQ {
		return fmt.Errorf("invalid parameters for RabbitMQ broker")
	}

	rabbitParams := params.(RabbitMQParams)

	conn, err := amqp.Dial(dns)
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to open a channel: %v", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to declare a BrokerType: %v", err)
		return err
	}

	exchange := ""
	if rabbitParams.Exchange != nil {
		exchange = *rabbitParams.Exchange
	}

	routingKey := q.Name
	if rabbitParams.RoutingKey != nil {
		routingKey = *rabbitParams.RoutingKey
	}

	contentType := "text/plain"
	if rabbitParams.ContentType != nil {
		contentType = *rabbitParams.ContentType
	}

	expiration := ""
	if rabbitParams.Expiration != nil {
		expiration = *rabbitParams.Expiration
	}

	err = ch.Publish(
		exchange,
		routingKey,
		rabbitParams.Mandatory,
		rabbitParams.Immediate,
		amqp.Publishing{
			ContentType: contentType,
			Body:        message,
			Expiration:  expiration,
		},
	)
	if err != nil {
		fmt.Printf("Failed to publish message: %v", err)
		return err
	}

	return nil
}
