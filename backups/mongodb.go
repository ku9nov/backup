package backups

import (
	"github.com/sirupsen/logrus"
)

func CreateMongoBackup(host string, port int, auth bool, username, password string) {
	logrus.Info("MongoDB is enabled, executing MongoDB example code...")
	if auth {
		logrus.Debugf("Authentification is enabled.\nConnecting to MongoDB at %s:%d with username '%s' and password '%s'\n", host, port, username, password)
	} else {
		logrus.Debug("Authentification is disabled.\nConnecting to MongoDB without auth.")
	}
}
