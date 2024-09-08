package service

import (
	"context"
	"fmt"
	"log"

	"github.com/artcurty/kick-it-to-aws/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/oklog/ulid/v2"
)

type CloudFrontClient interface {
	CreateInvalidation(ctx context.Context, params *cloudfront.CreateInvalidationInput, optFns ...func(*cloudfront.Options)) (*cloudfront.CreateInvalidationOutput, error)
}

var cloudFrontClientFactory = func(cfg aws.Config) CloudFrontClient {
	return cloudfront.NewFromConfig(cfg)
}

func InvalidateCloudFrontCache(awsConfig *cloud.Config, distributionID string) error {
	cfClient := cloudFrontClientFactory(awsConfig.Config)

	invalidationBatch := &types.InvalidationBatch{
		CallerReference: aws.String(ulid.Make().String()),
		Paths: &types.Paths{
			Quantity: aws.Int32(1),
			Items:    []string{"/*"},
		},
	}

	_, err := cfClient.CreateInvalidation(context.TODO(), &cloudfront.CreateInvalidationInput{
		DistributionId:    &distributionID,
		InvalidationBatch: invalidationBatch,
	})

	if err != nil {
		return fmt.Errorf("failed to invalidate CloudFront cache: %v", err)
	}

	log.Printf("Successfully created CloudFront invalidation for distribution: %s", distributionID)
	return nil
}
