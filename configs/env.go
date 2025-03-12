package configs

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	dirRetryList := []string{``, `../`, `../../`, `../../../`}
	var err error
	for _, dirPrefix := range dirRetryList {
		envFile := dirPrefix + `.env`

		err = godotenv.Overload(envFile + `.dev`)
		if err == nil {
			slog.Info(`file .env.dev loaded (development environment)`)
			return
		}

		err = godotenv.Overload(envFile)
		if err == nil {
			slog.Info(`file .env loaded (staging environment)`)
			return
		}
	}
	panic("cannot load .env file")
}
