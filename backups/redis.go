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

func CreateRedisBackup(cfgValues configs.Config, currentDate string, s3Cfg, extraS3Cfg utils.StorageClient) bool {
	var files []string
	success := false

	success = utils.CheckToolIsExist(cfgValues.Redis.RedisCliTool)

	logrus.Info("Redis backup process initiated...")

	cmdArgs := []string{cfgValues.Redis.RedisCliTool, "-h", cfgValues.Redis.Host, "-p", cfgValues.Redis.Port}

	if cfgValues.Redis.Auth.Enabled {
		logrus.Info("Redis backup with authentication is required.")
		cmdArgs = append(cmdArgs, "-a", cfgValues.Redis.Auth.Password)
	}

	filePath := fmt.Sprintf("%s/redis.rdb", cfgValues.Default.BackupDir)
	cmdArgs = append(cmdArgs, "--rdb", filePath)
	file, err := os.Create(filePath)
	if err != nil {
		logrus.Errorf("Error creating file for Redis backup: %v\n", err)
		success = false
		return success
	}
	defer file.Close()

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = file

	err = cmd.Run()
	if err != nil {
		logrus.Debugf("%s %s", cmdArgs[0], strings.Join(cmdArgs[1:], " "))
		logrus.Errorf("Error backing up Redis: %v\n", err)
		success = false
		return success
	}

	logrus.Info("Backup for Redis created successfully.\n")
	files = append(files, filePath)
	if success {
		tarFilename := utils.TarFiles("redis", currentDate, cfgValues.Default.BackupDir, []string{filePath})
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
