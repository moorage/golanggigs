package main_test

import (
	"testing"
	"./scrapers"
	"log"
)

/* TODO: Implement! */
func TestScraping(t *testing.T) {
	job, err := scrapers.ResumatorJob("http://iron.theresumator.com/apply/8KLLHQ/Go-Developer.html")
	if err != nil { t.Errorf("Couldn't call scrapers.ResumatorJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)
	

	job, err = scrapers.GithubJob("https://jobs.github.com/positions/ca666b62-4b3a-11e3-81ad-9804406c2307")
	if err != nil { t.Errorf("Couldn't call scrapers.GithubJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)
	

	job, err = scrapers.LinkedinJob("http://www.linkedin.com/jobs2/view/9696846")
	if err != nil { t.Errorf("Couldn't call scrapers.LinkedinJob %s", err.Error()) }
	log.Printf("JOB::::::::::::::: %#v", job)
}