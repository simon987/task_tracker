package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

type Task struct {
	Id         int64     `json:"id"`
	Priority   int64     `json:"priority"`
	Project    *Project  `json:"project"`
	Assignee   uuid.UUID `json:"assignee"`
	Retries    int64     `json:"retries"`
	MaxRetries int64     `json:"max_retries"`
	Status     string    `json:"status"`
	Recipe     string    `json:"recipe"`
}

func (database *Database) SaveTask(task *Task, project int64) error {

	db := database.getDB()
	taskErr := saveTask(task, project, db)
	err := db.Close()
	handleErr(err)

	return taskErr
}

func saveTask(task *Task, project int64, db *sql.DB) error {

	res, err := db.Exec(`
	INSERT INTO task (project, max_retries, recipe, priority) 
	VALUES ($1,$2,$3,$4)`,
		project, task.MaxRetries, task.Recipe, task.Priority)
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
	task := getTask(worker, db)
	err := db.Close()
	handleErr(err)

	return task
}

func getTask(worker *Worker, db *sql.DB) *Task {

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project p on task.project = p.id
	WHERE assignee IS NULL
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
	SELECT * FROM task 
	  INNER JOIN project ON task.project = project.id
	WHERE task.id=$1`, id)
	task := scanTask(row)

	logrus.WithFields(logrus.Fields{
		"id":   id,
		"task": task,
	}).Trace("Database.getTaskById SELECT task")

	return task
}

func (database *Database) GetTaskFromProject(worker *Worker, project int64) *Task {

	db := database.getDB()
	task := getTaskFromProject(worker, project, db)
	err := db.Close()
	handleErr(err)

	return task
}

func getTaskFromProject(worker *Worker, projectId int64, db *sql.DB) *Task {

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project p on task.project = p.id
	WHERE assignee IS NULL AND p.id=$2
	ORDER BY p.priority DESC, task.priority DESC
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
		&task.Retries, &task.MaxRetries, &task.Status, &task.Recipe, &project.Id,
		&project.Priority, &project.Name, &project.GitUrl, &project.Version)
	handleErr(err)

	return task
}
