package storage

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Task struct {
	Id                int64      `json:"id"`
	Priority          int16      `json:"priority"`
	Project           *Project   `json:"project"`
	Assignee          int64      `json:"assignee"`
	Retries           int16      `json:"retries"`
	MaxRetries        int16      `json:"max_retries"`
	Status            TaskStatus `json:"status"`
	Recipe            string     `json:"recipe"`
	MaxAssignTime     int64      `json:"max_assign_time"`
	AssignTime        int64      `json:"assign_time"`
	VerificationCount int16      `json:"verification_count"`
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

type SaveTaskRequest struct {
	Task     *Task
	Project  int64
	Hash64   int64
	WorkerId int64
}

func (database *Database) checkAccess(workerId, projectId int64, assign, submit bool) bool {

	if database.submitAccessCache[workerId] == nil {
		database.submitAccessCache[workerId] = make(map[int64]bool)
		database.assignAccessCache[workerId] = make(map[int64]bool)
	} else {
		_, ok := database.submitAccessCache[workerId][projectId]
		if ok {
			if assign && !database.assignAccessCache[workerId][projectId] {
				return false
			}
			if submit && !database.submitAccessCache[workerId][projectId] {
				return false
			}
			return true
		}
	}

	db := database.getDB()

	row := db.QueryRow(`SELECT role_assign, role_submit FROM worker_access 
				WHERE worker=$1 and project=$2 AND NOT request`,
		workerId, projectId)

	var hasAssign, hasSubmit bool
	err := row.Scan(&hasAssign, &hasSubmit)

	database.submitAccessCache[workerId][projectId] = hasSubmit
	database.assignAccessCache[workerId][projectId] = hasAssign

	if err != nil {
		return false
	}
	if !hasAssign && assign {
		return false
	}
	if !hasSubmit && submit {
		return false
	}

	return true
}

func (database *Database) SaveTask(task *Task, project int64, hash64 int64, wid int64) error {

	if !database.checkAccess(wid, project, false, true) {
		return errors.New("unauthorized task submit")
	}

	db := database.getDB()

	_, err := db.Exec(`INSERT INTO task 
			(project, max_retries, recipe, priority, max_assign_time, hash64, verification_count) 
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		project, task.MaxRetries, task.Recipe, task.Priority, task.MaxAssignTime,
		makeNullableInt(hash64), task.VerificationCount)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"task": task,
		}).Trace("Database.saveTask INSERT task ERROR")
		return err
	}

	return nil
}

func makeNullableInt(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	} else {
		return sql.NullInt64{
			Valid: true,
			Int64: i,
		}
	}
}

func (database Database) BulkSaveTask(bulkSaveTaskReqs []SaveTaskRequest) []error {

	if !database.checkAccess(bulkSaveTaskReqs[0].WorkerId, bulkSaveTaskReqs[0].Project,
		false, true) {
		return []error{errors.New("unauthorized task submit")}
	}

	db := database.getDB()

	txn, err := db.Begin()
	handleErr(err)

	errs := make([]error, len(bulkSaveTaskReqs))

	stmt, _ := txn.Prepare(pq.CopyIn(
		"task",
		"project", "max_retries", "recipe", "priority",
		"max_assign_time", "hash64", "verification_count",
	))

	for i, req := range bulkSaveTaskReqs {
		_, err = stmt.Exec(req.Project, req.Task.MaxRetries, req.Task.Recipe,
			req.Task.Priority, req.Task.MaxAssignTime, makeNullableInt(req.Hash64),
			req.Task.VerificationCount)
		if err != nil {
			errs[i] = err
		}
	}

	_, err = stmt.Exec()
	err = stmt.Close()
	handleErr(err)
	err = txn.Commit()
	handleErr(err)

	return errs
}

func (database Database) ReleaseTask(id int64, workerId int64, result TaskResult, verification int64) bool {

	db := database.getDB()

	var taskUpdated bool
	if result == TR_OK {
		row := db.QueryRow(`SELECT release_task_ok($1,$2,$3)`, workerId, id, verification)

		err := row.Scan(&taskUpdated)
		if err != nil {
			taskUpdated = false
		}
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

	return taskUpdated
}

func (database *Database) GetTaskFromProject(worker *Worker, projectId int64) *Task {

	db := database.getDB()

	database.assignMutex.Lock()

	row := db.QueryRow(`
		UPDATE task
		SET assignee=$1, assign_time=extract(epoch from now() at time zone 'utc')
		WHERE task.id = (
			SELECT task.id
			FROM task
			INNER JOIN project on task.project = project.id AND project.id=$2 AND not paused
			LEFT JOIN worker_verifies_task wvt on task.id = wvt.task AND wvt.worker=$1
			LEFT JOIN worker_access wa on project.id = wa.project AND wa.worker=$1
			WHERE 
				assignee IS NULL 
				AND status=1
				AND (project.public OR (wa.role_assign AND NOT request))
				AND wvt.task IS NULL
			ORDER BY task.priority DESC
			LIMIT 1
		)
		RETURNING task.id, task.priority, assignee, retries, max_retries,
				status, recipe, max_assign_time, assign_time, verification_count`, worker.Id, projectId)

	database.assignMutex.Unlock()

	task := &Task{}
	err := row.Scan(&task.Id, &task.Priority, &task.Assignee,
		&task.Retries, &task.MaxRetries, &task.Status, &task.Recipe, &task.MaxAssignTime,
		&task.AssignTime, &task.VerificationCount)

	if err != nil {
		return nil
	}

	task.Project = database.GetProject(projectId)

	return task
}
