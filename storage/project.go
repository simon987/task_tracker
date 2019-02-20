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
	Hidden   bool   `json:"hidden"`
	Chain    int64  `json:"chain"`
}

type AssignedTasks struct {
	Assignee  string `json:"assignee"`
	TaskCount int64  `json:"task_count"`
}

func (database *Database) SaveProject(project *Project) (int64, error) {
	db := database.getDB()

	row := db.QueryRow(`INSERT INTO project (name, git_repo, clone_url, version, priority,
                     motd, public, hidden, chain)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9, 0)) RETURNING id`,
		project.Name, project.GitRepo, project.CloneUrl, project.Version, project.Priority, project.Motd,
		project.Public, project.Hidden, project.Chain)

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
	row := db.QueryRow(`SELECT id, priority, name, clone_url, git_repo, version,
       motd, public, hidden, COALESCE(chain, 0)
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

	p := &Project{}
	err := row.Scan(&p.Id, &p.Priority, &p.Name, &p.CloneUrl, &p.GitRepo, &p.Version,
		&p.Motd, &p.Public, &p.Hidden, &p.Chain)

	return p, err
}

func (database *Database) GetProjectWithRepoName(repoName string) *Project {

	db := database.getDB()
	row := db.QueryRow(`SELECT id, priority, name, clone_url, git_repo, version,
       motd, public, hidden, COALESCE(chain, 0) FROM project WHERE LOWER(git_repo)=$1`,
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
		SET (priority, name, clone_url, git_repo, version, motd, public, hidden, chain) =
		  ($1,$2,$3,$4,$5,$6,$7,$8,NULLIF($9, 0))
		WHERE id=$10`,
		project.Priority, project.Name, project.CloneUrl, project.GitRepo, project.Version, project.Motd,
		project.Public, project.Hidden, project.Chain, project.Id)
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

func (database Database) GetAllProjects(managerId int64) *[]Project {
	projects := make([]Project, 0)

	db := database.getDB()
	var rows *sql.Rows
	var err error
	if managerId == 0 {
		rows, err = db.Query(`SELECT 
       	Id, priority, name, clone_url, git_repo, version, motd, public, hidden, COALESCE(chain,0)
		FROM project
		WHERE NOT hidden
		ORDER BY name`)
	} else {
		rows, err = db.Query(`SELECT 
       	Id, priority, name, clone_url, git_repo, version, motd, public, hidden, COALESCE(chain,0)
		FROM project
		LEFT JOIN manager_has_role_on_project mhrop ON mhrop.project = id AND mhrop.manager=$1
		WHERE NOT hidden OR mhrop.role & 1 = 1 OR (SELECT tracker_admin FROM manager WHERE id=$1)
		ORDER BY name`, managerId)
	}
	handleErr(err)

	for rows.Next() {
		p := Project{}
		err := rows.Scan(&p.Id, &p.Priority, &p.Name, &p.CloneUrl,
			&p.GitRepo, &p.Version, &p.Motd, &p.Public, &p.Hidden,
			&p.Chain)
		handleErr(err)
		projects = append(projects, p)
	}

	logrus.WithFields(logrus.Fields{
		"projects": len(projects),
	}).Trace("Get all projects stats")

	return &projects
}

func (database *Database) GetAssigneeStats(pid int64, count int64) *[]AssignedTasks {

	db := database.getDB()
	assignees := make([]AssignedTasks, 0)

	rows, err := db.Query(`SELECT worker.alias, COUNT(*) as wc FROM TASK
  			LEFT JOIN worker ON TASK.assignee = worker.id WHERE project=$1 
			GROUP BY worker.id ORDER BY wc LIMIT $2`, pid, count)
	handleErr(err)

	for rows.Next() {
		assignee := AssignedTasks{}
		var assigneeAlias sql.NullString
		err = rows.Scan(&assigneeAlias, &assignee.TaskCount)
		handleErr(err)

		if assigneeAlias.Valid {
			assignee.Assignee = assigneeAlias.String
		} else {
			assignee.Assignee = "unassigned"
		}

		assignees = append(assignees, assignee)
	}

	return &assignees
}

func (database *Database) GetSecret(pid int64, workerId int64) (secret string, err error) {

	db := database.getDB()

	var row *sql.Row
	if workerId == 0 {
		row = db.QueryRow(`SELECT secret FROM project WHERE id=$1`, pid)
	} else {
		row = db.QueryRow(`SELECT secret FROM project 
		WHERE id =$1 AND (
		  	SELECT a.role_assign FROM worker_access a WHERE a.worker=$2 AND a.project=$1
		  )`, pid, workerId)
	}

	err = row.Scan(&secret)
	handleErr(err)
	return
}

func (database *Database) SetSecret(pid int64, secret string) {

	db := database.getDB()

	res, err := db.Exec(`UPDATE project SET secret=$1 WHERE id=$2`, secret, pid)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"project":      pid,
	}).Info("Set secret")
}
