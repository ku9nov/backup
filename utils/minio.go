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

func (c *MinioStorageClient) ListObjects(cfgValues configs.Config, isExtraClient bool) (interface{}, error) {
	bucketName := getBucketName(cfgValues, isExtraClient)
	ctx := context.TODO()
	objectsCh := c.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
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

func (c *MinioStorageClient) UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string, isExtraClient bool) error {
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
	bucketName := getBucketName(cfgValues, isExtraClient)
	// Upload file to MinIO
	_, err = c.Client.PutObject(context.TODO(), bucketName, dailyPrefix+filepath.Base(filename), file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to MinIO/Spaces: %v", filename, err)
	}

	logrus.Infof("%s was successfully uploaded to MinIO/Spaces.\n", filename)
	return nil
}

func (c *MinioStorageClient) RemoveFileFromS3(filename string, cfgValues configs.Config, isExtraClient bool) error {
	bucketName := getBucketName(cfgValues, isExtraClient)
	// Remove object from MinIO
	err := c.Client.RemoveObject(context.TODO(), bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove file %s from MinIO/Spaces: %v", filename, err)
	}

	logrus.Infof("%s was successfully removed from MinIO/Spaces.\n", filename)
	return nil
}
