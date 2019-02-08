package storage

import (
	"bytes"
	"crypto"
	"errors"
	"github.com/Sirupsen/logrus"
)

type ManagerRole int

const (
	ROLE_READ          ManagerRole = 1
	ROLE_EDIT          ManagerRole = 2
	ROLE_MANAGE_ACCESS ManagerRole = 3
)

type Manager struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	WebsiteAdmin bool   `json:"website_admin"`
}

func (database *Database) ValidateCredentials(username []byte, password []byte) (*Manager, error) {

	db := database.getDB()

	row := db.QueryRow(`SELECT id, password, website_admin FROM manager WHERE username=$1`,
		username)

	manager := &Manager{}
	var passwordHash []byte
	err := row.Scan(&manager.Id, &passwordHash, &manager.WebsiteAdmin)
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

	row := db.QueryRow(`INSERT INTO manager (username, password, website_admin) 
			VALUES ($1,$2,$3) RETURNING ID`,
		manager.Username, hashedPassword, manager.WebsiteAdmin)

	err := row.Scan(&manager.Id)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": manager,
		}).Warning("Database.SaveManager INSERT error")

		return err
	}

	logrus.WithFields(logrus.Fields{
		"manager": manager,
	}).Info("Database.SaveManager INSERT")

	return nil
}

func (database *Database) UpdateManager(manager *Manager) {

	db := database.getDB()

	res, err := db.Exec(`UPDATE manager SET website_admin=$1 WHERE id=$2`,
		manager.WebsiteAdmin, manager.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithError(err).WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"manager":      manager,
	}).Warning("Database.UpdateManager UPDATE")
}

func (database *Database) UpdateManagerPassword(manager *Manager, newPassword []byte) {

	hash := crypto.SHA512.New()
	hash.Write([]byte(manager.Username))
	hash.Write(newPassword)
	hashedPassword := hash.Sum(nil)

	db := database.getDB()

	res, err := db.Exec(`UPDATE manager SET password=$1 WHERE id=$2`,
		hashedPassword, manager.Id)
	handleErr(err)

	rowsAffected, _ := res.RowsAffected()

	logrus.WithError(err).WithFields(logrus.Fields{
		"rowsAffected": rowsAffected,
		"id":           manager.Id,
	}).Warning("Database.UpdateManagerPassword UPDATE")
}
