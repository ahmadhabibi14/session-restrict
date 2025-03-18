package configs

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	EnvDev  string = `dev`
	EnvStag string = `stag`
	EnvProd string = `prod`
)

func LoadEnv() {
	err := godotenv.Load(`.env`)
	if err != nil {
		log.Fatal(err)
	}
}
