//go:build integration
// +build integration

package gcp_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/gcp"
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestProject(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := cli.NewConfigFromScratch(name, "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	c.SetGCPOrgID("55685363496")
	c.SetGCPBillingID("016290-A416F4-EC4527")

	ctx := context.Background()
	g, err := gcp.New(ctx, c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	if err := g.CreateProject(ctx, name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to create project: %s\n", err)
	}
	// idempotence
	if err = g.CreateProject(ctx, name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to create project: %s\n", err)
	}
	p, err := g.GetProject(ctx, name, c.GCPOrgID)
	if err != nil {
		t.Fatalf("unable to Get project: %s\n", err)
	}

	if err := g.SetBillingAccount(ctx, name, c.GCPBillingID); err != nil {
		t.Errorf("unable to set billing ID: %s", err)
	}

	if p.Name != name && p.Parent != c.GCPOrgID {
		t.Errorf("want %s,%s got %s,%s\n", name, c.GCPOrgID, p.Name, p.Parent)
	}

	if err := g.EnableServiceAccess(ctx, name, []string{"serviceusage.googleapis.com", "compute.googleapis.com"}); err != nil {
		t.Errorf("unable to enable access to services: %s\n", err)
	}

	if err := g.EnableServiceAccess(ctx, name, []string{"serviceusage.googleapis.com", "compute.googleapis.com"}); err != nil {
		t.Errorf("unable to enable access to services: %s\n", err)
	}

	if err := g.DeleteProject(ctx, name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}
	//idempotence
	if err := g.DeleteProject(ctx, name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}

}

func TestStateStore(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := cli.NewConfigFromScratch("andrea-test-015", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	ctx := context.Background()
	g, err := gcp.New(ctx, c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	// create bucket
	if err := g.CreateStateStore(ctx, name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
	// idempotence
	if err := g.CreateStateStore(ctx, name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}

	if err := g.WriteStateStore(ctx, name, name, "some data"); err != nil {
		t.Errorf("unable get write: %s", err)
	}

	if err := g.EmptyStateStore(ctx, name); err != nil {
		t.Errorf("unable to empty state store: %s\n", err)
	}

	// delete bucket
	if err := g.DeleteStateStore(ctx, name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
	if err := g.DeleteStateStore(ctx, name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
}

func TestServiceAccount(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := cli.NewConfigFromScratch("andrea-test-015", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	ctx := context.Background()
	g, err := gcp.New(ctx, c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}

	if err := g.CreateServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}
	if err := g.CreateServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}
	_, err = g.CreateServiceAccountKeys(ctx, name, name)
	if err != nil {
		t.Errorf("unable to create service account key: %s\n", err)
	}

	if err := g.DeleteServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
	if err := g.DeleteServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
}

func TestAddPolicy(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	name = "andrea-test-015-terraform"
	c, err := cli.NewConfigFromScratch("andrea-test-015", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	ctx := context.Background()
	g, err := gcp.New(ctx, c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	if err := g.CreateServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}
	member := "serviceAccount:" + name + "@" + g.Caravan.Name + ".iam.gserviceaccount.com"
	if err := g.AddPolicyBinding(ctx, "projects", g.Caravan.Name, member, "roles/owner"); err != nil {
		t.Errorf("unable to add policy binding: %s\n", err)
	}

	if err := g.AddPolicyBinding(ctx, "projects", g.Caravan.Name, member, "roles/owner"); err != nil {
		t.Errorf("unable to add policy binding: %s\n", err)
	}
	if err := g.DeleteServiceAccount(ctx, name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
}

func TestGetUser(t *testing.T) {
	c, err := cli.NewConfigFromScratch("andrea-test-015", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	ctx := context.Background()
	g, err := gcp.New(ctx, c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}

	got, err := g.GetUserEmail("testdata/config_default.golden")
	if err != nil {
		t.Errorf("unable to get user mail: %s\n", err)
	}
	want := "test.user@test.me"
	if got != want {
		t.Errorf("wanted %s got %s\n", want, got)
	}
}
