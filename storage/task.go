package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

type Task struct {
	Id         int64
	Project    int64
	Assignee   uuid.UUID
	Retries    int64
	MaxRetries int64
	Status     string
	Recipe     string
}

func (database *Database) SaveTask(task *Task) error {

	db := database.getDB()
	taskErr := saveTask(task, db)
	err := db.Close()
	handleErr(err)

	return taskErr
}

func saveTask(task *Task, db *sql.DB) error {

	res, err := db.Exec("INSERT INTO task (project, max_retries, recipe) "+
		"VALUES ($1,$2,$3)",
		task.Project, task.MaxRetries, task.Recipe)
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
