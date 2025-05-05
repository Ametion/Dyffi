package main

import (
	"github.com/Ametion/dyffi"
	dyffiBroker "github.com/Ametion/dyffi/broker"
	"time"
)

func main() {
	// Define RabbitMQ broker configuration
	brokerConf := dyffiBroker.BrokerConfig{
		BrokerType: dyffiBroker.RabbitMQ,
		Host:       "localhost",
		Port:       "5672",
		Username:   "guest", // Default RabbitMQ credentials
		Password:   "guest", // Default RabbitMQ credentials
	}

	engine := dyffi.NewDyffiEngine()
	engine.UseBroker(brokerConf)
	engine.SetDevelopment()

	engine.Post("/publish", func(context *dyffi.Context) {

		rabbitMQParams := dyffiBroker.RabbitMQParams{
			Exchange:    dyffiBroker.String(""),
			RoutingKey:  dyffiBroker.String("test_queue"),
			Mandatory:   dyffiBroker.Bool(false),
			Immediate:   dyffiBroker.Bool(false),
			ContentType: dyffiBroker.String("text/plain"),
			Timeout:     dyffiBroker.Time(time.Now()),
		}

		err := context.PublishToQueue("test_queue", []byte("Hello, Dyffi with RabbitMQ!"), rabbitMQParams)
		if err != nil {
			context.SendJSON(500, "Failed to publish message: "+err.Error())
			return
		}

		context.SendJSON(200, "Message Published Successfully with RabbitMQ!")
	})

	engine.Run(":8080")
}
