package main

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	// 1️⃣ MQTT client options
	opts := mqtt.NewClientOptions().
		// AddBroker("tcp://mosquitto:1883").
		AddBroker("tcp://localhost:1883").
		SetClientID("vehicle-location-publisher").
		SetConnectTimeout(5 * time.Second)

	client := mqtt.NewClient(opts)

	// 2️⃣ Connect ke Mosquitto
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	// 3️⃣ Data lokasi (SAMA PERSIS)
	locations := []VehicleLocation{
		{"B1234XYZ", -6.254004764096112, 106.87614398105238, time.Now().Unix()},
		{"B1234XYZ", -6.253914764096112, 106.87614398105238, time.Now().Add(5 * time.Second).Unix()},  //10
		{"B1234XYZ", -6.253824764096112, 106.87614398105238, time.Now().Add(10 * time.Second).Unix()}, //20
		{"B1234XYZ", -6.253644764096112, 106.87614398105238, time.Now().Add(15 * time.Second).Unix()}, //40
		{"B1234XYZ", -6.253554764096112, 106.87614398105238, time.Now().Add(20 * time.Second).Unix()}, //50
		{"B1234XYZ", -6.253464764096112, 106.87614398105238, time.Now().Add(25 * time.Second).Unix()}, //60
		{"B1234XYZ", -6.253644764096112, 106.87614398105238, time.Now().Add(30 * time.Second).Unix()}, //40
		{"B1234XYZ", -6.253554764096112, 106.87614398105238, time.Now().Add(35 * time.Second).Unix()}, //50
		{"", -6.254004764096112, 106.87614398105238, 0},
		{"B1234XYZ", 0, 0, 0},
		{"B1234XYZ", -6.254004764096112, 106.87614398105238, 0},
		{"B1234XYZS", -6.254004764096112, 106.87614398105238, time.Now().Unix()},
		{"B1234XYZS", -6.253914764096112, 106.87614398105238, time.Now().Add(5 * time.Second).Unix()},  //10
		{"B1234XYZS", -6.253824764096112, 106.87614398105238, time.Now().Add(10 * time.Second).Unix()}, //20
		{"B1234XYZS", -6.253644764096112, 106.87614398105238, time.Now().Add(15 * time.Second).Unix()}, //40
		{"B1234XYZS", -6.253554764096112, 106.87614398105238, time.Now().Add(20 * time.Second).Unix()}, //50
		{"B1234XYZS", -6.253464764096112, 106.87614398105238, time.Now().Add(25 * time.Second).Unix()}, //60
		{"B1234XYZS", -6.253644764096112, 106.87614398105238, time.Now().Add(30 * time.Second).Unix()}, //40
		{"B1234XYZS", -6.253554764096112, 106.87614398105238, time.Now().Add(35 * time.Second).Unix()}, //50
		{"", -6.254004764096112, 106.87614398105238, 0},
		{"B1234XYZS", 0, 0, 0},
		{"B1234XYZS", -6.254004764096112, 106.87614398105238, 0},
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// 4️⃣ Publish loop
	for _, loc := range locations {
		<-ticker.C

		body, err := json.Marshal(loc)
		if err != nil {
			panic(err)
		}

		// 🔁 routingKey → topic
		topic := fmt.Sprintf(
			"fleet/vehicle/%s/location",
			loc.VehicleID,
		)

		token := client.Publish(
			topic,
			1,     // QoS (mirip at-least-once)
			false, // retained
			body,
		)

		// optional: tunggu publish selesai
		token.Wait()

		fmt.Println("Published:", topic, string(body))
	}
}
