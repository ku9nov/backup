package backups

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/configs"
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreatePostgreSQLBackup(cfgValues configs.Config, currentDate string, s3Cfg utils.StorageClient) bool {
	var files []string
	success := false
	success = utils.CheckToolIsExist(cfgValues.PostgreSQL.DumpTool)

	logrus.Info("PostgreSQL is enabled, executing PostgreSQL backup code...")
	for _, db := range cfgValues.PostgreSQL.Databases {
		cmdArgs := []string{cfgValues.PostgreSQL.DumpTool, "-h", cfgValues.PostgreSQL.Host, "-p", cfgValues.PostgreSQL.Port}
		if cfgValues.PostgreSQL.Auth.Enabled {
			logrus.Info("PostgreSQL backup with authentication is required.")
			cmdArgs = append(cmdArgs, "-U", cfgValues.PostgreSQL.Auth.Username)
		}

		cmdArgs = append(cmdArgs, "-d", db)

		filePath := fmt.Sprintf("%s/%s.sql", cfgValues.Default.BackupDir, db)
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Error creating file for database %s: %v\n", db, err)
			success = false
			break
		}
		defer file.Close()

		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Stdout = file
		if cfgValues.PostgreSQL.Auth.Enabled {
			cmd.Env = append(os.Environ(), "PGPASSWORD="+cfgValues.PostgreSQL.Auth.Password)
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
		tarFilename := utils.TarFiles("postgres", currentDate, cfgValues.Default.BackupDir, files)
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
