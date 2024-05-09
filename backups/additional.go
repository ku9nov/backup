package backups

import (
	"github.com/sirupsen/logrus"
)

func CreateAdditionalFilesBackup(files []string) {
	logrus.Info("Additional files enabled, processing files:")
	for _, file := range files {
		logrus.Info(file)
	}
}
