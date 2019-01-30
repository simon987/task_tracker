package storage

import (
	"database/sql"
	"errors"
	"github.com/Sirupsen/logrus"
)

type Identity struct {
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
}

type Worker struct {
	Id       int64     `json:"id"`
	Created  int64     `json:"created"`
	Identity *Identity `json:"identity"`
	Alias    string    `json:"alias,omitempty"`
	Secret   []byte    `json:"secret"`
}

func (database *Database) SaveWorker(worker *Worker) {

	db := database.getDB()

	identityId := getOrCreateIdentity(worker.Identity, db)

	row := db.QueryRow("INSERT INTO worker (created, identity, secret, alias) VALUES ($1,$2,$3,$4) RETURNING id",
		worker.Created, identityId, worker.Secret, worker.Alias)

	err := row.Scan(&worker.Id)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"newId": worker.Id,
	}).Trace("Database.saveWorker INSERT worker")
}

func (database *Database) GetWorker(id int64) *Worker {

	db := database.getDB()

	worker := &Worker{}
	var identityId int64

	row := db.QueryRow("SELECT id, created, identity, secret, alias FROM worker WHERE id=$1", id)
	err := row.Scan(&worker.Id, &worker.Created, &identityId, &worker.Secret, &worker.Alias)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Database.getWorker SELECT worker NOT FOUND")
		return nil
	}

	worker.Identity, err = getIdentity(identityId, db)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"worker": worker,
	}).Trace("Database.getWorker SELECT worker")

	return worker
}

func getIdentity(id int64, db *sql.DB) (*Identity, error) {

	identity := &Identity{}

	row := db.QueryRow("SELECT remote_addr, user_agent FROM worker_identity WHERE id=$1", id)
	err := row.Scan(&identity.RemoteAddr, &identity.UserAgent)

	if err != nil {
		return nil, errors.New("identity not found")
	}

	logrus.WithFields(logrus.Fields{
		"identity": identity,
	}).Trace("Database.getIdentity SELECT workerIdentity")

	return identity, nil
}

func getOrCreateIdentity(identity *Identity, db *sql.DB) int64 {

	res, err := db.Exec("INSERT INTO worker_identity (remote_addr, user_agent) VALUES ($1,$2) ON CONFLICT DO NOTHING",
		identity.RemoteAddr, identity.UserAgent)
	handleErr(err)

	rowsAffected, err := res.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.saveWorker INSERT workerIdentity")

	row := db.QueryRow("SELECT (id) FROM worker_identity WHERE remote_addr=$1", identity.RemoteAddr)

	var rowId int64
	err = row.Scan(&rowId)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"rowId": rowId,
	}).Trace("Database.saveWorker SELECT workerIdentity")

	return rowId
}

func (database *Database) GrantAccess(workerId int64, projectId int64) bool {

	db := database.getDB()
	res, err := db.Exec(`INSERT INTO worker_has_access_to_project (worker, project) VALUES ($1,$2)
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
	res, err := db.Exec(`DELETE FROM worker_has_access_to_project WHERE worker=$1 AND project=$2`,
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
