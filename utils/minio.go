package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ku9nov/backup/configs"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type MinioStorageClient struct {
	Client *minio.Client
}

func (c *MinioStorageClient) ListObjects(cfgValues configs.Config) (interface{}, error) {
	ctx := context.TODO()
	objectsCh := c.Client.ListObjects(ctx, cfgValues.Default.Bucket, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	var objects []minio.ObjectInfo
	for object := range objectsCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)

	}
	return objects, nil
}

func (c *MinioStorageClient) UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", filename, err)
	}

	// Upload file to MinIO
	_, err = c.Client.PutObject(context.TODO(), cfgValues.Default.Bucket, dailyPrefix+filepath.Base(filename), file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to MinIO: %v", filename, err)
	}

	logrus.Infof("%s was successfully uploaded to MinIO.\n", filename)
	return nil
}

func (c *MinioStorageClient) RemoveFileFromS3(filename string, cfgValues configs.Config) error {
	// Remove object from MinIO
	err := c.Client.RemoveObject(context.TODO(), cfgValues.Default.Bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove file %s from MinIO: %v", filename, err)
	}

	logrus.Infof("%s was successfully removed from MinIO.\n", filename)
	return nil
}
