//added some comments in jamesbranch
package main

import (
	"log"
	"time"
	"os"
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/hoisie/mustache"
	// "github.com/garyburd/redigo/redis"
	"github.com/moorage/golanggigs/scrapers"
	"github.com/moorage/golanggigs/config"
	"github.com/moorage/golanggigs/models"
)

func Index(res http.ResponseWriter, req *http.Request) {
	postgresDataSource, err := config.PostgresDataSource()
	if err != nil { panic(err) }
	
	db, err := sql.Open("postgres", postgresDataSource)
	if err != nil { panic(err) }
	defer db.Close()
	
	jobs, err := models.AllRecentJobs(db)
	if err != nil { panic(err) }

	fmt.Fprint(res, mustache.RenderFile("application.html.mustache", map[string][]*models.Job{"jobs": jobs}))
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
	postgresDataSource, err := config.PostgresDataSource()
	if err != nil { panic(err) }

	db, err := sql.Open("postgres", postgresDataSource)
	if err != nil { panic(err) }
	defer db.Close()

	// Find and save new jobs every hour.
	go func() {
		for {
			err = scrapers.ScrapeJobs("postgres", postgresDataSource)
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
