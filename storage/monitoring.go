package storage

import (
	"github.com/simon987/task_tracker/config"
	"github.com/sirupsen/logrus"
	"time"
)

type ProjectMonitoringSnapshot struct {
	NewTaskCount              int64 `json:"new_task_count"`
	FailedTaskCount           int64 `json:"failed_task_count"`
	ClosedTaskCount           int64 `json:"closed_task_count"`
	WorkerAccessCount         int64 `json:"worker_access_count"`
	AwaitingVerificationCount int64 `json:"awaiting_verification_count"`
	TimeStamp                 int64 `json:"time_stamp"`
}

func (database *Database) MakeProjectSnapshots() {

	startTime := time.Now()
	db := database.getDB()

	insertRes, err := db.Exec(`
		INSERT INTO project_monitoring_snapshot
		  (project, new_task_count, failed_task_count, closed_task_count, worker_access_count,
		   awaiting_verification_task_count, timestamp)
		SELECT id,
			   (SELECT COUNT(*) FROM task 
					LEFT JOIN worker_verifies_task wvt on task.id = wvt.task
			   		WHERE task.project = project.id AND status = 1 AND wvt.task IS NULL),
			   (SELECT COUNT(*) FROM task WHERE task.project = project.id AND status = 2),
			   closed_task_count,
			   (SELECT COUNT(*) FROM worker_access wa WHERE wa.project = project.id),
			   (SELECT COUNT(*) FROM worker_verifies_task INNER JOIN task t on worker_verifies_task.task = t.id
			  		WHERE t.project = project.id),
			   extract(epoch from now() at time zone 'utc')
		FROM project`)
	handleErr(err)
	inserted, _ := insertRes.RowsAffected()

	res, err := db.Exec(`DELETE FROM project_monitoring_snapshot WHERE timestamp < $1`,
		int64(time.Now().Unix())-int64(config.Cfg.MonitoringHistory.Seconds()))
	handleErr(err)
	deleted, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"took":   time.Now().Sub(startTime),
		"add":    inserted,
		"remove": deleted,
	}).Trace("Took project monitoring snapshot")
}

func (database *Database) GetMonitoringSnapshotsBetween(pid int64, from int, to int) (ss *[]ProjectMonitoringSnapshot) {

	db := database.getDB()

	snapshots := make([]ProjectMonitoringSnapshot, 0)

	rows, err := db.Query(`SELECT new_task_count, failed_task_count, closed_task_count,
		worker_access_count, awaiting_verification_task_count, timestamp FROM project_monitoring_snapshot 
		WHERE project=$1 AND timestamp BETWEEN $2 AND $3 ORDER BY TIMESTAMP DESC `, pid, from, to)
	handleErr(err)

	for rows.Next() {

		s := ProjectMonitoringSnapshot{}
		err := rows.Scan(&s.NewTaskCount, &s.FailedTaskCount, &s.ClosedTaskCount, &s.WorkerAccessCount,
			&s.AwaitingVerificationCount, &s.TimeStamp)
		handleErr(err)

		snapshots = append(snapshots, s)
	}

	logrus.WithFields(logrus.Fields{
		"snapshotCount": len(snapshots),
		"projectId":     pid,
		"from":          from,
		"to":            to,
	}).Trace("Database.GetMonitoringSnapshotsBetween SELECT")

	return &snapshots
}

func (database *Database) GetNMonitoringSnapshots(pid int64, count int) (ss *[]ProjectMonitoringSnapshot) {

	db := database.getDB()

	snapshots := make([]ProjectMonitoringSnapshot, 0)

	rows, err := db.Query(`SELECT new_task_count, failed_task_count, closed_task_count,
		worker_access_count, awaiting_verification_task_count, timestamp FROM project_monitoring_snapshot 
		WHERE project=$1 ORDER BY TIMESTAMP DESC LIMIT $2`, pid, count)
	handleErr(err)

	for rows.Next() {
		s := ProjectMonitoringSnapshot{}
		err := rows.Scan(&s.NewTaskCount, &s.FailedTaskCount, &s.ClosedTaskCount, &s.WorkerAccessCount,
			&s.AwaitingVerificationCount, &s.TimeStamp)
		handleErr(err)

		snapshots = append(snapshots, s)
	}

	logrus.WithFields(logrus.Fields{
		"snapshotCount": len(snapshots),
		"projectId":     pid,
		"count":         count,
	}).Trace("Database.GetNMonitoringSnapshots SELECT")

	return &snapshots
}
