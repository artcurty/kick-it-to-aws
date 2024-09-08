package mock

import "github.com/artcurty/kick-it-to-aws/internal/cloud"

type AWSServiceMock struct {
	UploadDirectoryToS3BatchFunc func(awsConfig cloud.Config, sourceDir, targetPath, bucket string) error
	InvalidateCloudFrontFunc     func(awsConfig cloud.Config, distributionID string) error
}

func (m *AWSServiceMock) UploadDirectoryToS3Batch(awsConfig cloud.Config, sourceDir, targetPath, bucket string) error {
	return m.UploadDirectoryToS3BatchFunc(awsConfig, sourceDir, targetPath, bucket)
}

func (m *AWSServiceMock) InvalidateCloudFront(awsConfig cloud.Config, distributionID string) error {
	return m.InvalidateCloudFrontFunc(awsConfig, distributionID)
}
