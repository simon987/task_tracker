package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
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
}

type AssignedTasks struct {
	Assignee  uuid.UUID `json:"assignee"`
	TaskCount int64     `json:"task_count"`
}

type ProjectStats struct {
	Project         *Project         `json:"project"`
	NewTaskCount    int64            `json:"new_task_count"`
	FailedTaskCount int64            `json:"failed_task_count"`
	ClosedTaskCount int64            `json:"closed_task_count"`
	Assignees       []*AssignedTasks `json:"assignees"`
}

func (database *Database) SaveProject(project *Project) (int64, error) {
	db := database.getDB()
	id, projectErr := saveProject(project, db)

	return id, projectErr
}

func saveProject(project *Project, db *sql.DB) (int64, error) {

	row := db.QueryRow(`INSERT INTO project (name, git_repo, clone_url, version, priority, motd)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		project.Name, project.GitRepo, project.CloneUrl, project.Version, project.Priority, project.Motd)

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
	err := row.Scan(&project.Id, &project.Priority, &project.Motd, &project.Name, &project.CloneUrl,
		&project.GitRepo, &project.Version)

	return project, err
}

func (database *Database) GetProjectWithRepoName(repoName string) *Project {

	db := database.getDB()
	row := db.QueryRow(`SELECT * FROM project WHERE LOWER(git_repo)=$1`, strings.ToLower(repoName))

	project, err := scanProject(row)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"repoName": repoName,
		}).Warn("Database.getProjectWithRepoName SELECT project NOT FOUND")
		return nil
	}

	return project
}

func (database *Database) UpdateProject(project *Project) {

	db := database.getDB()

	res, err := db.Exec(`UPDATE project 
		SET (priority, name, clone_url, git_repo, version, motd) = ($1,$2,$3,$4,$5,$6) WHERE id=$7`,
		project.Priority, project.Name, project.CloneUrl, project.GitRepo, project.Version, project.Motd, project.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"project":      project,
		"rowsAffected": rowsAffected,
	}).Trace("Database.updateProject UPDATE project")

	return
}

func (database *Database) GetProjectStats(id int64) *ProjectStats {

	db := database.getDB()
	stats := ProjectStats{}

	stats.Project = getProject(id, db)

	if stats.Project != nil {
		row := db.QueryRow(`SELECT 
       SUM(CASE WHEN status='new' THEN 1 ELSE 0 END) newCount,
       SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END) failedCount,
       SUM(CASE WHEN status='closed' THEN 1 ELSE 0 END) closedCount
       FROM task WHERE project=$1 GROUP BY project`, id)

		err := row.Scan(&stats.NewTaskCount, &stats.FailedTaskCount, &stats.ClosedTaskCount)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"id": id,
			}).Trace("Get project stats: No task for this project")
			return nil
		}

		rows, err := db.Query(`SELECT assignee, COUNT(*) FROM TASK
  			LEFT JOIN worker ON TASK.assignee = worker.id WHERE project=$1 GROUP BY assignee`, id)

		for rows.Next() {
			assignee := AssignedTasks{}
			err = rows.Scan(&assignee.Assignee, &assignee.TaskCount)
			handleErr(err)

			stats.Assignees = append(stats.Assignees, &assignee)
		}
	}

	return &stats
}

func (database Database) GetAllProjectsStats() *[]ProjectStats {
	var statsList []ProjectStats

	db := database.getDB()
	rows, err := db.Query(`SELECT 
       	SUM(CASE WHEN status='new' THEN 1 ELSE 0 END) newCount,
       	SUM(CASE WHEN status='failed' THEN 1 ELSE 0 END) failedCount,
		SUM(CASE WHEN status='closed' THEN 1 ELSE 0 END) closedCount,
       	p.*
		FROM task INNER JOIN project p on task.project = p.id
		GROUP BY p.id`)
	handleErr(err)

	for rows.Next() {

		stats := ProjectStats{}
		p := &Project{}
		err := rows.Scan(&stats.NewTaskCount, &stats.FailedTaskCount, &stats.ClosedTaskCount,
			&p.Id, &p.Priority, &p.Motd, &p.Name, &p.CloneUrl, &p.GitRepo, &p.Version)
		handleErr(err)

		stats.Project = p

		statsList = append(statsList, stats)
	}

	logrus.WithFields(logrus.Fields{
		"statsList": statsList,
	}).Trace("Get all projects stats")

	return &statsList
}
