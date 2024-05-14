package backups

import (
	"github.com/ku9nov/backup/utils"
	"github.com/sirupsen/logrus"
)

func CreateAdditionalFilesBackup(files []string, currentDate, backupDir string) {
	logrus.Info("Additional files enabled, processing files:")
	for _, file := range files {
		logrus.Info(file)
	}
	tarFilename := utils.TarFiles("etc", currentDate, backupDir, files)
	utils.CleanupFilesAndTar(tarFilename)
}
