package utils

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ku9nov/backup/configs"
	"github.com/minio/minio-go/v7"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateStorageClient(cfgValues configs.Config) interface{} {

	switch cfgValues.Default.StorageProvider {
	case "minio":
		minioClient, err := minio.New(cfgValues.Minio.S3Endpoint, &minio.Options{
			Creds:  minioCredentials.NewStaticV4(cfgValues.Default.AccessKey, cfgValues.Default.SecretKey, ""),
			Secure: cfgValues.Minio.Secure,
		})
		if err != nil {
			logrus.Errorf("error setting up Minio client: %v", err)
			return nil
		}
		return &MinioStorageClient{Client: minioClient}

	case "aws":
		if cfgValues.Default.UseProfile.Enabled {
			logrus.Infof("The use profile is: %t, profile: %s", cfgValues.Default.UseProfile.Enabled, cfgValues.Default.UseProfile.Profile)
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile(cfgValues.Default.UseProfile.Profile),
			)
			if err != nil {
				logrus.Panicf("Failed loading config, %v", err)
			}
			return &AWSS3StorageClient{Client: s3.NewFromConfig(cfg)}
		} else {
			logrus.Infof("The use profile is: %t. Keys will be used as credentials.", cfgValues.Default.UseProfile.Enabled)
			creds := credentials.NewStaticCredentialsProvider(cfgValues.Default.AccessKey, cfgValues.Default.SecretKey, "")
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(cfgValues.Default.Region))
			if err != nil {
				logrus.Panicf("Failed authentication using keys as credentials: %v", err)
			}
			return &AWSS3StorageClient{Client: s3.NewFromConfig(cfg)}
		}
	default:
		logrus.Errorf("Unknown storage driver: %s", cfgValues.Default.StorageProvider)
		return nil
	}
}

func CheckOldFilesInS3(cfgValues configs.Config, s3Client interface{}) {
	var objectKeys interface{}
	var err error

	switch client := s3Client.(type) {
	case *MinioStorageClient:
		objectKeys, err = client.ListObjects(cfgValues)
		if err != nil {
			logrus.Fatal(err)
		}

	case *AWSS3StorageClient:
		objectKeys, err = client.ListObjects(cfgValues)
		if err != nil {
			logrus.Fatal(err)
		}
	default:
		logrus.Errorf("Unknown storage client type.")
	}
	currentTime := time.Now()
	logrus.Debug("Retention candidates:")
	switch keys := objectKeys.(type) {
	case []minio.ObjectInfo:
		processFiles(keys, cfgValues, currentTime)
	case *s3.ListObjectsV2Output:
		processFiles(keys, cfgValues, currentTime)
	default:
		logrus.Errorf("Unknown object type.")
	}
}

func UploadToS3(cfgValues configs.Config, tarFilename []string, s3Client interface{}) {
	today := time.Now()
	dailyPrefix := "daily/"
	switch today.Weekday() {
	case time.Monday:
		dailyPrefix = "weekly/"
	}
	if today.Day() == 1 {
		dailyPrefix = "monthly/"
	}
	switch client := s3Client.(type) {
	case *MinioStorageClient:
		for _, filename := range tarFilename {
			client.UploadFileToS3(filename, cfgValues, dailyPrefix)
		}
	case *AWSS3StorageClient:
		for _, filename := range tarFilename {
			client.UploadFileToS3(filename, cfgValues, dailyPrefix)
		}
	default:
		logrus.Errorf("Unknown storage client type.")
	}
}
