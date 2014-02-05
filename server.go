//added some comments in jamesbranch
package main

import (
	"github.com/moorage/golanggigs/scrapers"
	"database/sql"
	"log"
	"time"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strings"
	// "github.com/hoisie/mustache"
	// "github.com/garyburd/redigo/redis"
)

func Index(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "application.html")
}

func IndexJson(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// jsonified, err := json.Marshal(results)
	// if err != nil {
	// 	fmt.Fprint(res, "{error:true,message:\""+err.Error()+"\"}")
	// } else {
	// 	fmt.Fprint(res, string(jsonified))
	// }
}

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

	// Find and save new jobs every hour.
	go func() {
		for {
			err = scrapers.ScrapeJobs("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=require port="+dbPort+" host="+dbHost)
			if err != nil {
				log.Printf("ERROR occurred when calling main models.ScrapeJobs: %#v\n", err.Error())
			}
			time.Sleep(60 * time.Minute)
		}
	}()

	http.HandleFunc("/", Index)
	http.HandleFunc("/index.json", IndexJson)
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("public/stylesheets"))))
	http.Handle("/javascripts/", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("public/javascripts"))))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		panic(err)
	}
}
