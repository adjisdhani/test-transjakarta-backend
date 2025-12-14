package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	db "subscribe_backend/db"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rabbitmq/amqp091-go"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type GeofenceEvent struct {
	VehicleID string `json:"vehicle_id"`
	Event     string `json:"event"`
	Location  struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Timestamp int64 `json:"timestamp"`
}

func decodeStrict(body []byte) (VehicleLocation, error) {
	var loc VehicleLocation

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&loc); err != nil {
		return loc, err
	}

	return loc, nil
}

func validateLocation(loc VehicleLocation) error {
	if loc.VehicleID == "" {
		return errors.New("vehicle_id kosong")
	}

	if loc.Latitude == 0 || loc.Longitude == 0 {
		return errors.New("latitude / longitude kosong")
	}

	if loc.Latitude < -90 || loc.Latitude > 90 {
		return errors.New("latitude tidak valid")
	}

	if loc.Longitude < -180 || loc.Longitude > 180 {
		return errors.New("longitude tidak valid")
	}

	if loc.Timestamp <= 0 {
		return errors.New("timestamp tidak valid")
	}

	return nil
}

func saveLocationToDB(loc VehicleLocation) error {
	_, err := db.DB.Exec(`
		INSERT INTO vehicle_locations
		(vehicle_id, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4)
	`,
		loc.VehicleID,
		loc.Latitude,
		loc.Longitude,
		loc.Timestamp,
	)

	return err
}

func getFirstLocationFromDB(vehicleID string) (*VehicleLocation, error) {
	row := db.DB.QueryRow(`
		SELECT latitude, longitude, timestamp
		FROM vehicle_locations
		WHERE vehicle_id = $1
		ORDER BY timestamp ASC
		LIMIT 1
	`, vehicleID)

	var loc VehicleLocation
	err := row.Scan(&loc.Latitude, &loc.Longitude, &loc.Timestamp)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	loc.VehicleID = vehicleID
	return &loc, nil
}

func publishGeofenceEvent(channel *amqp091.Channel, loc VehicleLocation) {
	event := GeofenceEvent{
		VehicleID: loc.VehicleID,
		Event:     "geofence_entry",
		Location: struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}{
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		},
		Timestamp: loc.Timestamp,
	}

	body, _ := json.Marshal(event)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,
		"fleet.events",
		"geofence_alerts",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Passing geofence event:", "geofence_alerts", string(body))
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // meter

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func main() {
	db.InitDB()

	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmqfix:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	opts := mqtt.NewClientOptions().
		AddBroker("tcp://mosquitto:1883").
		SetClientID("location-subscriber")

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mqttClient.Subscribe(
		"fleet/vehicle/+/location",
		1,
		func(c mqtt.Client, m mqtt.Message) {

			loc, err := decodeStrict(m.Payload())
			if err != nil {
				return
			}

			if err := validateLocation(loc); err != nil {
				return
			}

			if err := saveLocationToDB(loc); err != nil {
				return
			}

			firstLoc, err := getFirstLocationFromDB(loc.VehicleID)
			if err != nil {
				return
			}

			fmt.Println("Data pertama:", firstLoc)

			if firstLoc != nil {
				distance := calculateDistance(
					firstLoc.Latitude,
					firstLoc.Longitude,
					loc.Latitude,
					loc.Longitude,
				)

				fmt.Println("Jaraknya:", distance)

				if distance >= 50 {
					publishGeofenceEvent(channel, loc)
				}
			}

			fmt.Println("Topic:", m.Topic())
			fmt.Println("Payload:", string(m.Payload()))
		},
	)

	select {}
}
