package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/ku9nov/backup/configs"
	"github.com/sirupsen/logrus"
)

type AzureStorageClient struct {
	Client         *azblob.Client
	BlobServiceURL string
}
type FileMetadata struct {
	Name         string
	Size         int64
	LastModified time.Time
}

func (c *AzureStorageClient) ListObjects(cfgValues configs.Config, isExtraClient bool) (interface{}, error) {
	bucketName := getBucketName(cfgValues, isExtraClient)
	pager := c.Client.NewListBlobsFlatPager(bucketName, nil)
	var files []FileMetadata
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, v := range resp.Segment.BlobItems {
			files = append(files, FileMetadata{
				Name:         *v.Name,
				Size:         *v.Properties.ContentLength,
				LastModified: *v.Properties.LastModified,
			})
		}
	}
	return files, nil
}
func (c *AzureStorageClient) UploadFileToS3(filename string, cfgValues configs.Config, dailyPrefix string, isExtraClient bool) error {
	bucketName := getBucketName(cfgValues, isExtraClient)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()
	_, err = c.Client.UploadFile(context.TODO(), bucketName, dailyPrefix+filepath.Base(filename), file,
		&azblob.UploadFileOptions{
			BlockSize:   int64(1024),
			Concurrency: uint16(3),
			Progress:    func(bytesTransferred int64) {},
		})
	if err != nil {
		return fmt.Errorf("failed to upload file %s to Azure Blob Storage: %v", filename, err)
	}

	logrus.Infof("%s was successfully uploaded to Azure Blob Storage.\n", filename)
	return nil
}

func (c *AzureStorageClient) RemoveFileFromS3(filename string, cfgValues configs.Config, isExtraClient bool) error {
	bucketName := getBucketName(cfgValues, isExtraClient)
	_, err := c.Client.DeleteBlob(context.TODO(), bucketName, filename, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob from Azure storage: %v", err)
	}

	logrus.Infof("%s was successfully removed from Azure Blob Storage.\n", filename)
	return nil
}
