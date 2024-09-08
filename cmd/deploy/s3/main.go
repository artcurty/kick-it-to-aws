package main

import (
	"log"
	"os"
	"strconv"

	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/artcurty/kick-it-to-aws/internal/service"
)

func main() {
	awsConfig, err := cloud.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	sourceDir := os.Getenv("S3_FILES_PATH")
	targetPath := os.Getenv("S3_PATH")
	bucket := os.Getenv("S3_BUCKET")
	maxParallelUploadsStr := os.Getenv("MAX_PARALLEL_UPLOADS")
	maxParallelUploads, err := strconv.Atoi(maxParallelUploadsStr)
	if err != nil {
		log.Fatalf("Invalid value for MAX_PARALLEL_UPLOADS: %v", err)
	}

	err = service.UploadDirectoryToS3Parallel(awsConfig, sourceDir, targetPath, bucket, maxParallelUploads)
	if err != nil {
		log.Fatalf("Failed to upload to S3: %v", err)
	}

	log.Println("Files successfully uploaded to S3!")
}
