package main

import (
	"github.com/Ametion/dyffi"
	dyffiBroker "github.com/Ametion/dyffi/broker"
	"time"
)

func main() {
	// Define Kafka broker configuration
	brokerConf := dyffiBroker.BrokerConfig{
		BrokerType: dyffiBroker.Kafka,
		Host:       "localhost",
		Port:       "9092",
	}

	engine := dyffi.NewDyffiEngine()
	engine.UseBroker(brokerConf)
	engine.SetDevelopment()

	engine.Post("/publish", func(context *dyffi.Context) {

		kafkaParams := dyffiBroker.KafkaParams{
			Key:       dyffiBroker.String("key"),
			Partition: dyffiBroker.Int(0),
			Offset:    dyffiBroker.Int64(0),
			Headers:   nil,
			Timeout:   dyffiBroker.Time(time.Now()),
		}

		err := context.PublishToQueue("test_topic", []byte("Hello, Dyffi!"), kafkaParams)
		if err != nil {
			context.SendJSON(500, "Failed to publish message: "+err.Error())
			return
		}

		context.SendJSON(200, "Message Published Successfully!")
	})

	engine.Run(":8080")
}
