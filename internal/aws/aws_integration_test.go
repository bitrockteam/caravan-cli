// +build integration

package aws_test

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestStateStore(t *testing.T) {
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(caravan.Config{
		Name:   "test-" + id,
		Region: "eu-south-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}
	if err := aws.CreateStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}
	if err := aws.CreateStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error idempotent creating bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestEmptyStateStore(t *testing.T) {
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(caravan.Config{
		Name:   "test-" + id,
		Region: "eu-south-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}
	if err := aws.CreateStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}
	if err := aws.EmptyStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error emptying bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestLock(t *testing.T) {
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(caravan.Config{
		Name:   "test-" + id,
		Region: "eu-south-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}

	if err := aws.CreateLock(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.CreateLock(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.DeleteLock(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.DeleteLock(aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
}
