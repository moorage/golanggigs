package main

import (
	"log"
	"strings"
	"os"
	"database/sql"
	_ "github.com/bmizerany/pq"
)

func main() {
	postgresUrl := os.Getenv("HEROKU_POSTGRESQL_AQUA_URL")
	if len(postgresUrl) < 1 {
		panic("HEROKU_POSTGRESQL_AQUA_URL is not set")
	}

	dbUser := strings.Split(strings.Split(postgresUrl, "postgres://")[1], ":")[0]
	dbPw := strings.Split(strings.Split(strings.Split(postgresUrl, "postgres://")[1], ":")[1], "@")[0]
	dbHost := strings.Split(strings.Split(postgresUrl, "@")[1], ":")[0]
	dbPort := strings.Split(strings.Split(strings.Split(postgresUrl, "@")[1], ":")[1], "/")[0]
	dbName := strings.Split(strings.Split(postgresUrl, "@")[1], "/")[1]

	db, err := sql.Open("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=verify-full port="+dbPort+" host="+dbHost)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	log.Println("Creating Jobs Table.")
	_, err = db.Exec("create table jobs (" + 
	"Id               BIGSERIAL     PRIMARY KEY," + 
	"JobTitle         VARCHAR(512)  NOT NULL," + 
	"JobLocation      VARCHAR(512)          ," + 
	"JobDescription   TEXT                  ," + 
	"HowToApply       VARCHAR(512)          ," + 
	"CompanyLocation  VARCHAR(512)          ," + 
	"CompanyName      VARCHAR(512)  NOT NULL," + 
	"CompanyUrl       VARCHAR(1024)         ," + 
	"SourceUrl        VARCHAR(512)          ," + 
	"SourceName       VARCHAR(512)  NOT NULL," + 
	"PostedAt         DATETIME              ," + 
	"CreatedAt        DATETIME      NOT NULL," + 
	");")
	
	if err != nil {
		panic(err)
	}
	
	log.Println("Done with seeds.go.")
}