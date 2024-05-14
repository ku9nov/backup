package backups

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateMongoBackup(host, port, username, password, authDatabase, tool, backupDir, currentDate string, databases []string, auth bool) {
	var files []string
	success := false
	logrus.Info("MongoDB is enabled, executing MongoDB backup code...")

	for _, db := range databases {
		var uri string
		if auth {
			logrus.Info("MongoDB backup with authentication is required.")
			uri = fmt.Sprintf(`"mongodb://%s:%s@%s:%s/?authSource=%s"`, username, password, host, port, authDatabase)
		} else {
			uri = fmt.Sprintf(`"mongodb://%s:%s"`, host, port)
		}

		cmdArgs := []string{tool, "--out", fmt.Sprintf("%s/%s", backupDir, db), "--uri", uri, "--db", db}

		output, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
		if err != nil {
			logrus.Debugf("%s %s", cmdArgs[0], strings.Join(cmdArgs[1:], " "))
			logrus.Errorf("Error backing up database %s: %v\n", db, err)
			logrus.Errorf("Stderr output for database %s: %s\n", db, output)
			success = false
			break
		}

		logrus.Infof("Backup for database %s created successfully.\n", db)
		files = append(files, fmt.Sprintf("%s/%s", backupDir, db))
		success = true
	}
	if success {
		tarFilename := utils.TarFiles("mongo", currentDate, backupDir, files)
		files = append(files, tarFilename...)
		utils.CleanupFilesAndTar(files)
	}

}
