package storage

import (
	"bytes"
	"crypto"
	"errors"
	"github.com/Sirupsen/logrus"
)

type ManagerRole int

const (
	ROLE_NONE          ManagerRole = 0
	ROLE_READ          ManagerRole = 1
	ROLE_EDIT          ManagerRole = 2
	ROLE_MANAGE_ACCESS ManagerRole = 4
)

type Manager struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	WebsiteAdmin bool   `json:"tracker_admin"`
	RegisterTime int64  `json:"register_time"`
}

func (database *Database) ValidateCredentials(username []byte, password []byte) (*Manager, error) {

	db := database.getDB()

	row := db.QueryRow(`SELECT id, password, tracker_admin, register_time FROM manager WHERE username=$1`,
		username)

	manager := &Manager{}
	var passwordHash []byte
	err := row.Scan(&manager.Id, &passwordHash, &manager.WebsiteAdmin, &manager.RegisterTime)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": username,
		}).Warning("Database.ValidateCredentials: user not found")

		return nil, errors.New("username not found")
	}

	hash := crypto.SHA512.New()
	hash.Write([]byte(password))
	hash.Write([]byte(username))

	if bytes.Compare(passwordHash, hash.Sum(nil)) != 0 {
		logrus.WithFields(logrus.Fields{
			"username": username,
		}).Warning("Database.ValidateCredentials: password does not match")

		return nil, errors.New("password does not match")
	}

	manager.Username = string(username)

	return manager, nil
}

func (database *Database) SaveManager(manager *Manager, password []byte) error {

	db := database.getDB()

	hash := crypto.SHA512.New()
	hash.Write(password)
	hash.Write([]byte(manager.Username))
	hashedPassword := hash.Sum(nil)

	row := db.QueryRow(`INSERT INTO manager (username, password, tracker_admin, register_time) 
			VALUES ($1,$2,$3, extract(epoch from now() at time zone 'utc')) RETURNING ID, register_time`,
		manager.Username, hashedPassword, manager.WebsiteAdmin)

	err := row.Scan(&manager.Id, &manager.RegisterTime)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": manager,
		}).Warning("Database.SaveManager INSERT error")

		return err
	}

	manager.WebsiteAdmin = manager.Id == 1

	logrus.WithFields(logrus.Fields{
		"manager": manager,
	}).Info("Database.SaveManager INSERT")

	return nil
}

func (database *Database) UpdateManager(manager *Manager) {

	db := database.getDB()

	res, err := db.Exec(`UPDATE manager SET tracker_admin=$1 WHERE id=$2`,
		manager.WebsiteAdmin, manager.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"manager":      manager,
	}).Trace("Database.UpdateManager UPDATE")
}

func (database *Database) UpdateManagerPassword(manager *Manager, newPassword []byte) {

	hash := crypto.SHA512.New()
	hash.Write(newPassword)
	hash.Write([]byte(manager.Username))
	hashedPassword := hash.Sum(nil)

	db := database.getDB()

	res, err := db.Exec(`UPDATE manager SET password=$1 WHERE id=$2`,
		hashedPassword, manager.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"id":           manager.Id,
	}).Trace("Database.UpdateManagerPassword UPDATE")
}

func (database *Database) GetManagerRoleOn(manager *Manager, projectId int64) ManagerRole {

	db := database.getDB()

	row := db.QueryRow(`SELECT role FROM manager_has_role_on_project 
		WHERE project=$1 AND manager=$2`, projectId, manager.Id)

	var role ManagerRole
	err := row.Scan(&role)
	if err != nil {
		return ROLE_NONE
	}

	return role
}

func (database *Database) SetManagerRoleOn(manager *Manager, projectId int64, role ManagerRole) {

	db := database.getDB()

	res, err := db.Exec(`INSERT INTO manager_has_role_on_project (manager, role, project) 
		VALUES ($1,$2,$3) ON CONFLICT (manager, project) DO UPDATE SET role=$2`,
		manager.Id, role, projectId)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithFields(logrus.Fields{
		"role":         role,
		"manager":      manager.Username,
		"rowsAffected": rowsAffected,
		"project":      projectId,
	}).Info("Set manager role on project")
}

func (database *Database) GetManagerList() *[]Manager {

	db := database.getDB()

	rows, _ := db.Query(`SELECT id, register_time, tracker_admin, username FROM manager`)

	managers := make([]Manager, 0)

	for rows.Next() {
		m := Manager{}
		_ = rows.Scan(&m.Id, &m.RegisterTime, &m.WebsiteAdmin, &m.Username)
		managers = append(managers, m)
	}

	return &managers
}
