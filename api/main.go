package main

import (
	db "api/db"
	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	db.InitDB()

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	})

	app.Get("/vehicles/:vehicle_id/location", func(c *fiber.Ctx) error {
		vehicle_id := c.Params("vehicle_id")

		row := db.DB.QueryRow(`
			SELECT vehicle_id, latitude, longitude, timestamp
			FROM vehicle_locations
			WHERE vehicle_id = $1
			ORDER BY timestamp DESC
			LIMIT 1
		`, vehicle_id)

		var loc VehicleLocation
		err := row.Scan(
			&loc.VehicleID,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Timestamp,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).JSON(fiber.Map{
					"message": "vehicle not found",
				})
			}

			return c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(loc)
	})

	app.Get("/vehicles/:vehicle_id/history", func(c *fiber.Ctx) error {
		vehicle_id := c.Params("vehicle_id")
		startStr := c.Query("start")
		endStr := c.Query("end")

		if startStr == "" || endStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "format salah",
			})
		}

		start, err1 := strconv.ParseInt(startStr, 10, 64)
		end, err2 := strconv.ParseInt(endStr, 10, 64)

		if err1 != nil || err2 != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "format salah",
			})
		}

		if start <= 0 || end <= 0 || start > end {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "format salah",
			})
		}

		rows, err := db.DB.Query(`
			SELECT vehicle_id, latitude, longitude, timestamp
			FROM vehicle_locations
			WHERE vehicle_id = $1
			AND timestamp BETWEEN $2 AND $3
			ORDER BY timestamp ASC
		`, vehicle_id, start, end)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		defer rows.Close()

		locations := []VehicleLocation{}

		for rows.Next() {
			var loc VehicleLocation
			if err := rows.Scan(
				&loc.VehicleID,
				&loc.Latitude,
				&loc.Longitude,
				&loc.Timestamp,
			); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": err.Error()})
			}
			locations = append(locations, loc)
		}

		return c.JSON(locations)
	})

	err := app.Listen(":3000")

	if err != nil {
		panic(err)
	}
}
