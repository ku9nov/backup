package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ku9nov/backup/configs"
)

func AWSAuth(cfgValues configs.Config) (awsCfg aws.Config) {
	var err error

	if cfgValues.Default.UseProfile.Enabled {
		// Load config using shared profile
		logrus.Infof("The use profile is: %t, profile: %s", cfgValues.Default.UseProfile.Enabled, cfgValues.Default.UseProfile.Profile)
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithSharedConfigProfile(cfgValues.Default.UseProfile.Profile),
		)
		if err != nil {
			panic(fmt.Sprintf("failed loading config, %v", err))
		}
	} else {
		// Load config using static credentials
		logrus.Infof("The use profile is: %t. Keys will be used as credentials.", cfgValues.Default.UseProfile.Enabled)
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfgValues.Default.AccessKey, cfgValues.Default.SecretKey, "")),
			config.WithRegion(cfgValues.Default.Region),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed authentication using keys as credentials: %v", err))
		}
	}
	return awsCfg
}

func CheckOldFilesInS3(awsCfg aws.Config, cfgValues configs.Config) {
	client := s3.NewFromConfig(awsCfg)

	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(cfgValues.Default.Bucket),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	currentTime := time.Now()
	logrus.Debug("Retention candidates:")
	for _, object := range output.Contents {
		// Calculate the age of the object
		age := currentTime.Sub(*object.LastModified)

		// Check if the object is in the daily folder and older than 1 day
		if strings.HasPrefix(aws.ToString(object.Key), "daily/") && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodDaily)*24*time.Hour {
			logrus.Debugf("Key: %s, Size: %d, Last Modified: %v, Age: %v", aws.ToString(object.Key), object.Size, object.LastModified, age)
			// if err := deleteObject(client, cfgValues.Default.Bucket, aws.ToString(object.Key)); err != nil {
			//     logrus.Errorf("Delete failed: %v", err)
			// } else {
			//     logrus.Infof("Deleted object: %s", aws.ToString(object.Key))
			// }
		}

		// Check if the object is in the weekly folder and older than 1 month
		if strings.HasPrefix(aws.ToString(object.Key), "weekly/") && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodWeekly)*24*time.Hour*7 {
			logrus.Debugf("Key: %s, Size: %d, Last Modified: %v, Age: %v", aws.ToString(object.Key), object.Size, object.LastModified, age)
			// if err := deleteObject(client, cfgValues.Default.Bucket, aws.ToString(object.Key)); err != nil {
			//     logrus.Errorf("Delete failed: %v", err)
			// } else {
			//     logrus.Infof("Deleted object: %s", aws.ToString(object.Key))
			// }
		}

		// Check if the object is in the monthly folder and older than 6 months
		if strings.HasPrefix(aws.ToString(object.Key), "monthly/") && age > time.Duration(cfgValues.Default.Retention.RetentionPeriodMonthly)*24*time.Hour*30 {
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
