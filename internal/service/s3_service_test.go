package service

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/artcurty/kick-it-to-aws/internal/service/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
)

func TestUploadDirectoryToS3Batch(t *testing.T) {
	type args struct {
		awsConfig  *cloud.Config
		localDir   string
		s3BasePath string
		s3Bucket   string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockS3Client func(*mock.MockS3Client)
		setup        func() (string, error)
		teardown     func(string) error
	}{
		{
			name: "Successful upload",
			args: args{
				awsConfig:  &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				localDir:   "testdata",
				s3BasePath: "basepath",
				s3Bucket:   "testbucket",
			},
			wantErr: false,
			mockS3Client: func(m *mock.MockS3Client) {
				m.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(&s3.PutObjectOutput{}, nil).AnyTimes()
			},
			setup: func() (string, error) {
				dir := "testdata"
				err := os.Mkdir(dir, 0755)
				if err != nil {
					return "", err
				}
				file, err := os.Create(filepath.Join(dir, "testfile.txt"))
				if err != nil {
					return "", err
				}
				file.Close()
				return dir, nil
			},
			teardown: func(dir string) error {
				return os.RemoveAll(dir)
			},
		},
		{
			name: "Failed to fetch files",
			args: args{
				awsConfig:  &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				localDir:   "invalidDir",
				s3BasePath: "basepath",
				s3Bucket:   "testbucket",
			},
			wantErr:      true,
			mockS3Client: func(m *mock.MockS3Client) {},
			setup:        func() (string, error) { return "", nil },
			teardown:     func(dir string) error { return nil },
		},
		{
			name: "Failed to upload file",
			args: args{
				awsConfig:  &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				localDir:   "testdata",
				s3BasePath: "basepath",
				s3Bucket:   "testbucket",
			},
			wantErr: true,
			mockS3Client: func(m *mock.MockS3Client) {
				m.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, errors.New("upload error")).AnyTimes()
			},
			setup: func() (string, error) {
				dir := "testdata"
				err := os.Mkdir(dir, 0755)
				if err != nil {
					return "", err
				}
				file, err := os.Create(filepath.Join(dir, "testfile.txt"))
				if err != nil {
					return "", err
				}
				file.Close()
				return dir, nil
			},
			teardown: func(dir string) error {
				return os.RemoveAll(dir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockS3Client := mock.NewMockS3Client(ctrl)

			s3ClientFactory = func(cfg aws.Config) S3Client {
				return mockS3Client
			}

			tt.mockS3Client(mockS3Client)

			if tt.setup != nil {
				dir, err := tt.setup()
				if err != nil {
					t.Fatalf("setup() error = %v", err)
				}
				tt.args.localDir = dir
				defer tt.teardown(dir)
			}

			if err := UploadDirectoryToS3Batch(tt.args.awsConfig, tt.args.localDir, tt.args.s3BasePath, tt.args.s3Bucket); (err != nil) != tt.wantErr {
				t.Errorf("UploadDirectoryToS3Batch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
