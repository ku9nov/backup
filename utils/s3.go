package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
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
		return minioClient

	case "aws":
		if cfgValues.Default.UseProfile.Enabled {
			logrus.Infof("The use profile is: %t, profile: %s", cfgValues.Default.UseProfile.Enabled, cfgValues.Default.UseProfile.Profile)
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithSharedConfigProfile(cfgValues.Default.UseProfile.Profile),
			)
			if err != nil {
				panic(fmt.Sprintf("Failed loading config, %v", err))
			}
			return s3.NewFromConfig(cfg)
		} else {
			logrus.Infof("The use profile is: %t. Keys will be used as credentials.", cfgValues.Default.UseProfile.Enabled)
			creds := credentials.NewStaticCredentialsProvider(cfgValues.Default.AccessKey, cfgValues.Default.SecretKey, "")
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(cfgValues.Default.Region))
			if err != nil {
				panic(fmt.Sprintf("Failed authentication using keys as credentials: %v", err))
			}
			return s3.NewFromConfig(cfg)
		}
	default:
		logrus.Errorf("Unknown storage driver: %s", cfgValues.Default.StorageProvider)
		return nil
	}
}

func CheckOldFilesInS3(cfgValues configs.Config, s3Client interface{}) {
	var output *s3.ListObjectsV2Output
	var err error
	switch client := s3Client.(type) {
	case *minio.Client:
		logrus.Info("Minio catched.")

	case *s3.Client:
		output, err = client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(cfgValues.Default.Bucket),
		})
		if err != nil {
			logrus.Fatal(err)
		}
	default:
		logrus.Errorf("Unknown storage client type.")
	}

	currentTime := time.Now()
	logrus.Debug("Retention candidates:")
	for _, object := range output.Contents {
		// Calculate the age of the object
		age := currentTime.Sub(*object.LastModified)
		isFolder := strings.HasSuffix(aws.ToString(object.Key), "/")
		// Check if the object is in the daily folder and older than 1 day
		if strings.HasPrefix(aws.ToString(object.Key), "daily/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodDaily)*24*time.Hour {
			logrus.Debugf("Key: %s, Size: %d, Last Modified: %v, Age: %v", aws.ToString(object.Key), object.Size, object.LastModified, age)
			// if err := deleteObject(client, cfgValues.Default.Bucket, aws.ToString(object.Key)); err != nil {
			//     logrus.Errorf("Delete failed: %v", err)
			// } else {
			//     logrus.Infof("Deleted object: %s", aws.ToString(object.Key))
			// }
		}

		// Check if the object is in the weekly folder and older than 1 month
		if strings.HasPrefix(aws.ToString(object.Key), "weekly/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodWeekly)*24*time.Hour*7 {
			logrus.Debugf("Key: %s, Size: %d, Last Modified: %v, Age: %v", aws.ToString(object.Key), object.Size, object.LastModified, age)
			// if err := deleteObject(client, cfgValues.Default.Bucket, aws.ToString(object.Key)); err != nil {
			//     logrus.Errorf("Delete failed: %v", err)
			// } else {
			//     logrus.Infof("Deleted object: %s", aws.ToString(object.Key))
			// }
		}

		// Check if the object is in the monthly folder and older than 6 months
		if strings.HasPrefix(aws.ToString(object.Key), "monthly/") && !isFolder && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodMonthly)*24*time.Hour*30 {
			logrus.Debugf("Key: %s, Size: %d, Last Modified: %v, Age: %v", aws.ToString(object.Key), object.Size, object.LastModified, age)
			// if err := deleteObject(client, cfgValues.Default.Bucket, aws.ToString(object.Key)); err != nil {
			//     logrus.Errorf("Delete failed: %v", err)
			// } else {
			//     logrus.Infof("Deleted object: %s", aws.ToString(object.Key))
			// }
		}
	}
}

// func deleteObject(client *s3.Client, bucket string, key string) error {
// 	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		return logrus.Errorf("failed to delete object %s: %w", key, err)
// 	}
// 	return nil
// }

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
	case *minio.Client:
		logrus.Info("Minio catched.")

	case *s3.Client:
		for _, filename := range tarFilename {
			file, err := os.Open(filename)
			if err != nil {
				logrus.Fatalf("Failed to open file %s: %v", filename, err)
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				logrus.Fatalf("Failed to get file info for %s: %v", filename, err)
			}

			_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket:        aws.String(cfgValues.Default.Bucket),
				Key:           aws.String(dailyPrefix + filepath.Base(filename)), // Adjust key to include daily, weekly, or monthly prefix
				Body:          file,
				ContentType:   aws.String("application/octet-stream"),
				ContentLength: aws.Int64(fileInfo.Size()),
			})
			if err != nil {
				logrus.Fatalf("Failed to upload file %s to S3: %v", filename, err)
			}

			logrus.Infof("%s was successfully uploaded to S3.\n", filename)
		}

	default:
		logrus.Errorf("Unknown storage client type.")
	}

}
