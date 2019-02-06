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
	Motd     string `json:"motd"`
	Public   bool   `json:"public"`
}

type AssignedTasks struct {
	Assignee  string `json:"assignee"`
	TaskCount int64  `json:"task_count"`
}

func (database *Database) SaveProject(project *Project) (int64, error) {
	db := database.getDB()
	id, projectErr := saveProject(project, db)

	return id, projectErr
}

func saveProject(project *Project, db *sql.DB) (int64, error) {

	row := db.QueryRow(`INSERT INTO project (name, git_repo, clone_url, version, priority, motd, public)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		project.Name, project.GitRepo, project.CloneUrl, project.Version, project.Priority, project.Motd, project.Public)

	var id int64
	err := row.Scan(&id)

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"project": project,
		}).Warn("Database.saveProject INSERT project ERROR")
		return -1, err
	}

	project.Id = id

	logrus.WithFields(logrus.Fields{
		"id":      id,
		"project": project,
	}).Trace("Database.saveProject INSERT project")

	return id, nil
}

func (database *Database) GetProject(id int64) *Project {

	db := database.getDB()
	project := getProject(id, db)
	return project
}

func getProject(id int64, db *sql.DB) *Project {

	row := db.QueryRow(`SELECT id, priority, name, clone_url, git_repo, version, motd, public
		FROM project WHERE id=$1`, id)

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
	err := row.Scan(&project.Id, &project.Priority, &project.Name, &project.CloneUrl,
		&project.GitRepo, &project.Version, &project.Motd, &project.Public)

	return project, err
}

func (database *Database) GetProjectWithRepoName(repoName string) *Project {

	db := database.getDB()
	row := db.QueryRow(`SELECT id, priority, name, clone_url, git_repo, version, motd, public 
			FROM project WHERE LOWER(git_repo)=$1`,
		strings.ToLower(repoName))

	project, err := scanProject(row)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"repoName": repoName,
		}).Warn("Database.getProjectWithRepoName SELECT project NOT FOUND")
		return nil
	}

	return project
}

func (database *Database) UpdateProject(project *Project) error {

	db := database.getDB()

	res, err := db.Exec(`UPDATE project 
		SET (priority, name, clone_url, git_repo, version, motd, public) = ($1,$2,$3,$4,$5,$6,$7) WHERE id=$8`,
		project.Priority, project.Name, project.CloneUrl, project.GitRepo, project.Version, project.Motd, project.Public, project.Id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"project":      project,
		"rowsAffected": rowsAffected,
	}).Trace("Database.updateProject UPDATE project")

	return nil
}

func (database Database) GetAllProjects() *[]Project {
	var projects []Project

	db := database.getDB()
	rows, err := db.Query(`SELECT 
       	Id, priority, name, clone_url, git_repo, version, motd, public
		FROM project
		ORDER BY name`)
	handleErr(err)

	for rows.Next() {

		p := Project{}
		err := rows.Scan(&p.Id, &p.Priority, &p.Name, &p.CloneUrl,
			&p.GitRepo, &p.Version, &p.Motd, &p.Public)
		handleErr(err)
		projects = append(projects, p)
	}

	logrus.WithFields(logrus.Fields{
		"projects": projects,
	}).Trace("Get all projects stats")

	return &projects
}
