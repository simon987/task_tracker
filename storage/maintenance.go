package storage

import (
	"github.com/sirupsen/logrus"
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

func (database Database) HardReset(pid int64) int64 {

	db := database.getDB()

	_, err := db.Exec(`UPDATE task SET assignee=NULL WHERE project=$1`, pid)
	handleErr(err)
	res, err := db.Exec(`DELETE FROM task WHERE project=$1`, pid)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"project":      pid,
	}).Info("Hard reset")

	return rowsAffected
}
