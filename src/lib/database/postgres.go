package database

import (
	"fmt"
	"os"
	"session-restrict/src/lib/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ConnPg *sqlx.DB

func ConnectPostgresSQL() {
	postgresDbName := os.Getenv("POSTGRES_DB")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	postgresURL := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		postgresUser, postgresPassword, postgresDbName, postgresHost, postgresPort,
	)

	ConnPg = sqlx.MustConnect("postgres", postgresURL)

	err := ConnPg.Ping()
	if err != nil {
		ConnPg.Close()
		logger.Log.Fatal(err, `failed to test connection`)
	}
}
