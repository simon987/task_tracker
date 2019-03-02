package storage

import (
	"github.com/Sirupsen/logrus"
)

func (database *Database) ResetFailedTasks(pid int64) int64 {

	db := database.getDB()

	res, err := db.Exec(`UPDATE task SET status=1, retries=0, assign_time=NULL, assignee=NULL 
		WHERE project=$1 AND status=2`, pid)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected
}

func (database *Database) ResetTimedOutTasks() {

	db := database.getDB()

	res, err := db.Exec(`
		UPDATE task SET assignee=NULL, assign_time=NULL
		WHERE status=1 AND assignee IS NOT NULL
		AND extract(epoch from now() at time zone 'utc') > (assign_time + max_assign_time);`)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
	}).Info("Reset timed out tasks")
}
