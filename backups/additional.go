package backups

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func CreateAdditionalFilesBackup(files []string, currentDate string) {
	logrus.Info("Additional files enabled, processing files:")
	for _, file := range files {
		logrus.Info(file)
	}

	tarGzFilename := "etc-" + currentDate + ".tgz"

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
		if err := addFileToTar(tarWriter, file); err != nil {
			logrus.Error("Error adding file to tar.gz:", err)
		}
	}

	logrus.Info("Files archived to:", tarGzFilename)
}

func addFileToTar(tarWriter *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    filepath.Base(filename),
		Size:    fileInfo.Size(),
		Mode:    int64(fileInfo.Mode()),
		ModTime: fileInfo.ModTime(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return err
	}

	return nil
}
