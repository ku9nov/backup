package backups

import (
	"github.com/ku9nov/backup/configs"
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateAdditionalFilesBackup(cfgValues configs.Config, currentDate string, s3Cfg interface{}) bool {
	success := true
	logrus.Info("Additional files enabled, processing files:")
	for _, file := range cfgValues.Additional.Files {
		logrus.Info(file)
	}
	tarFilename := utils.TarFiles("etc", currentDate, cfgValues.Default.BackupDir, cfgValues.Additional.Files)
	if len(tarFilename) == 0 {
		logrus.Error("Error creating tar file.")
		success = false
	} else {
		utils.UploadToS3(cfgValues, tarFilename, s3Cfg)
		utils.CleanupFilesAndTar(tarFilename)
	}
	return success
}
