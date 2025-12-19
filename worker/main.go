package main

import (
	"encoding/json"
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

func main() {
	// conn, err := amqp091.Dial("amqp://guest:guest@rabbitmqfix:5672")
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	channel.ExchangeDeclare(
		"fleet.events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	// queue, err := channel.QueueDeclare(
	// 	"geofence_alerts",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// channel.QueueBind(
	// 	queue.Name,
	// 	"geofence_alerts",
	// 	"fleet.events",
	// 	false,
	// 	nil,
	// )

	err = channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(
		"geofence_alerts",
		"consumer-geofence",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	var geo GeofenceEvent
	for message := range msgs {
		err := json.Unmarshal(message.Body, &geo)
		if err != nil {
			fmt.Println("ERROR decode:", err)
			continue
		}

		fmt.Println("Received message Vehicle ID: ", geo.VehicleID)
		fmt.Println(string(message.Body))

		message.Ack(false)
	}
}
