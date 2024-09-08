package service

import (
	"errors"
	"testing"

	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/artcurty/kick-it-to-aws/internal/service/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/golang/mock/gomock"
)

func TestInvalidateCloudFrontCache(t *testing.T) {
	type args struct {
		awsConfig      *cloud.Config
		distributionID string
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		mockCFClient func(*mock.MockCloudFrontClient)
	}{
		{
			name: "Successful invalidation",
			args: args{
				awsConfig:      &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				distributionID: "validDistributionID",
			},
			wantErr: false,
			mockCFClient: func(m *mock.MockCloudFrontClient) {
				m.EXPECT().CreateInvalidation(gomock.Any(), gomock.Any()).Return(&cloudfront.CreateInvalidationOutput{}, nil)
			},
		},
		{
			name: "Invalid distribution ID",
			args: args{
				awsConfig:      &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				distributionID: "",
			},
			wantErr: true,
			mockCFClient: func(m *mock.MockCloudFrontClient) {
				m.EXPECT().CreateInvalidation(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid distribution ID"))
			},
		},
		{
			name: "Failed invalidation",
			args: args{
				awsConfig:      &cloud.Config{Config: aws.Config{Region: "us-west-2"}},
				distributionID: "invalidDistributionID",
			},
			wantErr: true,
			mockCFClient: func(m *mock.MockCloudFrontClient) {
				m.EXPECT().CreateInvalidation(gomock.Any(), gomock.Any()).Return(nil, errors.New("invalidation error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCFClient := mock.NewMockCloudFrontClient(ctrl)

			cloudFrontClientFactory = func(cfg aws.Config) CloudFrontClient {
				return mockCFClient
			}

			tt.mockCFClient(mockCFClient)

			if err := InvalidateCloudFrontCache(tt.args.awsConfig, tt.args.distributionID); (err != nil) != tt.wantErr {
				t.Errorf("InvalidateCloudFrontCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
