package storage

import (
	"github.com/sirupsen/logrus"
)

type Worker struct {
	Id      int64  `json:"id"`
	Created int64  `json:"created"`
	Alias   string `json:"alias,omitempty"`
	Secret  []byte `json:"secret"`
	Paused  bool   `json:"paused"`
}

type WorkerStats struct {
	Alias           string `json:"alias"`
	ClosedTaskCount int64  `json:"closed_task_count"`
}

type WorkerAccess struct {
	Submit  bool   `json:"submit"`
	Assign  bool   `json:"assign"`
	Request bool   `json:"request"`
	Worker  Worker `json:"worker"`
	Project int64  `json:"project"`
}

func (database *Database) SaveWorker(worker *Worker) {

	db := database.getDB()

	row := db.QueryRow(`INSERT INTO worker (created, secret, alias) 
		VALUES ($1,$2,$3) RETURNING id`,
		worker.Created, worker.Secret, worker.Alias)

	err := row.Scan(&worker.Id)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"newId": worker.Id,
	}).Trace("Database.saveWorker INSERT worker")
}

func (database *Database) GetWorker(id int64) *Worker {

	if database.workerCache[id] != nil {
		return database.workerCache[id]
	}

	db := database.getDB()

	worker := &Worker{}

	row := db.QueryRow("SELECT id, created, secret, alias FROM worker WHERE id=$1", id)
	err := row.Scan(&worker.Id, &worker.Created, &worker.Secret, &worker.Alias)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Database.getWorker SELECT worker NOT FOUND")
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"worker": worker,
	}).Trace("Database.getWorker SELECT worker")

	database.workerCache[id] = worker

	return worker
}

func (database *Database) GrantAccess(workerId int64, projectId int64) bool {

	db := database.getDB()
	res, err := db.Exec(`UPDATE worker_access SET
  		request=FALSE WHERE worker=$1 AND project=$2`,
		workerId, projectId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"workerId":  workerId,
			"projectId": projectId,
		}).WithError(err).Warn("Database.GrantAccess INSERT")
		return false
	}

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"workerId":     workerId,
		"projectId":    projectId,
	}).Trace("Database.GrantAccess INSERT")

	return rowsAffected == 1
}

func (database *Database) UpdateWorker(worker *Worker) bool {

	db := database.getDB()
	res, err := db.Exec(`UPDATE worker SET alias=$1, paused=$2 WHERE id=$3`,
		worker.Alias, worker.Paused, worker.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"worker":       worker,
	}).Trace("Database.UpdateWorker UPDATE worker")

	database.workerCache[worker.Id] = worker

	return rowsAffected == 1
}

func (database *Database) SaveAccessRequest(wa *WorkerAccess) bool {

	db := database.getDB()

	res, err := db.Exec(`INSERT INTO worker_access(worker, project, role_assign,
                          role_submit, request)
 		VALUES ($1,$2,$3,$4,TRUE)`,
		wa.Worker.Id, wa.Project, wa.Assign, wa.Submit)
	if err != nil {
		return false
	}

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.SaveAccessRequest INSERT")

	return rowsAffected == 1
}

func (database *Database) AcceptAccessRequest(worker int64, projectId int64) bool {

	db := database.getDB()

	res, err := db.Exec(`UPDATE worker_access SET request=FALSE 
		WHERE worker=$1 AND project=$2`, worker, projectId)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.AcceptAccessRequest")

	return rowsAffected == 1
}

func (database *Database) RejectAccessRequest(workerId int64, projectId int64) bool {

	db := database.getDB()
	res, err := db.Exec(`DELETE FROM worker_access WHERE worker=$1 AND project=$2`,
		workerId, projectId)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"workerId":     workerId,
		"projectId":    projectId,
	}).Trace("Database.RejectAccessRequest DELETE")

	return rowsAffected == 1
}

func (database *Database) GetAllAccesses(projectId int64) *[]WorkerAccess {

	db := database.getDB()

	rows, err := db.Query(`SELECT id, alias, created, role_assign, role_submit, request
		FROM worker_access
		INNER JOIN worker w on worker_access.worker = w.id
		WHERE project=$1 ORDER BY request, alias`,
		projectId)
	handleErr(err)

	requests := make([]WorkerAccess, 0)

	for rows.Next() {
		wa := WorkerAccess{
			Project: projectId,
		}
		_ = rows.Scan(&wa.Worker.Id, &wa.Worker.Alias, &wa.Worker.Created,
			&wa.Assign, &wa.Submit, &wa.Request)
		requests = append(requests, wa)
	}

	return &requests
}

func (database *Database) GetAllWorkerStats() *[]WorkerStats {

	db := database.getDB()
	rows, err := db.Query(`SELECT alias, closed_task_count FROM worker WHERE closed_task_count>0 LIMIT 50`)
	handleErr(err)

	stats := make([]WorkerStats, 0)
	for rows.Next() {
		s := WorkerStats{}
		_ = rows.Scan(&s.Alias, &s.ClosedTaskCount)
		stats = append(stats, s)
	}

	return &stats
}
