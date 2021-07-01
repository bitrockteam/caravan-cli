package aws

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"caravan/internal/caravan"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dytypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AWS struct {
	Caravan   caravan.Config
	AWSConfig aws.Config
	Templates []Template
}

func NewAWS(conf caravan.Config) (a AWS, err error) {
	a.Caravan = conf
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if a.Caravan.Region != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithDefaultRegion(a.Caravan.Region),
		)
	}
	if err != nil {
		return a, err
	}
	if cfg.Region == "" {
		return a, fmt.Errorf("please provide a region")
	}
	a.Caravan.Region = cfg.Region
	a.AWSConfig = cfg
	return a, nil
}

func (a *AWS) CreateBucket(name string) (err error) {
	var bae *s3types.BucketAlreadyExists
	var bao *s3types.BucketAlreadyOwnedByYou

	svc := s3.NewFromConfig(a.AWSConfig)

	_, err = svc.CreateBucket(
		context.TODO(),
		&s3.CreateBucketInput{
			Bucket: aws.String(name),
			CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
				LocationConstraint: s3types.BucketLocationConstraint(a.Caravan.Region),
			},
		})

	if err != nil {
		if !errors.As(err, &bae) && !errors.As(err, &bao) {
			return fmt.Errorf("unable to create bucket %s: %w", name, err)
		}
	}

	_, err = svc.PutBucketVersioning(
		context.TODO(),
		&s3.PutBucketVersioningInput{
			Bucket: aws.String(name),
			VersioningConfiguration: &s3types.VersioningConfiguration{
				Status: s3types.BucketVersioningStatusEnabled,
			},
		})

	if err != nil {
		return fmt.Errorf("unable to enable versioning on bucket %s: %w", name, err)
	}

	return nil
}

func (a *AWS) EmptyBucket(name string) (err error) {
	var nsb *s3types.NoSuchBucket

	svc := s3.NewFromConfig(a.AWSConfig)
	vers, err := svc.ListObjectVersions(
		context.TODO(),
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
			context.TODO(),
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
			context.TODO(),
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

func (a *AWS) DeleteBucket(name string) (err error) {
	var nsb *s3types.NoSuchBucket

	svc := s3.NewFromConfig(a.AWSConfig)
	_, err = svc.DeleteBucket(
		context.TODO(),
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

func (a *AWS) CreateLockTable(name string) (err error) {
	var riu *dytypes.ResourceInUseException

	retry := 10
	sleep := 1
	svc := dynamodb.NewFromConfig(a.AWSConfig)
	i := 0
	for i = 0; i <= retry; i++ {
		_, err = svc.CreateTable(
			context.TODO(),
			&dynamodb.CreateTableInput{
				TableName: aws.String(name),
				KeySchema: []dytypes.KeySchemaElement{
					{
						KeyType:       dytypes.KeyTypeHash,
						AttributeName: aws.String("LockID"),
					},
				},
				AttributeDefinitions: []dytypes.AttributeDefinition{
					{
						AttributeName: aws.String("LockID"),
						AttributeType: dytypes.ScalarAttributeTypeS,
					},
				},
				BillingMode: dytypes.BillingModePayPerRequest,
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

func (a *AWS) DeleteLockTable(name string) (err error) {
	var riu *dytypes.ResourceInUseException
	var rnf *dytypes.ResourceNotFoundException

	retry := 10
	sleep := 1
	svc := dynamodb.NewFromConfig(a.AWSConfig)
	i := 0
	for i = 0; i <= retry; i++ {
		_, err = svc.DeleteTable(
			context.TODO(),
			&dynamodb.DeleteTableInput{
				TableName: aws.String(name),
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
