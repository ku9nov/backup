package backups

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/configs"
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateMySQLBackup(cfgValues configs.Config, currentDate string, s3Cfg utils.StorageClient) bool {
	var files []string
	success := false
	success = utils.CheckToolIsExist(cfgValues.MySQL.DumpTool)

	logrus.Info("MySQL is enabled, executing MySQL backup code...")
	for _, db := range cfgValues.MySQL.Databases {
		cmdArgs := []string{cfgValues.MySQL.DumpTool, "-h", cfgValues.MySQL.Host, "-P", cfgValues.MySQL.Port}
		if cfgValues.MySQL.Auth.Enabled {
			logrus.Info("MySQL backup with authentication is required.")
			cmdArgs = append(cmdArgs, "-u", cfgValues.MySQL.Auth.Username, fmt.Sprintf("-p%s", cfgValues.MySQL.Auth.Password))
		}

		cmdArgs = append(cmdArgs, db)

		filePath := fmt.Sprintf("%s/%s.sql", cfgValues.Default.BackupDir, db)
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
		tarFilename := utils.TarFiles("mysql", currentDate, cfgValues.Default.BackupDir, files)
		if len(tarFilename) == 0 {
			logrus.Error("Error creating tar file.")
			success = false
		} else {
			utils.UploadToS3(cfgValues, tarFilename, s3Cfg)
			files = append(files, tarFilename...)
			utils.CleanupFilesAndTar(files)
		}
	}
	return success
}
