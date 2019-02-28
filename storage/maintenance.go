package storage

func (database *Database) ResetFailedTasks(pid int64) int64 {

	db := database.getDB()

	res, err := db.Exec(`UPDATE task SET status=1, retries=0, assign_time=NULL, assignee=NULL 
		WHERE project=$1 AND status=2`, pid)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected
}
