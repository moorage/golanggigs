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
	log.Printf("JOB::::::::::::::: %#v", job)
}