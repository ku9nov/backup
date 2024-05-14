package backups

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateMySQLBackup(host, port, username, password, tool, backupDir, currentDate string, databases []string, auth bool) {
	var files []string
	success := false
	utils.CheckToolIsExist(tool)

	logrus.Info("MySQL is enabled, executing MySQL backup code...")
	for _, db := range databases {
		cmdArgs := []string{tool, "-h", host, "-P", port}
		if auth {
			logrus.Info("MySQL backup with authentication is required.")
			cmdArgs = append(cmdArgs, "-u", username, fmt.Sprintf("-p%s", password))
		}

		cmdArgs = append(cmdArgs, db)

		filePath := fmt.Sprintf("%s/%s.sql", backupDir, db)
		file, err := os.Create(filePath)
		if err != nil {
			logrus.Errorf("Error creating file for database %s: %v\n", db, err)
			success = false
			break
		}
		defer file.Close()

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = file

		err = cmd.Run()
		if err != nil {
			logrus.Debugf("%s %s", cmdArgs[0], strings.Join(cmdArgs[1:], " "))
			logrus.Errorf("Error backing up database %s: %v\n", db, err)
			success = false
			break
		}

		logrus.Infof("Backup for database %s created successfully.\n", db)
		files = append(files, filePath)
		success = true
	}
	if success {
		tarFilename := utils.TarFiles("mysql", currentDate, backupDir, files)
		files = append(files, tarFilename...)
		utils.CleanupFilesAndTar(files)
	}

}
