package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
)

type Project struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	GitUrl string `json:"git_url"`
	Version string `json:"version"`
}

func (database *Database) SaveProject(project *Project) (int64, error) {
	db := database.getDB()
	id, projectErr := saveProject(project, db)
	err := db.Close()
	handleErr(err)

	return id, projectErr
}

func saveProject(project *Project, db *sql.DB) (int64, error) {

	row := db.QueryRow("INSERT INTO project (name, git_url, version) VALUES ($1,$2,$3) RETURNING id",
		project.Name, project.GitUrl, project.Version)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project": project,
		}).Warn("Database.saveProject INSERT project ERROR")
		return -1, err
	}

	logrus.WithFields(logrus.Fields{
		"id": id,
		"project": project,
	}).Trace("Database.saveProject INSERT project")

	return id, nil
}

func (database *Database) GetProject(id int64) *Project {

	db := database.getDB()
	project := getProject(id, db)
	err := db.Close()
	handleErr(err)
	return project
}

func getProject(id int64, db *sql.DB) *Project {

	project := &Project{}

	row := db.QueryRow("SELECT id, name, git_url, version FROM project WHERE id=$1",
		id)

	err := row.Scan(&project.Id, &project.Name, &project.GitUrl, &project.Version)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Database.getProject SELECT project NOT FOUND")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"id": id,
		"project": project,
	}).Trace("Database.saveProject SELECT project")

	return project
}
