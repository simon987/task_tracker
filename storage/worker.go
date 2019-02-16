package storage

import (
	"github.com/Sirupsen/logrus"
)

type Worker struct {
	Id      int64  `json:"id"`
	Created int64  `json:"created"`
	Alias   string `json:"alias,omitempty"`
	Secret  []byte `json:"secret"`
}

type WorkerStats struct {
	Alias           string `json:"alias"`
	ClosedTaskCount int64  `json:"closed_task_count"`
}

func (database *Database) SaveWorker(worker *Worker) {

	db := database.getDB()

	row := db.QueryRow("INSERT INTO worker (created, secret, alias) VALUES ($1,$2,$3) RETURNING id",
		worker.Created, worker.Secret, worker.Alias)

	err := row.Scan(&worker.Id)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"newId": worker.Id,
	}).Trace("Database.saveWorker INSERT worker")
}

func (database *Database) GetWorker(id int64) *Worker {

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

	return worker
}

func (database *Database) GrantAccess(workerId int64, projectId int64) bool {

	db := database.getDB()
	res, err := db.Exec(`INSERT INTO worker_has_access_to_project (worker, projectChange) VALUES ($1,$2)
		ON CONFLICT DO NOTHING`,
		workerId, projectId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"workerId":  workerId,
			"projectId": projectId,
		}).WithError(err).Warn("Database.GrantAccess INSERT worker_hase_access_to_project")
		return false
	}

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"workerId":     workerId,
		"projectId":    projectId,
	}).Trace("Database.GrantAccess INSERT worker_has_access_to_project")

	return rowsAffected == 1
}

func (database *Database) RemoveAccess(workerId int64, projectId int64) bool {

	db := database.getDB()
	res, err := db.Exec(`DELETE FROM worker_has_access_to_project WHERE worker=$1 AND projectChange=$2`,
		workerId, projectId)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"workerId":     workerId,
		"projectId":    projectId,
	}).Trace("Database.RemoveAccess DELETE worker_has_access_to_project")

	return rowsAffected == 1
}

func (database *Database) UpdateWorker(worker *Worker) bool {

	db := database.getDB()
	res, err := db.Exec(`UPDATE worker SET alias=$1 WHERE id=$2`,
		worker.Alias, worker.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"worker":       worker,
	}).Trace("Database.UpdateWorker UPDATE worker")

	return rowsAffected == 1
}

func (database *Database) SaveAccessRequest(worker *Worker, projectId int64) bool {

	db := database.getDB()

	res, err := db.Exec(`INSERT INTO worker_requests_access_to_project 
  		SELECT $1, id FROM projectChange WHERE id=$2 AND NOT projectChange.public 
		AND NOT EXISTS(SELECT * FROM worker_has_access_to_project WHERE worker=$1 AND projectChange=$2)`,
		worker.Id, projectId)
	if err != nil {
		return false
	}

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.SaveAccessRequest INSERT")

	return rowsAffected == 1
}

func (database *Database) AcceptAccessRequest(worker *Worker, projectId int64) bool {

	db := database.getDB()

	res, err := db.Exec(`DELETE FROM worker_requests_access_to_project 
		WHERE worker=$1 AND projectChange=$2`)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 1 {
		_, err := db.Exec(`INSERT INTO worker_has_access_to_project 
  			(worker, projectChange) VALUES ($1,$2)`,
			worker.Id, projectId)
		handleErr(err)
	}

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.AcceptAccessRequest")

	return rowsAffected == 1
}

func (database *Database) RejectAccessRequest(worker *Worker, projectId int64) bool {

	db := database.getDB()

	res, err := db.Exec(`DELETE FROM worker_requests_access_to_project 
		  WHERE worker=$1 AND projectChange=$2`, worker.Id, projectId)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.AcceptAccessRequest")

	return rowsAffected == 1
}

func (database *Database) GetAllAccessRequests(projectId int64) *[]Worker {

	db := database.getDB()

	rows, err := db.Query(`SELECT id, alias, created FROM worker_requests_access_to_project
		INNER JOIN worker w on worker_requests_access_to_project.worker = w.id
		WHERE projectChange=$1`,
		projectId)
	handleErr(err)

	requests := make([]Worker, 0)

	for rows.Next() {
		w := Worker{}
		_ = rows.Scan(&w.Id, &w.Alias, &w.Created)
		requests = append(requests, w)
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
