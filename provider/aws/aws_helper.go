package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	types2 "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (a AWS) CreateStateStore(ctx context.Context, name string) (err error) {
	var bae *types.BucketAlreadyExists
	var bao *types.BucketAlreadyOwnedByYou

	svc := s3.NewFromConfig(a.AWSConfig)

	if a.Caravan.Region != "us-east-1" {
		_, err = svc.CreateBucket(
			ctx,
			&s3.CreateBucketInput{
				Bucket: aws2.String(name),
				CreateBucketConfiguration: &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraint(a.Caravan.Region),
				},
			})
	} else {
		_, err = svc.CreateBucket(
			ctx,
			&s3.CreateBucketInput{
				Bucket: aws2.String(name),
			})
	}

	if err != nil {
		if !errors.As(err, &bae) && !errors.As(err, &bao) {
			return fmt.Errorf("unable to create bucket %s: %w", name, err)
		}
	}

	_, err = svc.PutBucketVersioning(
		ctx,
		&s3.PutBucketVersioningInput{
			Bucket: aws2.String(name),
			VersioningConfiguration: &types.VersioningConfiguration{
				Status: types.BucketVersioningStatusEnabled,
			},
		})

	if err != nil {
		return fmt.Errorf("unable to enable versioning on bucket %s: %w", name, err)
	}

	return nil
}

func (a AWS) EmptyStateStore(ctx context.Context, name string) (err error) {
	var nsb *types.NoSuchBucket

	svc := s3.NewFromConfig(a.AWSConfig)
	vers, err := svc.ListObjectVersions(
		ctx,
		&s3.ListObjectVersionsInput{
			Bucket: &name,
		})
	if err != nil {
		if errors.As(err, &nsb) {
			return nil
		}
		return fmt.Errorf("error listing object versions: %w", err)
	}

	for _, k := range vers.Versions {
		_, err := svc.DeleteObject(
			ctx,
			&s3.DeleteObjectInput{
				Bucket:    &name,
				Key:       k.Key,
				VersionId: k.VersionId,
			})
		if err != nil {
			if !errors.As(err, &nsb) {
				return fmt.Errorf("error deleting object %v: %w", k.Key, err)
			}
		}
	}
	for _, k := range vers.DeleteMarkers {
		_, err := svc.DeleteObject(
			ctx,
			&s3.DeleteObjectInput{
				Bucket:    &name,
				Key:       k.Key,
				VersionId: k.VersionId,
			})
		if err != nil {
			if !errors.As(err, &nsb) {
				return fmt.Errorf("error removing delete marker %v: %w", k.Key, err)
			}
		}
	}
	return nil
}

func (a AWS) DeleteStateStore(ctx context.Context, name string) (err error) {
	var nsb *types.NoSuchBucket

	svc := s3.NewFromConfig(a.AWSConfig)
	_, err = svc.DeleteBucket(
		ctx,
		&s3.DeleteBucketInput{
			Bucket: &name,
		})
	if err != nil {
		// TODO why is error.As not working as expected ?
		if !strings.Contains(err.Error(), "NoSuchBucket") && !errors.As(err, &nsb) {
			return fmt.Errorf("error deleting bucket %s: %w", name, err)
		}
	}
	return nil
}

func (a AWS) CreateLock(ctx context.Context, name string) (err error) {
	var riu *types2.ResourceInUseException

	retry := 10
	sleep := 1
	svc := dynamodb.NewFromConfig(a.AWSConfig)
	i := 0
	for i = 0; i <= retry; i++ {
		_, err = svc.CreateTable(
			ctx,
			&dynamodb.CreateTableInput{
				TableName: aws2.String(name),
				KeySchema: []types2.KeySchemaElement{
					{
						KeyType:       types2.KeyTypeHash,
						AttributeName: aws2.String("LockID"),
					},
				},
				AttributeDefinitions: []types2.AttributeDefinition{
					{
						AttributeName: aws2.String("LockID"),
						AttributeType: types2.ScalarAttributeTypeS,
					},
				},
				BillingMode: types2.BillingModePayPerRequest,
			})

		if err != nil {
			if errors.As(err, &riu) {
				time.Sleep(time.Duration(sleep) * time.Second)
				continue
			}
			return fmt.Errorf("error creating table %s: %w", name, err)
		}
		if i >= retry {
			return fmt.Errorf("maximum number of retries reached creating table %s: %d", name, retry)
		}
	}
	return nil
}

func (a AWS) DeleteLock(ctx context.Context, name string) (err error) {
	var riu *types2.ResourceInUseException
	var rnf *types2.ResourceNotFoundException

	retry := 10
	sleep := 1
	svc := dynamodb.NewFromConfig(a.AWSConfig)
	i := 0
	for i = 0; i <= retry; i++ {
		_, err = svc.DeleteTable(
			ctx,
			&dynamodb.DeleteTableInput{
				TableName: aws2.String(name),
			})
		if err != nil {
			if errors.As(err, &riu) {
				time.Sleep(time.Duration(sleep) * time.Second)
				continue
			}
			if errors.As(err, &rnf) {
				return nil
			}

			return fmt.Errorf("unable to delete lock table %s: %w", name, err)
		}
	}

	if i >= retry {
		return fmt.Errorf("maximum number of retries reached deleting %s: %d", name, retry)
	}
	return nil
}
