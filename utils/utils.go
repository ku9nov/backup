package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func CheckToolIsExist(tool string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		logrus.Errorf("%s command not found. Please ensure %s is installed and properly configured.", tool, tool)
		return
	}
}
func TarFiles(backupSource, currentDate, backupDir string, files []string) {
	tarGzFilename := filepath.Join(backupDir, backupSource+"-"+currentDate+".tgz")

	tarGzFile, err := os.Create(tarGzFilename)
	if err != nil {
		logrus.Error("Error creating tar.gz file:", err)
		return
	}
	defer tarGzFile.Close()

	gzipWriter := gzip.NewWriter(tarGzFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, file := range files {
		if err := AddToTar(tarWriter, file); err != nil {
			logrus.Error("Error adding file/folder to tar.gz:", err)
		}
	}

	logrus.Info("Files/folders archived to:", tarGzFilename)
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
