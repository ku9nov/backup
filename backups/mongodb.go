package backups

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ku9nov/backup/configs"
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateMongoBackup(cfgValues configs.Config, currentDate string, s3Cfg, extraS3Cfg utils.StorageClient) bool {
	var files []string
	success := false
	logrus.Info("MongoDB is enabled, executing MongoDB backup code...")
	success = utils.CheckToolIsExist(cfgValues.Mongo.DumpTool)
	for _, db := range cfgValues.Mongo.Databases {
		var uri string
		if cfgValues.Mongo.Auth.Enabled {
			logrus.Info("MongoDB backup with authentication is required.")
			uri = fmt.Sprintf(`"mongodb://%s:%s@%s:%s/?authSource=%s"`, cfgValues.Mongo.Auth.Username, cfgValues.Mongo.Auth.Password, cfgValues.Mongo.Host, cfgValues.Mongo.Port, cfgValues.Mongo.Auth.AuthDatabase)
		} else {
			uri = fmt.Sprintf(`"mongodb://%s:%s"`, cfgValues.Mongo.Host, cfgValues.Mongo.Port)
		}

		cmdArgs := []string{cfgValues.Mongo.DumpTool, "--out", fmt.Sprintf("%s/%s", cfgValues.Default.BackupDir, db), "--uri", uri, "--db", db}

		output, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
		if err != nil {
			logrus.Debugf("%s %s", cmdArgs[0], strings.Join(cmdArgs[1:], " "))
			logrus.Errorf("Error backing up database %s: %v\n", db, err)
			logrus.Errorf("Stderr output for database %s: %s\n", db, output)
			success = false
			break
		}

		logrus.Infof("Backup for database %s created successfully.\n", db)
		files = append(files, fmt.Sprintf("%s/%s", cfgValues.Default.BackupDir, db))
		success = true
	}
	if success {
		tarFilename := utils.TarFiles("mongo", currentDate, cfgValues.Default.BackupDir, files)
		if len(tarFilename) == 0 {
			logrus.Error("Error creating tar file.")
			success = false
		} else {
			utils.UploadToS3(cfgValues, tarFilename, s3Cfg, extraS3Cfg)
			files = append(files, tarFilename...)
			utils.CleanupFilesAndTar(files)
		}
	}
	return success
}
