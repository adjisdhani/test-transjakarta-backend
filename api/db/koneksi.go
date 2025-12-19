package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	// DB, err = sql.Open(
	// 	"postgres",
	// 	os.Getenv("DATABASE_URL"),
	// )

	DB, err = sql.Open(
		"postgres",
		"postgres://postgres:adjis@localhost:5432/backend_transjakarta?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL connected")
}
