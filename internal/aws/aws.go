package aws

import (
	"context"
	"errors"
	"fmt"
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
	CaravanConfig caravan.Config
	AWSConfig     aws.Config
}

func NewAWS(conf caravan.Config) (a AWS) {

	a.CaravanConfig = conf
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(a.CaravanConfig.Region),
	)

	if err != nil {
		fmt.Printf("error creating config: %s\n", err)
		return
	}
	a.AWSConfig = cfg
	return a
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
				LocationConstraint: s3types.BucketLocationConstraint(a.CaravanConfig.Region),
			},
		})

	if err != nil {
		if !errors.As(err, &bae) && !errors.As(err, &bao) {
			return err
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
		return fmt.Errorf("unable to enable versioning: %s\n", err)
	}

	return nil
}

func (a *AWS) DeleteBucket(name string) (err error) {

	svc := s3.NewFromConfig(a.AWSConfig)
	_, err = svc.DeleteBucket(
		context.TODO(),
		&s3.DeleteBucketInput{
			Bucket: &name,
		})
	if err != nil {
		return err
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
					dytypes.KeySchemaElement{
						KeyType:       dytypes.KeyTypeHash,
						AttributeName: aws.String("LockID"),
					},
				},
				AttributeDefinitions: []dytypes.AttributeDefinition{
					dytypes.AttributeDefinition{
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
			return err
		}
		if i >= retry {
			return fmt.Errorf("maximum number of retries reached: %d\n", retry)
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

			return err
		}
	}

	if i >= retry {
		return fmt.Errorf("maximum number of retries reached: %d\n", retry)
	}
	return nil

}
