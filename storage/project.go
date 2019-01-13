package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"strings"
)

type Project struct {
	Id       int64  `json:"id"`
	Priority int64  `json:"priority"`
	Name     string `json:"name"`
	CloneUrl string `json:"clone_url"`
	GitRepo  string `json:"git_repo"`
	Version  string `json:"version"`
}

func (database *Database) SaveProject(project *Project) (int64, error) {
	db := database.getDB()
	id, projectErr := saveProject(project, db)
	err := db.Close()
	handleErr(err)

	return id, projectErr
}

func saveProject(project *Project, db *sql.DB) (int64, error) {

	row := db.QueryRow(`INSERT INTO project (name, git_repo, clone_url, version, priority)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		project.Name, project.GitRepo, project.CloneUrl, project.Version, project.Priority)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project": project,
		}).Warn("Database.saveProject INSERT project ERROR")
		return -1, err
	}

	logrus.WithFields(logrus.Fields{
		"id":      id,
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

	row := db.QueryRow(`SELECT * FROM project WHERE id=$1`, id)

	project, err := scanProject(row)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"id": id,
		}).Warn("Database.getProject SELECT project NOT FOUND")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"id":      id,
		"project": project,
	}).Trace("Database.saveProject SELECT project")

	return project
}

func scanProject(row *sql.Row) (*Project, error) {

	project := &Project{}
	err := row.Scan(&project.Id, &project.Priority, &project.Name, &project.CloneUrl, &project.GitRepo,
		&project.Version)

	return project, err
}

func (database *Database) GetProjectWithRepoName(repoName string) *Project {

	db := database.getDB()
	project := getProjectWithRepoName(repoName, db)
	err := db.Close()
	handleErr(err)
	return project
}

func getProjectWithRepoName(repoName string, db *sql.DB) *Project {

	row := db.QueryRow(`SELECT * FROm project WHERE LOWER(git_repo)=$1`, strings.ToLower(repoName))

	project, err := scanProject(row)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"repoName": repoName,
		}).Error("Database.getProjectWithRepoName SELECT project NOT FOUND")
		return nil
	}

	return project
}

func (database *Database) UpdateProject(project *Project) {

	db := database.getDB()
	updateProject(project, db)
	err := db.Close()
	handleErr(err)
}

func updateProject(project *Project, db *sql.DB) {

	res, err := db.Exec(`UPDATE project 
		SET (priority, name, clone_url, git_repo, version) = ($1,$2,$3,$4,$5) WHERE id=$6`,
		project.Priority, project.Name, project.CloneUrl, project.GitRepo, project.Version, project.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"project":      project,
		"rowsAffected": rowsAffected,
	}).Trace("Database.updateProject UPDATE project")

	return
}
