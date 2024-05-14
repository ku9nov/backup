package backups

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreatePostgreSQLBackup(host, port, username, password, tool, backupDir, currentDate string, databases []string, auth bool) {
	var files []string
	success := false
	utils.CheckToolIsExist(tool)

	logrus.Info("PostgreSQL is enabled, executing PostgreSQL backup code...")
	for _, db := range databases {
		cmdArgs := []string{tool, "-h", host, "-p", port}
		if auth {
			logrus.Info("PostgreSQL backup with authentication is required.")
			cmdArgs = append(cmdArgs, "-U", username)
		}

		cmdArgs = append(cmdArgs, "-d", db)

		filePath := fmt.Sprintf("%s/%s.sql", backupDir, db)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file for database %s: %v\n", db, err)
			success = false
			break
		}
		defer file.Close()

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = file
		if auth {
			cmd.Env = append(os.Environ(), "PGPASSWORD="+password)
		}
		err = cmd.Run()
		if err != nil {
			logrus.Errorf("%s %s", cmdArgs[0], strings.Join(cmdArgs[1:], " "))
			logrus.Errorf("Error backing up database %s: %v\n", db, err)
			success = false
			break
		}

		logrus.Infof("Backup for database %s created successfully.\n", db)
		files = append(files, filePath)
		success = true
	}

	if success {
		tarFilename := utils.TarFiles("postgres", currentDate, backupDir, files)
		files = append(files, tarFilename...)
		utils.CleanupFilesAndTar(files)
	}
}
