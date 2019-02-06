package storage

type ProjectMonitoringSnapshot struct {
	NewTaskCount      int64
	FailedTaskCount   int64
	ClosedTaskCount   int64
	WorkerAccessCount int64
	TimeStamp         int64
}

func (database *Database) MakeProjectSnapshots() {

	db := database.getDB()

	_, err := db.Exec(`
		INSERT INTO project_monitoring_snapshot
		  (project, new_task_count, failed_task_count, closed_task_count, worker_access_count, timestamp)
		SELECT id,
			   (SELECT COUNT(*) FROM task WHERE task.project = project.id AND status = 1),
			   (SELECT COUNT(*) FROM task WHERE task.project = project.id AND status = 2),
			   closed_task_count,
			   (SELECT COUNT(*) FROM worker_has_access_to_project wa WHERE wa.project = project.id),
			   extract(epoch from now() at time zone 'utc')
		FROM project`)
	handleErr(err)
}

func (database *Database) GetMonitoringSnapshots(pid int64, from int64, to int64) (ss *[]ProjectMonitoringSnapshot) {

	db := database.getDB()

	rows, err := db.Query(`SELECT new_task_count, failed_task_count, closed_task_count,
		worker_access_count, timestamp FROM project_monitoring_snapshot 
		WHERE project=$1 AND timestamp BETWEEN $2 AND $3`, pid, from, to)
	handleErr(err)

	for rows.Next() {

		s := ProjectMonitoringSnapshot{}
		err := rows.Scan(&s.NewTaskCount, &s.FailedTaskCount, &s.ClosedTaskCount, &s.WorkerAccessCount, &s.TimeStamp)
		handleErr(err)
	}
	return nil
}
