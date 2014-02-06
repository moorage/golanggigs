package models

import (
	"time"
	"errors"
	"database/sql"
	_ "github.com/lib/pq"
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
	CreatedAt       time.Time
}


func AllRecentJobs(db *sql.DB) ([]*Job, error) {
	jobs := []*Job{}
	
	rows, err := db.Query("SELECT Id, JobTitle, JobLocation, JobDescription, HowToApply, CompanyLocation, CompanyName, CompanyUrl, SourceUrl, SourceName, PostedAt, CreatedAt FROM jobs ORDER BY CreatedAt DESC LIMIT $1", 1000)
	if err != nil { return jobs, err }
	
	defer rows.Close()
	for rows.Next() {
		job := &Job{}
		jobs = append(jobs, job)
		
		err = job.scanRow(rows)
		if err != nil { return jobs, err }
	}

	err = rows.Err()
	if err != nil { return jobs, err }
	
	return jobs, nil
}

func (self *Job) scanRow(rows *sql.Rows) error {
	return rows.Scan(
		&self.Id,
		&self.JobTitle,
		&self.JobLocation,
		&self.JobDescription,
		&self.HowToApply,
		&self.CompanyLocation,
		&self.CompanyName,
		&self.CompanyUrl,
		&self.SourceUrl,
		&self.SourceName,
		&self.PostedAt,
		&self.CreatedAt,
	)
}

func (self *Job) FormattedJobDate() string {
	displayedDate := self.PostedAt
	if self.PostedAt.IsZero() {
		displayedDate = self.CreatedAt
	}
	
	layout := "Jan 2"
	if time.Now().Year() != displayedDate.Year() {
		layout = "Jan 2, 2006"
	}
	
	return displayedDate.Format(layout)
}

func (self *Job) PassesCreationValidation(db *sql.DB) (bool, error) {
	var stmt *sql.Stmt
	var err error
	if (self.Id > 0) {
		stmt, err = db.Prepare("select Id from jobs where ((JobTitle = $1 AND CompanyName = $2) OR (SourceUrl = $3)) AND Id <> $4 LIMIT 1")
	} else {
		stmt, err = db.Prepare("select Id from jobs where ((JobTitle = $1 AND CompanyName = $2) OR (SourceUrl = $3)) LIMIT 1")
	}
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if (self.Id > 0) {
		rows, err = stmt.Query(self.JobTitle, self.CompanyName, self.SourceUrl, self.Id)
	} else {
		rows, err = stmt.Query(self.JobTitle, self.CompanyName, self.SourceUrl)
	}
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// if any rows are returned, this doesn't pass validation
	for rows.Next() {
		return false, nil
	}
	return true, nil
}

func (self *Job) Create(db *sql.DB) error {
	validationPassed, err := self.PassesCreationValidation(db)
	if err != nil { return err }
	if (!validationPassed) {
		return errors.New("called create on an object that doesn't pass validation")
	}

	stmt, err := db.Prepare("insert into jobs (JobTitle, JobLocation, JobDescription, HowToApply, CompanyLocation, CompanyName, CompanyUrl, SourceUrl, SourceName, PostedAt) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")
	if err != nil { return err }
	defer stmt.Close()

	_, err = stmt.Exec(self.JobTitle, self.JobLocation, self.JobDescription, self.HowToApply, self.CompanyLocation, self.CompanyName, self.CompanyUrl, self.SourceUrl, self.SourceName, self.PostedAt)
	if err != nil {
		return err
	}

	return nil
}