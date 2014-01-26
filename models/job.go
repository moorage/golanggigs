package models

import (
	"time"
)

type Job struct {
	Id              int64
	JobTitle        string
	JobLocation     string
	JobDescription  string
	HowToApply      string
	CompanyLocation string
	CompanyName     string
	CompanyUrl      string
	SourceUrl       string
	SourceName      string
	PostedAt        time.Time
}

