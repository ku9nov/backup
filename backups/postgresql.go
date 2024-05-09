package backups

import (
	"github.com/sirupsen/logrus"
)

func CreatePostgreSQLBackup(host string, port int, auth bool, username, password string) {
	logrus.Info("PostgreSQL is enabled, executing PostgreSQL example code...")
	if auth {
		logrus.Debugf("Authentification is enabled.\nConnecting to PostgreSQL at %s:%d with username '%s' and password '%s'\n", host, port, username, password)
	} else {
		logrus.Debug("Authentification is disabled.\nConnecting to PostgreSQL without auth.")
	}
}
