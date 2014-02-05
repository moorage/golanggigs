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
	CreatedAt       time.Time
}

func (self *Job) PassesCreationValidation() (bool, error) {
	var stmt *db.Stmt
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

	var rows *db.Rows
	if (self.Id > 0) {
		rows, err = stmt.Query(self.JobTitle, self.CompanyName, self.SourceUrl, self.Id)
	} else {
		rows, err = stmt.Query(self.JobTitle, self.CompanyName, self.SourceUrl, self.Id)
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

func (self *Job) Create() error {
	if (!self.PassesCreationValidation()) {
		return errors.New("called create on an object that doesn't pass validation")
	}
	
	stmt, err := db.Prepare("insert into jobs (JobTitle, JobLocation, JobDescription, HowToApply, CompanyLocation, CompanyName, CompanyUrl, SourceUrl, SourceName, PostedAt, CreatedAt) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	res, err := stmt.Exec(self.JobTitle, self.JobLocation, self.JobDescription, self.HowToApply, self.CompanyLocation, self.CompanyName, self.CompanyUrl, self.SourceUrl, self.SourceName, self.PostedAt, self.CreatedAt)
	if err != nil {
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	self.Id = lastId
	return nil
}