package utils

import (
	"context"
	"fmt"

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
		)
		if err != nil {
			panic(fmt.Sprintf("Failed authentication using keys as credentials: %v", err))
		}
	}
	return awsCfg
}

func ConnectToS3(awsCfg aws.Config, cfgValues configs.Config) {
	client := s3.NewFromConfig(awsCfg)

	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(cfgValues.Default.Bucket),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debug("s3 results:")
	for _, object := range output.Contents {
		logrus.Debugf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}
}
