package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ku9nov/backup/configs"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func CheckToolIsExist(tool string) bool {
	_, err := exec.LookPath(tool)
	if err != nil {
		logrus.Errorf("%s command not found. Please ensure %s is installed and properly configured.", tool, tool)
		return false
	}
	return true
}
func TarFiles(backupSource, currentDate, backupDir string, files []string) []string {
	tarGzFilename := filepath.Join(backupDir, backupSource+"-"+currentDate+".tgz")

	tarGzFile, err := os.Create(tarGzFilename)
	if err != nil {
		logrus.Error("Error creating tar.gz file:", err)
		return nil
	}
	defer tarGzFile.Close()

	gzipWriter := gzip.NewWriter(tarGzFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, file := range files {
		if err := AddToTar(tarWriter, file); err != nil {
			logrus.Error("Error adding file/folder to tar.gz:", err)
			return nil
		}
	}

	logrus.Info("Files/folders archived to:", tarGzFilename)
	return []string{tarGzFilename}
}

func AddToTar(tarWriter *tar.Writer, path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() {
				relPath, err := filepath.Rel(filepath.Dir(path), filePath)
				if err != nil {
					return err
				}
				return addFileToTar(tarWriter, filePath, relPath, fileInfo)
			}
			return nil
		})
	}

	relPath := filepath.Base(path)
	return addFileToTar(tarWriter, path, relPath, info)
}

func addFileToTar(tarWriter *tar.Writer, path string, relPath string, info os.FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(relPath)

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}

func CleanupFilesAndTar(files []string) {
	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			logrus.Errorf("Error removing file %s: %v", file, err)
		} else {
			logrus.Infof("%s was removed successfully.", file)
		}
	}
}

func processFiles(object interface{}, cfgValues configs.Config, currentTime time.Time) {
	var key string
	var lastModified time.Time

	switch obj := object.(type) {
	case *s3.ListObjectsV2Output:
		for _, object := range obj.Contents {
			key = *object.Key
			lastModified = *object.LastModified
			processFileData(key, lastModified, cfgValues, currentTime)
		}
	case []minio.ObjectInfo:
		for _, object := range obj {
			key = object.Key
			lastModified = object.LastModified
			processFileData(key, lastModified, cfgValues, currentTime)
		}
	default:
		logrus.Info("Unknown object type.")
	}
}

func processFileData(key string, lastModified time.Time, cfgValues configs.Config, currentTime time.Time) {
	age := currentTime.Sub(lastModified)
	isFolder := strings.HasSuffix(key, "/")

	if strings.HasPrefix(key, "daily/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodDaily)*24*time.Hour {
		logrus.Debugf("Key: %s, Last Modified: %v, Age: %v", key, lastModified, age)

	}

	if strings.HasPrefix(key, "weekly/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodWeekly)*24*time.Hour*7 {
		logrus.Debugf("Key: %s, Last Modified: %v, Age: %v", key, lastModified, age)

	}

	if strings.HasPrefix(key, "monthly/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodMonthly)*24*time.Hour*30 {
		logrus.Debugf("Key: %s, Last Modified: %v, Age: %v", key, lastModified, age)

	}
}
