package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ku9nov/backup/configs"
	"github.com/minio/minio-go/v7"
	minioCredentials "github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type StorageClient interface {
	ListObjects(cfgValues configs.Config) (interface{}, error)
	RemoveFileFromS3(filename string, cfgValues configs.Config) error
	UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string) error
}

func CreateStorageClient(cfgValues configs.Config) StorageClient {

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
	case "azure":
		svcUrl := fmt.Sprintf("https://%s.blob.core.windows.net", cfgValues.Azure.StorageAccountName)
		azureAccountCred, err := azblob.NewSharedKeyCredential(cfgValues.Azure.StorageAccountName, cfgValues.Azure.StorageAccountKey)
		if err != nil {
			logrus.Errorf("error creating Azure Blob Storage shared key credential: %v", err)
			return nil
		}
		credential, err := azblob.NewClientWithSharedKeyCredential(svcUrl, azureAccountCred, nil)
		if err != nil {
			logrus.Errorf("error creating Azure Blob Storage shared key credential: %v", err)
			return nil
		}
		return &AzureStorageClient{Client: credential, BlobServiceURL: svcUrl}
	default:
		logrus.Errorf("Unknown storage driver: %s", cfgValues.Default.StorageProvider)
		return nil
	}
}

func CheckOldFilesInS3(cfgValues configs.Config, s3Client StorageClient) {
	objectKeys, err := s3Client.ListObjects(cfgValues)
	if err != nil {
		logrus.Fatal(err)
	}
	currentTime := time.Now()
	logrus.Debug("Retention candidates:")
	oldFiles := processFiles(objectKeys, cfgValues, currentTime)
	for _, oldFile := range oldFiles {
		if oldFile.Key != "" {
			logrus.Debugf("Key: %s, Last Modified: %v, Age: %v", oldFile.Key, oldFile.LastModified, oldFile.Age)
			if !cfgValues.Default.Retention.DryRun {
				s3Client.RemoveFileFromS3(oldFile.Key, cfgValues)
			}
		}
	}
}

func UploadToS3(cfgValues configs.Config, tarFilename []string, s3Client StorageClient) {
	today := time.Now()
	dailyPrefix := "daily/"
	switch today.Weekday() {
	case time.Monday:
		dailyPrefix = "weekly/"
	}
	if today.Day() == 1 {
		dailyPrefix = "monthly/"
	}
	for _, filename := range tarFilename {
		s3Client.UploadFileToS3(filename, cfgValues, dailyPrefix)
	}
}
