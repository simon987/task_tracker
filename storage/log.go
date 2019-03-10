package storage

import (
	"encoding/json"
	"github.com/simon987/task_tracker/config"
	"github.com/sirupsen/logrus"
)

type LogEntry struct {
	Message   string   `json:"message"`
	Timestamp int64    `json:"timestamp"`
	Data      string   `json:"data"`
	Level     LogLevel `json:"level"`
}

type LogLevel int

const (
	FATAL LogLevel = 1
	PANIC LogLevel = 2
	ERROR LogLevel = 3
	WARN  LogLevel = 4
	INFO  LogLevel = 5
	DEBUG LogLevel = 6
	TRACE LogLevel = 7
)

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

	var logLevel LogLevel

	switch entry.Level {
	case logrus.TraceLevel:
		logLevel = TRACE
	case logrus.DebugLevel:
		logLevel = DEBUG
	case logrus.InfoLevel:
		logLevel = INFO
	case logrus.WarnLevel:
		logLevel = WARN
	case logrus.ErrorLevel:
		logLevel = ERROR
	case logrus.FatalLevel:
		logLevel = FATAL
	case logrus.PanicLevel:
		logLevel = PANIC
	}

	_, err = db.Exec("INSERT INTO log_entry (message, level, message_data, timestamp) VALUES ($1,$2,$3,$4)",
		entry.Message, logLevel, jsonData, entry.Time.Unix())
	return err
}

func (database *Database) SetupLoggerHook() {
	hook := sqlLogHook{}
	hook.database = database
	logrus.AddHook(hook)
}

func (database *Database) GetLogs(since int64, level LogLevel) *[]LogEntry {

	db := database.getDB()

	var logs []LogEntry

	rows, err := db.Query("SELECT * FROM log_entry WHERE timestamp > $1 AND level=$2",
		since, level)
	handleErr(err)

	for rows.Next() {

		e := LogEntry{}

		err := rows.Scan(&e.Level, &e.Message, &e.Data, &e.Timestamp)
		handleErr(err)

		logs = append(logs, e)
	}

	return &logs
}
