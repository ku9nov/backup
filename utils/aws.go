package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ku9nov/backup/configs"
	"github.com/sirupsen/logrus"
)

type AWSS3StorageClient struct {
	Client *s3.Client
}

func (c *AWSS3StorageClient) ListObjects(cfgValues configs.Config, isExtraClient bool) (interface{}, error) {
	bucketName := getBucketName(cfgValues, isExtraClient)
	output, err := c.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (c *AWSS3StorageClient) UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string, isExtraClient bool) error {
	bucketName := getBucketName(cfgValues, isExtraClient)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %v", filename, err)
	}
	_, err = c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(dailyPrefix + filepath.Base(filename)), // Adjust key to include daily, weekly, or monthly prefix
		Body:          file,
		ContentType:   aws.String("application/octet-stream"),
		ContentLength: aws.Int64(fileInfo.Size()),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to S3: %v", filename, err)
	}

	logrus.Infof("%s was successfully uploaded to S3.\n", filename)
	return nil
}
func (c *AWSS3StorageClient) RemoveFileFromS3(filename string, cfgValues configs.Config, isExtraClient bool) error {
	bucketName := getBucketName(cfgValues, isExtraClient)
	_, err := c.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{
				{Key: aws.String(filename)},
			},
		},
	})
	if err != nil {
		logrus.Errorf("Couldn't delete objects from bucket %v. Here's why: %v\n", cfgValues.Default.Bucket, err)
	}
	logrus.Infof("%s was successfully removed from AWS s3.\n", filename)
	return nil
}
