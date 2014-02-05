//added some comments in jamesbranch
package main

import (
	"github.com/moorage/golanggigs/models"
	"github.com/moorage/golanggigs/scrapers"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/bmizerany/pq"
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

	results := []models.Job{}

	result, err := scrapers.Google("golang", "theresumator.com")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.ResumatorJob(result[i])
		if err != nil {
			panic(err)
		}
		results = append(results, job)
	}
	// fmt.Printf("The Resumator [%d results]: %#v\n", len(result), result)

	result, err = scrapers.Google("golang", "jobscore.com")
	if err != nil {
		panic(err)
	}
	//results = append(results, result...)
	fmt.Printf("Jobscore [%d results]: %#v\n", len(result), result)

	result, err = scrapers.Google("golang", "jobvite.com")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.JobviteJob(result[i])
		if err != nil {
			panic(err)
		}
		if len(job.JobTitle) > 0 {
			results = append(results, job)
		}
	}
	// fmt.Printf("Jobvite [%d results]: %#v\n", len(result), result)

	result, err = scrapers.Google("golang +\"/jobs2/\"", "linkedin.com")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.LinkedinJob(result[i])
		if err != nil {
			panic(err)
		}
		if len(job.JobTitle) > 0 {
			results = append(results, job)
		}
	}
	// fmt.Printf("LinkedIn [%d results]: %#v\n", len(result), result)

	result, err = scrapers.Github("golang")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.GithubJob(result[i])
		if err != nil {
			panic(err)
		}
		results = append(results, job)
	}
	// fmt.Printf("GitHub [%d results]: %#v\n", len(result), result)

	result, err = scrapers.StackOverflow("golang")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.StackOverflowJob(result[i])
		if err != nil {
			panic(err)
		}
		if len(job.JobTitle) > 0 {
			results = append(results, job)
		}
	}
	// fmt.Printf("StackOverflow [%d results]: %#v\n", len(result), result)

	result, err = scrapers.Dice("golang")
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(result); i++ {
		job, err := scrapers.DiceJob(result[i])
		if err != nil {
			panic(err)
		}
		if len(job.JobTitle) > 0 {
			results = append(results, job)
		}
	}
	// fmt.Printf("Dice [%d results]: %#v\n", len(result), result)

	jsonified, err := json.Marshal(results)
	if err != nil {
		fmt.Fprint(res, "{error:true,message:\""+err.Error()+"\"}")
	} else {
		fmt.Fprint(res, string(jsonified))
	}
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

	db, err := sql.Open("postgres", "user="+dbUser+" password="+dbPw+" dbname="+dbName+" sslmode=verify-full port="+dbPort+" host="+dbHost)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", Index)
	http.HandleFunc("/index.json", IndexJson)
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("public/stylesheets"))))
	http.Handle("/javascripts/", http.StripPrefix("/javascripts/", http.FileServer(http.Dir("public/javascripts"))))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	if err != nil {
		panic(err)
	}
}
