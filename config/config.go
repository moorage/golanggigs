package config

import (
	"os"
	"errors"
	"strings"
)

const POSTGRES_ENV_VAR string = "HEROKU_POSTGRESQL_AQUA_URL"

func PostgresDataSource() (string, error) {
	postgresUrl := os.Getenv(POSTGRES_ENV_VAR)
	if len(postgresUrl) < 1 {
		return "", errors.New(POSTGRES_ENV_VAR + " is not set")
	}

	dbUser := strings.Split(strings.Split(postgresUrl, "postgres://")[1], ":")[0]
	dbPw := strings.Split(strings.Split(strings.Split(postgresUrl, "postgres://")[1], ":")[1], "@")[0]
	dbHost := strings.Split(strings.Split(postgresUrl, "@")[1], ":")[0]
	dbPort := strings.Split(strings.Split(strings.Split(postgresUrl, "@")[1], ":")[1], "/")[0]
	dbName := strings.Split(strings.Split(postgresUrl, "@")[1], "/")[1]
	
	return "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=require port="+dbPort+" host="+dbHost, nil
}
