package main

import (
	"log"
	"strings"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/moorage/golanggigs/scrapers"
	// "github.com/moorage/golanggigs/models"
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

	db, err := sql.Open("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=require port="+dbPort+" host="+dbHost)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	// db.Exec("drop table jobs")
	
	log.Println("Checking if Jobs Table Exists.")
	rows, err := db.Query("select exists(select * from information_schema.tables where table_name = $1)", "jobs")
	
	tableExists := false
	
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&tableExists)
		if err != nil {
			panic(err)
		}
	}
	
	if !tableExists {
		log.Println("Jobs table missing. Creating Jobs Table.")
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
		"PostedAt         timestamp             ," + 
		"CreatedAt        timestamp DEFAULT current_timestamp NOT NULL" + 
		");")
	
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("Jobs table exists. Skipping creation.")
	}

	// results := []models.Job{}
	// 
	// result, err := scrapers.Google("golang +\"/jobs2/\"", "linkedin.com")
	// if err != nil { panic(err) }
	// 
	// for i := 0; i < len(result); i++ {
	// 	job, err := scrapers.LinkedinJob(result[i])
	// 	if err != nil { panic(err) }
	// 	
	// 	if len(job.JobTitle) > 0 {
	// 		results = append(results, job)
	// 	}
	// }
	// // fmt.Printf("LinkedIn [%d results]: %#v\n", len(result), result)
	// 
	// if len(results) > 0 {
	// 	db, err := sql.Open("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=require port="+dbPort+" host="+dbHost)
	// 	if err != nil { panic(err) }
	// 
	// 	defer db.Close()
	// 	
	// 	err = results[0].Create(db)
	// 	if err != nil { panic(err) }
	// }

	log.Println("Scraping jobs.")
	err = scrapers.ScrapeJobs("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=require port="+dbPort+" host="+dbHost)
	if err != nil {
		panic(err)
	}
	
	log.Println("Done with seeds.go.")
}