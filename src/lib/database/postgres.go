package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ConnPg *sqlx.DB

func ConnectPostgresSQL() error {
	postgresDbName := os.Getenv("POSTGRES_DB")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	postgresURL := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		postgresUser, postgresPassword, postgresDbName, postgresHost, postgresPort,
	)

	var err error
	ConnPg, err = sqlx.Connect("postgres", postgresURL)

	return err
}
