package backups

import "github.com/sirupsen/logrus"

func CreateMySQLBackup(host string, port int, auth bool, username, password string) {
	logrus.Info("MySQL is enabled, executing MySQL example code...")
	if auth {
		logrus.Debugf("Authentification is enabled.\nConnecting to MySQL at %s:%d with username '%s' and password '%s'\n", host, port, username, password)
	} else {
		logrus.Debug("Authentification is disabled.\nConnecting to MySQL without auth.")
	}
}
