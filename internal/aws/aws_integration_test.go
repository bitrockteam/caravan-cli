// +build integration

package aws_test

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"testing"

	"github.com/google/uuid"
)

func TestBucketAWS(t *testing.T) {
	uid := uuid.New().String()
	c, err := caravan.NewConfigFromScratch("name", "aws", "")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	aws, _ := aws.NewAWS(*c)
	if err := aws.CreateBucket("caravan-aws-test-" + uid); err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}

	if err := aws.DeleteBucket("caravan-aws-test-" + uid); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestIdempotentBucketAWS(t *testing.T) {
	uid := uuid.New().String()
	aws, _ := aws.NewAWS(caravan.Config{
		Region: "eu-south-1",
	})
	err := aws.CreateBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}

	err = aws.CreateBucket("caravan-aws-test-" + uid)

	if err != nil {
		t.Fatalf("error idempotent creating bucket: %s\n", err)
	}

	err = aws.DeleteBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}

	err = aws.DeleteBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestEmptyBucketAWS(t *testing.T) {
	uid := uuid.New().String()

	aws, _ := aws.NewAWS(caravan.Config{
		Region: "eu-south-1",
	})

	err := aws.CreateBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}

	err = aws.EmptyBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error emptying bucket: %s\n", err)
	}

	err = aws.DeleteBucket("caravan-aws-test-" + uid)

	if err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestLockTableAWS(t *testing.T) {
	aws, _ := aws.NewAWS(caravan.Config{
		Region: "eu-south-1",
	})
	uid := uuid.New().String()
	err := aws.CreateLockTable("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	err = aws.DeleteLockTable("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
}

func TestIdempotentLockTableAWS(t *testing.T) {
	aws, _ := aws.NewAWS(caravan.Config{
		Region: "eu-south-1",
	})
	uid := uuid.New().String()
	err := aws.CreateLockTable("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	err = aws.CreateLockTable("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	err = aws.DeleteLockTable("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
}
