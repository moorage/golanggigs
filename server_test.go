package main_test

import (
	"testing"
	"./scrapers"
)

/* TODO: Implement! */
func TestScraping(t *testing.T) {
	result, err := scrapers.Google("golang", "theresumator.com")
	if (err) { t.Errorf("Couldn't call scrapers.Google on %s, %s: %s", "golang", "thresumator.com", err.Message()) }
	fmt.Printf(result)
}