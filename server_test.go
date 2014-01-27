package main_test

import (
	"testing"
	"./scrapers"
	"log"
)

/* TODO: Implement! */
func TestScraping(t *testing.T) {
	log.Printf("Running Resumator")
	job, err := scrapers.ResumatorJob("http://iron.theresumator.com/apply/8KLLHQ/Go-Developer.html")
	if err != nil { t.Errorf("Couldn't call scrapers.ResumatorJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)

	log.Printf("Running Github")
	job, err = scrapers.GithubJob("https://jobs.github.com/positions/ca666b62-4b3a-11e3-81ad-9804406c2307")
	if err != nil { t.Errorf("Couldn't call scrapers.GithubJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)

	log.Printf("Running linkedin")
	job, err = scrapers.LinkedinJob("http://www.linkedin.com/jobs2/view/9696846")
	if err != nil { 
		log.Printf("Got err in linkedin")
		log.Printf("Couldn't call scrapers.LinkedinJob %#v", err)
		t.Errorf("Couldn't call scrapers.LinkedinJob %s", err.Error()) 
	}
	log.Printf("Finished Linkedin")
	log.Printf("JOB::::::::::::::: %#v", job)

	log.Printf("Running stack overflow")
	job, err = scrapers.StackOverflowJob("http://careers.stackoverflow.com/jobs/11130/sr-web-developer-video-and-mobile-beachfront-media?a=dzBGurSM&searchTerm=golang")
	if err != nil { t.Errorf("Couldn't call scrapers.StackOverflowJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)

	log.Printf("Running Dive")
	job, err = scrapers.DiceJob("http://www.dice.com/job/result/cybercod/ST-JRSYDEV-AUS-12-106?src=19&q=golang")
	if err != nil { t.Errorf("Couldn't call scrapers.DiceJob %s", err.Error()) }
	// log.Printf("JOB::::::::::::::: %#v", job)

	log.Printf("Running Linkedin Google Search")
	_, err = scrapers.Google("golang +\"/jobs2/\"", "linkedin.com")
	if err != nil {
		panic(err)
	}
	// log.Printf("LinkedIn [%d results]: %#v\n", len(result), result)
}