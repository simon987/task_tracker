package storage

import (
	"database/sql"
	"fmt"
	"github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"src/task_tracker/config"
)

type Database struct {
}

func (database *Database) Reset() {

	file, err := os.Open("./schema.sql")
	handleErr(err)

	buffer, err := ioutil.ReadAll(file)
	handleErr(err)

	db := database.getDB()
	_, err = db.Exec(string(buffer))
	handleErr(err)

	db.Close()
	file.Close()

	logrus.Info("Database has been reset")
}

func (database *Database) getDB () *sql.DB {
	db, err := sql.Open("postgres", config.Cfg.DbConnStr)
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}

func (database *Database) Test() {

	db := database.getDB()

	rows, err := db.Query("SELECT name FROM Task")
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(rows)

}