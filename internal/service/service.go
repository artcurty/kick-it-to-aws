package service

import (
	"github.com/artcurty/kick-it-to-aws/internal/cloud"
)

type AWSService interface {
	UploadDirectoryToS3Batch(awsConfig *cloud.Config, sourceDir, targetPath, bucket string) error
	InvalidateCloudFront(awsConfig *cloud.Config, distributionID string) error
}
