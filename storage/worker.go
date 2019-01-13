package storage

import (
	"database/sql"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

type Identity struct {
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
}

type Worker struct {
	Id       uuid.UUID `json:"id"`
	Created  int64     `json:"created"`
	Identity *Identity `json:"identity"`
}

func (database *Database) SaveWorker(worker *Worker) {

	db := database.getDB()
	saveWorker(worker, db)
	err := db.Close()
	handleErr(err)
}

func (database *Database) GetWorker(id uuid.UUID) *Worker {

	db := database.getDB()
	worker := getWorker(id, db)
	err := db.Close()
	handleErr(err)
	return worker
}

func saveWorker(worker *Worker, db *sql.DB) {

	identityId := getOrCreateIdentity(worker.Identity, db)

	res, err := db.Exec("INSERT INTO worker (id, created, identity) VALUES ($1,$2,$3)",
		worker.Id, worker.Created, identityId)
	handleErr(err)

	var rowsAffected, _ = res.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.saveWorker INSERT worker")
}

func getWorker(id uuid.UUID, db *sql.DB) *Worker {

	worker := &Worker{}
	var identityId int64

	row := db.QueryRow("SELECT id, created, identity FROM worker WHERE id=$1", id)
	err := row.Scan(&worker.Id, &worker.Created, &identityId)
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

	row := db.QueryRow("SELECT remote_addr, user_agent FROM workeridentity WHERE id=$1", id)
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

	res, err := db.Exec("INSERT INTO workeridentity (remote_addr, user_agent) VALUES ($1,$2) ON CONFLICT DO NOTHING",
		identity.RemoteAddr, identity.UserAgent)
	handleErr(err)

	rowsAffected, err := res.RowsAffected()
	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Trace("Database.saveWorker INSERT workerIdentity")

	row := db.QueryRow("SELECT (id) FROM workeridentity WHERE remote_addr=$1", identity.RemoteAddr)

	var rowId int64
	err = row.Scan(&rowId)
	handleErr(err)

	logrus.WithFields(logrus.Fields{
		"rowId": rowId,
	}).Trace("Database.saveWorker SELECT workerIdentity")

	return rowId
}
