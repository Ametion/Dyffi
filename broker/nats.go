package dyffiBroker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

type natsBroker struct {
	brokerConfig BrokerConfig
}

type NatsParams struct {
	Subject *string
	Timeout *time.Duration
	Headers nats.Header
}

func (n NatsParams) GetQueueType() QueueType {
	return Nats
}

func (n *natsBroker) Publish(topic string, message []byte, params DefaultParams) error {
	if params.GetQueueType() != Nats {
		return fmt.Errorf("invalid parameters for NATS broker")
	}

	natsParams, ok := params.(NatsParams)
	if !ok {
		return fmt.Errorf("invalid NATS parameters provided")
	}

	dns := fmt.Sprintf("nats://%s:%s", n.brokerConfig.Host, n.brokerConfig.Port)

	nc, err := nats.Connect(dns)
	if err != nil {
		return err
	}
	defer nc.Close()

	var subject string
	if natsParams.Subject != nil {
		subject = *natsParams.Subject
	} else {
		subject = topic
	}

	if natsParams.Headers != nil {
		err = nc.PublishMsg(&nats.Msg{
			Subject: subject,
			Data:    message,
			Header:  natsParams.Headers,
		})
	} else {
		err = nc.Publish(subject, message)
	}

	if err != nil {
		return err
	}

	return nil
}
