package storage

import (
	"database/sql"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"src/task_tracker/config"
)

type LogEntry struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
	Level     string `json:"level"`
}
type sqlLogHook struct {
	database *Database
}

func (h sqlLogHook) Levels() []logrus.Level {
	return config.Cfg.DbLogLevels
}

func (h sqlLogHook) Fire(entry *logrus.Entry) error {

	db := h.database.getDB()

	jsonData, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO log_entry (message, level, message_data, timestamp) VALUES ($1,$2,$3,$4)",
		entry.Message, entry.Level.String(), jsonData, entry.Time.Unix())
	return err
}

func (database *Database) SetupLoggerHook() {
	hook := sqlLogHook{}
	hook.database = database
	logrus.AddHook(hook)
}

func (database *Database) GetLogs(since int64, level logrus.Level) *[]LogEntry {

	db := database.getDB()
	logs := getLogs(since, level, db)

	return logs
}

func getLogs(since int64, level logrus.Level, db *sql.DB) *[]LogEntry {

	var logs []LogEntry

	rows, err := db.Query("SELECT * FROM log_entry WHERE timestamp > $1 AND level=$2",
		since, level.String())
	handleErr(err)

	for rows.Next() {

		e := LogEntry{}

		err := rows.Scan(&e.Level, &e.Message, &e.Data, &e.Timestamp)
		handleErr(err)

		logs = append(logs, e)
	}

	return &logs
}
