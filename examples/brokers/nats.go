package main

import (
	"github.com/Ametion/dyffi"
	dyffiBroker "github.com/Ametion/dyffi/broker"
	"time"
)

func main() {
	// Define NATS broker configuration
	brokerConf := dyffiBroker.BrokerConfig{
		BrokerType: dyffiBroker.Nats,
		Host:       "localhost",
		Port:       "4222",
	}

	engine := dyffi.NewDyffiEngine()
	engine.UseBroker(brokerConf)
	engine.IsDevelopment()

	engine.Post("/publish", func(context *dyffi.Context) {

		natsParams := dyffiBroker.NatsParams{
			Subject: dyffiBroker.String("test_subject"),
			Timeout: dyffiBroker.Duration(0 * time.Second),
			Headers: nil,
		}

		err := context.PublishToQueue("test_subject", []byte("Hello, Dyffi with NATS!"), natsParams)
		if err != nil {
			context.SendJSON(500, "Failed to publish message: "+err.Error())
			return
		}
		context.SendJSON(200, "Message Published Successfully with NATS!")
	})

	engine.Run(":8080")
}
