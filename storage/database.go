package storage

import (
	"database/sql"
	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"src/task_tracker/config"
)

type Database struct {
	db           *sql.DB
	saveTaskStmt *sql.Stmt
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

		database.db = db
	} else {
		err := database.db.Ping()
		handleErr(err)
	}

	return database.db
}
