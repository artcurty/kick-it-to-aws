package service

import (
	"context"
	"fmt"
	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/sync/semaphore"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

var s3ClientFactory = func(cfg aws.Config) S3Client {
	return s3.NewFromConfig(cfg)
}

func UploadDirectoryToS3Batch(awsConfig *cloud.Config, localDir string, s3BasePath string, s3Bucket string) error {
	s3Client := s3ClientFactory(awsConfig.Config)
	files, err := fetchFiles(localDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		s3Key, err := buildS3Key(localDir, file, s3BasePath)
		if err != nil {
			return err
		}

		err = uploadFile(s3Client, s3Bucket, s3Key, file)
		if err != nil {
			return err
		}
	}

	log.Println("All files successfully uploaded to S3!")
	return nil
}

func fetchFiles(localDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error fetching files: %v", err)
	}

	return files, nil
}

func buildS3Key(localDir, filePath, s3BasePath string) (string, error) {
	relativePath, err := filepath.Rel(localDir, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path for %s: %v", filePath, err)
	}

	s3Key := strings.TrimRight(s3BasePath, "/") + "/" + relativePath
	return s3Key, nil
}

func uploadFile(s3Client S3Client, bucketName, s3Key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %v: %v", filePath, err)
	}
	defer file.Close()

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Key),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file %v: %v", filePath, err)
	}

	log.Printf("Successfully uploaded %s to S3 as %s", filePath, s3Key)
	return nil
}

func UploadDirectoryToS3Parallel(awsConfig *cloud.Config, localDir string, s3BasePath string, s3Bucket string, maxParallelUploads int) error {
	s3Client := s3ClientFactory(awsConfig.Config)
	files, err := fetchFiles(localDir)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	sem := semaphore.NewWeighted(int64(maxParallelUploads))

	for _, file := range files {
		wg.Add(1)
		if err := sem.Acquire(context.Background(), 1); err != nil {
			return err
		}
		go func(filePath string) {
			defer wg.Done()
			defer sem.Release(1)
			s3Key, err := buildS3Key(localDir, filePath, s3BasePath)
			if err != nil {
				errChan <- err
				return
			}

			err = uploadFile(s3Client, s3Bucket, s3Key, filePath)
			if err != nil {
				errChan <- err
			}
		}(file)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	log.Println("All files successfully uploaded to S3!")
	return nil
}
