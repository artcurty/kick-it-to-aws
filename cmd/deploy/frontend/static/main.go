package main

import (
	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/artcurty/kick-it-to-aws/internal/service"
	"log"
	"os"
)

func main() {
	awsConfig, err := cloud.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	s3Path := os.Getenv("S3_PATH")
	staticFilesPath := os.Getenv("S3_FILES_PATH")
	cloudFrontDistributionID := os.Getenv("CLOUDFRONT_DISTRIBUTION_ID")
	s3Bucket := os.Getenv("S3_BUCKET")

	err = service.UploadDirectoryToS3Batch(awsConfig, staticFilesPath, s3Path, s3Bucket)
	if err != nil {
		log.Fatalf("Failed to upload to S3: %v", err)
	}

	err = service.InvalidateCloudFrontCache(awsConfig, cloudFrontDistributionID)
	if err != nil {
		log.Fatalf("Failed to invalidate CloudFront: %v", err)
	}

	log.Println("Frontend successfully deployed!")
}
