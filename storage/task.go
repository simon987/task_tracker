package storage

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

type Task struct {
	Id            int64     `json:"id"`
	Priority      int64     `json:"priority"`
	Project       *Project  `json:"project"`
	Assignee      uuid.UUID `json:"assignee"`
	Retries       int64     `json:"retries"`
	MaxRetries    int64     `json:"max_retries"`
	Status        string    `json:"status"`
	Recipe        string    `json:"recipe"`
	MaxAssignTime int64     `json:"max_assign_time"`
	AssignTime    int64     `json:"assign_time"`
}

func (database *Database) SaveTask(task *Task, project int64, hash64 int64) error {

	db := database.getDB()

	//TODO: For some reason it refuses to insert the 64-bit value unless I do that...
	res, err := db.Exec(fmt.Sprintf(`
	INSERT INTO task (project, max_retries, recipe, priority, max_assign_time, hash64) 
	VALUES ($1,$2,$3,$4,$5,NULLIF(%d, 0))`, hash64),
		project, task.MaxRetries, task.Recipe, task.Priority, task.MaxAssignTime)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"task": task,
		}).Warn("Database.saveTask INSERT task ERROR")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"task":         task,
	}).Trace("Database.saveTask INSERT task")

	return nil
}

func (database *Database) GetTask(worker *Worker) *Task {

	db := database.getDB()

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project p on task.project = p.id
	WHERE assignee IS NULL AND task.status='new'
		AND (p.public OR EXISTS (
		  SELECT 1 FROM worker_has_access_to_project a WHERE a.worker=$1 AND a.project=p.id
		))
	ORDER BY p.priority DESC, task.priority DESC
	LIMIT 1
	)
	RETURNING id`, worker.Id)
	var id int64

	err := row.Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"worker": worker,
		}).Trace("No task available")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"id":     id,
		"worker": worker,
	}).Trace("Database.getTask UPDATE task")

	task := getTaskById(id, db)

	return task
}

func getTaskById(id int64, db *sql.DB) *Task {

	row := db.QueryRow(`
	SELECT task.id, task.priority, task.project, assignee, retries, max_retries,
	        status, recipe, max_assign_time, assign_time, project.* FROM task 
	  INNER JOIN project ON task.project = project.id
	WHERE task.id=$1`, id)
	task := scanTask(row)

	logrus.WithFields(logrus.Fields{
		"id":   id,
		"task": task,
	}).Trace("Database.getTaskById SELECT task")

	return task
}

func (database Database) ReleaseTask(id int64, workerId *uuid.UUID, success bool) bool {

	db := database.getDB()

	var res sql.Result
	var err error
	if success {
		res, err = db.Exec(`UPDATE task SET (status, assignee) = ('closed', NULL)
		WHERE id=$1 AND task.assignee=$2`, id, workerId)
	} else {
		res, err = db.Exec(`UPDATE task SET (status, assignee, retries) = 
  		(CASE WHEN retries+1 >= max_retries THEN 'failed' ELSE 'new' END, NULL, retries+1)
		WHERE id=$1 AND assignee=$2`, id, workerId)
	}
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	})

	return rowsAffected == 1
}

func (database *Database) GetTaskFromProject(worker *Worker, projectId int64) *Task {

	db := database.getDB()

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project p on task.project = p.id
	WHERE assignee IS NULL AND p.id=$2 AND status='new'
		AND (p.public OR EXISTS (
		  SELECT 1 FROM worker_has_access_to_project a WHERE a.worker=$1 AND a.project=$2
		))
	ORDER BY task.priority DESC
	LIMIT 1
	)
	RETURNING id`, worker.Id, projectId)
	var id int64

	err := row.Scan(&id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"worker": worker,
		}).Trace("No task available")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"id":     id,
		"worker": worker,
	}).Trace("Database.getTask UPDATE task")

	task := getTaskById(id, db)

	return task
}

func scanTask(row *sql.Row) *Task {

	project := &Project{}
	task := &Task{}
	task.Project = project

	err := row.Scan(&task.Id, &task.Priority, &project.Id, &task.Assignee,
		&task.Retries, &task.MaxRetries, &task.Status, &task.Recipe, &task.MaxAssignTime,
		&task.AssignTime, &project.Id, &project.Priority, &project.Name,
		&project.CloneUrl, &project.GitRepo, &project.Version, &project.Motd, &project.Public)
	handleErr(err)

	return task
}
