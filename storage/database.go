package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/simon987/task_tracker/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sync"
)

type Database struct {
	db           *sql.DB
	saveTaskStmt *sql.Stmt

	workerCache  map[int64]*Worker
	projectCache map[int64]*Project

	// [worker][project]
	submitAccessCache map[int64]map[int64]bool
	assignAccessCache map[int64]map[int64]bool

	assignMutex *sync.Mutex
}

func New() *Database {

	d := Database{}
	d.workerCache = make(map[int64]*Worker)
	d.projectCache = make(map[int64]*Project)
	d.assignAccessCache = make(map[int64]map[int64]bool)
	d.submitAccessCache = make(map[int64]map[int64]bool)
	d.assignMutex = &sync.Mutex{}

	d.init()

	return &d
}

func (database *Database) init() {

	db := database.getDB()

	_, err := db.Exec(`SELECT * FROM project`)
	if err != nil {
		logrus.Info("Database first time setup")
		database.Reset()
	}
}

func (database *Database) Reset() {

	file, err := os.Open("./schema.sql")
	handleErr(err)

	buffer, err := ioutil.ReadAll(file)
	handleErr(err)

	db := database.getDB()
	_, err = db.Exec(string(buffer))
	handleErr(err)

	file.Close()

	logrus.Info("Database has been reset")
}

func (database *Database) getDB() *sql.DB {

	if database.db == nil {
		db, err := sql.Open("postgres", config.Cfg.DbConnStr)
		if err != nil {
			logrus.Fatal(err)
		}

		db.SetMaxOpenConns(50)

		database.db = db
	}

	return database.db
}

func (database *Database) invalidateAccessCache(worker int64) {
	delete(database.submitAccessCache, worker)
	delete(database.assignAccessCache, worker)
}
