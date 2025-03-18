package tests

import (
	"fmt"
	"net/url"
	"os"
	"session-restrict/src/lib/database"
	"session-restrict/src/lib/logger"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	logger.InitLogger()

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		logger.Log.Fatal(err, "could not construct pool")
	}

	err = dockerPool.Client.Ping()
	if err != nil {
		logger.Log.Fatal(err, "could not connect to Docker")
	}

	setupPostgreSQL(dockerPool)
	setupRedis(dockerPool)

	// run tests
	m.Run()
}

const (
	testPostgresDb   = `testDb`
	testPostgresUser = `testUser`
	testPostgresPass = `testPass123`
)

func setupPostgreSQL(dockerPool *dockertest.Pool) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.6-alpine3.20",
		Env: []string{
			"POSTGRES_PASSWORD=" + testPostgresPass,
			"POSTGRES_USER=" + testPostgresUser,
			"POSTGRES_DB=" + testPostgresDb,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logger.Log.Fatal(err, "could not start resource")
	}

	dockerPool.MaxWait = 120 * time.Second

	os.Setenv("POSTGRES_DB", testPostgresDb)
	os.Setenv("POSTGRES_USER", testPostgresUser)
	os.Setenv("POSTGRES_PASSWORD", testPostgresPass)
	os.Setenv("POSTGRES_HOST", "127.0.0.1")
	os.Setenv("POSTGRES_PORT", resource.GetPort("5432/tcp"))

	if err = dockerPool.Retry(func() error {
		return database.ConnectPostgresSQL()
	}); err != nil {
		logger.Log.Fatal(err, "could not connect to docker")
	}

	logger.Log.Info("Connected to PostgreSQL")

	postgresUrl, err := url.Parse(fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		testPostgresUser, testPostgresPass, resource.GetHostPort("5432/tcp"), testPostgresDb,
	))
	if err != nil {
		logger.Log.Fatal(err, "failed to parse postgres url")
	}

	db := dbmate.New(postgresUrl)
	db.SchemaFile = "../db/schema.sql"
	db.MigrationsDir = []string{"../db/migrations"}

	migrations, err := db.FindMigrations()
	if err != nil {
		panic(err)
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}

	fmt.Println("\nApplying...")
	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}

func setupRedis(dockerPool *dockertest.Pool) {
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "8.0-M02-alpine",
		Env: []string{
			"ALLOW_EMPTY_PASSWORD=yes",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logger.Log.Fatal(err, "could not start resource")
	}

	dockerPool.MaxWait = 120 * time.Second

	hostAndPort := resource.GetHostPort("6379/tcp")

	os.Setenv("REDIS_ADDR", hostAndPort)

	database.ConnectRedis()

	if err = dockerPool.Retry(func() error {
		return database.ConnRd.Ping().Err()
	}); err != nil {
		logger.Log.Fatal(err, "could not connect to docker")
	}

	logger.Log.Info("Connected to Redis")
}
