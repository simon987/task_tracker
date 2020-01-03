package storage

import "github.com/sirupsen/logrus"

func handleErr(err error) {
	if err != nil {
		logrus.WithError(err).Error("Error during database operation!")
	}
}
