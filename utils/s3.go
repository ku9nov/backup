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
	ListObjects(cfgValues configs.Config, isExtraClient bool) (interface{}, error)
	RemoveFileFromS3(filename string, cfgValues configs.Config, isExtraClient bool) error
	UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string, isExtraClient bool) error
}

func SetStorageClient(cfgValues configs.Config) (StorageClient, StorageClient) {
	var mainClient, extraClient StorageClient

	switch cfgValues.Default.StorageProvider {
	case "minio", "aws", "azure":
		mainClient = CreateStorageClient(cfgValues, cfgValues.Default.StorageProvider, false)
	default:
		logrus.Errorf("Unknown default storage driver: %s", cfgValues.Default.StorageProvider)
		return nil, nil
	}

	if cfgValues.ExtraBackups.Enabled {
		switch cfgValues.ExtraBackups.StorageProvider {
		case "minio", "aws", "azure":
			extraClient = CreateStorageClient(cfgValues, cfgValues.ExtraBackups.StorageProvider, true)
		default:
			logrus.Errorf("Unknown extra backup storage driver: %s", cfgValues.ExtraBackups.StorageProvider)
			return nil, nil
		}
	}

	return mainClient, extraClient
}

func CreateStorageClient(cfgValues configs.Config, provider string, isExtraClient bool) StorageClient {
	var accessKey, secretKey, profile, region string
	var useProfile bool
	if !isExtraClient {
		accessKey = cfgValues.Default.AccessKey
		secretKey = cfgValues.Default.SecretKey
		useProfile = cfgValues.Default.UseProfile.Enabled
		profile = cfgValues.Default.UseProfile.Profile
		region = cfgValues.Default.Region
	} else {
		accessKey = cfgValues.ExtraBackups.AccessKey
		secretKey = cfgValues.ExtraBackups.SecretKey
		useProfile = cfgValues.ExtraBackups.UseProfile.Enabled
		profile = cfgValues.ExtraBackups.UseProfile.Profile
		region = cfgValues.ExtraBackups.Region
	}
	switch provider {
	case "minio":
		minioClient, err := minio.New(cfgValues.Minio.S3Endpoint, &minio.Options{
			Creds:  minioCredentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: cfgValues.Minio.Secure,
		})
		if err != nil {
			logrus.Errorf("error setting up Minio client: %v", err)
			return nil
		}
		return &MinioStorageClient{Client: minioClient}

	case "aws":
		if cfgValues.Default.UseProfile.Enabled {
			logrus.Infof("The use profile is: %t, profile: %s", useProfile, profile)
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile(profile),
			)
			if err != nil {
				logrus.Panicf("Failed loading config, %v", err)
			}
			return &AWSS3StorageClient{Client: s3.NewFromConfig(cfg)}
		} else {
			logrus.Infof("The use profile is: %t. Keys will be used as credentials.", useProfile)
			creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(region))
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
		logrus.Errorf("Unknown storage driver: %s", provider)
		return nil
	}
}

func CheckOldFilesInS3(cfgValues configs.Config, s3Client StorageClient, isExtraClient bool) {
	objectKeys, err := s3Client.ListObjects(cfgValues, isExtraClient)
	if err != nil {
		logrus.Fatal(err)
	}
	currentTime := time.Now()
	logrus.Debug("Retention candidates:")
	oldFiles := processFiles(objectKeys, cfgValues, currentTime, isExtraClient)
	for _, oldFile := range oldFiles {
		if oldFile.Key != "" {
			logrus.Debugf("Key: %s, Last Modified: %v, Age: %v", oldFile.Key, oldFile.LastModified, oldFile.Age)
			if !cfgValues.Default.Retention.DryRun {
				s3Client.RemoveFileFromS3(oldFile.Key, cfgValues, false)
			}
			if isExtraClient && !cfgValues.ExtraBackups.Retention.DryRun {
				s3Client.RemoveFileFromS3(oldFile.Key, cfgValues, true)
			}
		}
	}
}

func UploadToS3(cfgValues configs.Config, tarFilename []string, s3Client, extraS3Client StorageClient) {
	dailyPrefix := setDailyPrefix()
	for _, filename := range tarFilename {
		s3Client.UploadFileToS3(filename, cfgValues, dailyPrefix, false)
		if cfgValues.ExtraBackups.Enabled {
			extraS3Client.UploadFileToS3(filename, cfgValues, fmt.Sprintf(cfgValues.Default.Bucket+"/"), true)
		}
	}
}
