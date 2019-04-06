package storage

import "github.com/sirupsen/logrus"

func handleErr(err error) {
	if err != nil {
		logrus.WithError(err).Fatal("Error during database operation!")
	}
}
