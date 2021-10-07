//go:build integration
// +build integration

package aws_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/aws"
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestStateStore(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(ctx, &cli.Config{
		Name:   "test-" + id,
		Region: "us-east-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}
	if err := aws.CreateStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}
	if err := aws.CreateStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error idempotent creating bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestEmptyStateStore(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(ctx, &cli.Config{
		Name:   "test-" + id,
		Region: "eu-south-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}
	if err := aws.CreateStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating bucket: %s\n", err)
	}
	if err := aws.EmptyStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error emptying bucket: %s\n", err)
	}
	if err := aws.DeleteStateStore(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error deleting bucket: %s\n", err)
	}
}

func TestLock(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New().String()
	id := strings.Split(uid, "-")[3]
	aws, err := aws.New(ctx, &cli.Config{
		Name:   "test-" + id,
		Region: "eu-south-1",
	})
	if err != nil {
		t.Fatalf("error creating aws config: %s\n", err)
	}

	if err := aws.CreateLock(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.CreateLock(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.DeleteLock(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
	if err := aws.DeleteLock(ctx, aws.Caravan.Name); err != nil {
		t.Fatalf("error creating lock: %s\n", err)
	}
}
