package aws_test

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"testing"

	"github.com/google/uuid"
)

func TestBucketAWS(t *testing.T) {

	uid := uuid.New().String()
	aws := aws.NewAWS(caravan.Config{
		Region: "eu-south-1",
	})
	err := aws.CreateBucket("caravan-aws-test-" + uid)

	if err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}

	err = aws.DeleteBucket("caravan-aws-test-" + uid)
	if err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestIdempotentBucketAWS(t *testing.T) {

	uid := uuid.New().String()
	aws := aws.NewAWS(caravan.Config{
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
}

func TestLockTableAWS(t *testing.T) {

	aws := aws.NewAWS(caravan.Config{
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

	aws := aws.NewAWS(caravan.Config{
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
