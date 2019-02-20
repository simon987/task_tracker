package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type Task struct {
	Id                int64      `json:"id"`
	Priority          int64      `json:"priority"`
	Project           *Project   `json:"project"`
	Assignee          int64      `json:"assignee"`
	Retries           int64      `json:"retries"`
	MaxRetries        int64      `json:"max_retries"`
	Status            TaskStatus `json:"status"`
	Recipe            string     `json:"recipe"`
	MaxAssignTime     int64      `json:"max_assign_time"`
	AssignTime        int64      `json:"assign_time"`
	VerificationCount int64      `json:"verification_count"`
}

type TaskStatus int

const (
	NEW    TaskStatus = 1
	FAILED TaskStatus = 2
)

type TaskResult int

const (
	TR_OK   TaskResult = 0
	TR_FAIL TaskResult = 1
	TR_SKIP TaskResult = 2
)

func (database *Database) SaveTask(task *Task, project int64, hash64 int64, wid int64) error {

	db := database.getDB()

	//TODO: For some reason it refuses to insert the 64-bit value unless I do that...
	res, err := db.Exec(fmt.Sprintf(`
	INSERT INTO task (project, max_retries, recipe, priority, max_assign_time, hash64,verification_count) 
	SELECT $1,$2,$3,$4,$5,NULLIF(%d, 0),$6 FROM worker_access 
	WHERE role_submit AND worker=$7 AND project=$1`, hash64),
		project, task.MaxRetries, task.Recipe, task.Priority, task.MaxAssignTime, task.VerificationCount,
		wid)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"task": task,
		}).Trace("Database.saveTask INSERT task ERROR")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"task":         task,
	}).Trace("Database.saveTask INSERT task")

	if rowsAffected == 0 {
		return errors.New("unauthorized task submit")
	}

	return nil
}

func (database *Database) GetTask(worker *Worker) *Task {

	db := database.getDB()

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1, assign_time=extract(epoch from now() at time zone 'utc')
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project project on task.project = project.id
	LEFT JOIN worker_verifies_task wvt on task.id = wvt.task AND wvt.worker=$1
	WHERE assignee IS NULL AND task.status=1
		AND (project.public OR (
		  SELECT a.role_assign FROM worker_access a WHERE a.worker=$1 AND a.project=project.id
		))
		AND wvt.task IS NULL
	ORDER BY project.priority DESC, task.priority DESC
	LIMIT 1
	)
	RETURNING id`, worker.Id)
	var id int64

	err := row.Scan(&id)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
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
	        status, recipe, max_assign_time, assign_time, verification_count, project.priority, project.name,
	       project.clone_url, project.git_repo, project.version, project.motd, project.public, COALESCE(project.chain,0) FROM task 
	  INNER JOIN project project ON task.project = project.id
	WHERE task.id=$1`, id)
	project := &Project{}
	task := &Task{}
	task.Project = project

	err := row.Scan(&task.Id, &task.Priority, &project.Id, &task.Assignee,
		&task.Retries, &task.MaxRetries, &task.Status, &task.Recipe, &task.MaxAssignTime,
		&task.AssignTime, &task.VerificationCount, &project.Priority, &project.Name,
		&project.CloneUrl, &project.GitRepo, &project.Version, &project.Motd, &project.Public,
		&project.Chain)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"id":   id,
		"task": task,
	}).Trace("Database.getTaskById SELECT task")

	return task
}

func (database Database) ReleaseTask(id int64, workerId int64, result TaskResult, verification int64) bool {

	db := database.getDB()

	var taskUpdated bool
	if result == TR_OK {
		row := db.QueryRow(fmt.Sprintf(`SELECT release_task_ok(%d,%d,%d)`, workerId, id, verification))

		err := row.Scan(&taskUpdated)
		handleErr(err)
	} else if result == TR_FAIL {
		res, err := db.Exec(`UPDATE task SET (status, assignee, retries) = 
			(CASE WHEN retries+1 >= max_retries THEN 2 ELSE 1 END, NULL, retries+1)
			WHERE id=$1 AND assignee=$2`, id, workerId)
		handleErr(err)
		rowsAffected, _ := res.RowsAffected()
		taskUpdated = rowsAffected == 1
	} else if result == TR_SKIP {
		res, err := db.Exec(`UPDATE task SET (status, assignee) = (1, NULL)
			WHERE id=$1 AND assignee=$2`, id, workerId)
		handleErr(err)
		rowsAffected, _ := res.RowsAffected()
		taskUpdated = rowsAffected == 1
	}

	logrus.WithFields(logrus.Fields{
		"taskUpdated":  taskUpdated,
		"taskId":       id,
		"workerId":     workerId,
		"verification": verification,
	}).Trace("Database.ReleaseTask")

	return taskUpdated
}

func (database *Database) GetTaskFromProject(worker *Worker, projectId int64) *Task {

	db := database.getDB()

	row := db.QueryRow(`
	UPDATE task
	SET assignee=$1, assign_time=extract(epoch from now() at time zone 'utc')
	WHERE id IN
	(
		SELECT task.id
	FROM task
	INNER JOIN project project on task.project = project.id
	LEFT JOIN worker_verifies_task wvt on task.id = wvt.task AND wvt.worker=$1
	WHERE assignee IS NULL AND project.id=$2 AND status=1
		AND (project.public OR (
		  SELECT a.role_assign FROM worker_access a WHERE a.worker=$1 AND a.project=$2
		))
		AND wvt.task IS NULL
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
