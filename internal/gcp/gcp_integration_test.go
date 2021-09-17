// +build integration

package gcp_test

import (
	"caravan/internal/caravan"
	"caravan/internal/gcp"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestProject(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := caravan.NewConfigFromScratch(name, "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	c.SetGCPOrgID("55685363496")
	c.SetGCPBillingID("016290-A416F4-EC4527")

	g, err := gcp.New(c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	if err := g.CreateProject(name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to create project: %s\n", err)
	}
	// idempotence
	if err = g.CreateProject(name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to create project: %s\n", err)
	}
	p, err := g.GetProject(name, c.GCPOrgID)
	if err != nil {
		t.Fatalf("unable to Get project: %s\n", err)
	}

	if err := g.SetBillingAccount(name, c.GCPBillingID); err != nil {
		t.Errorf("unable to set billing ID: %s", err)
	}

	if p.Name != name && p.Parent != c.GCPOrgID {
		t.Errorf("want %s,%s got %s,%s\n", name, c.GCPOrgID, p.Name, p.Parent)
	}

	if err := g.DeleteProject(name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}
	//idempotence
	if err := g.DeleteProject(name, c.GCPOrgID); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}
}

func TestStateStore(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := caravan.NewConfigFromScratch("andrea-test-003", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	g, err := gcp.New(c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	// create bucket
	if err := g.CreateStateStore(name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
	if err := g.CreateStateStore(name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}

	// delete bucket
	if err := g.DeleteStateStore(name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
	if err := g.DeleteStateStore(name); err != nil {
		t.Errorf("unable to create the bucket: %s", err)
	}
}

func TestServiceAccount(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := caravan.NewConfigFromScratch("andrea-test-005", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	g, err := gcp.New(c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}

	if err := g.CreateServiceAccount(name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}
	if err := g.CreateServiceAccount(name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}

	_, err = g.CreateServiceAccountKeys(name, name)
	if err != nil {
		t.Errorf("unable to create service account key: %s\n", err)
	}

	if err := g.DeleteServiceAccount(name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
	if err := g.DeleteServiceAccount(name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
}

func TestAddPolicy(t *testing.T) {
	uid := strings.Split(uuid.New().String(), "-")[0]
	name := "name-" + uid
	c, err := caravan.NewConfigFromScratch("andrea-test-008", "gcp", "europe-west6")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}

	g, err := gcp.New(c)
	if err != nil {
		t.Fatalf("unable to create GCP: %s\n", err)
	}
	if err := g.CreateServiceAccount(name); err != nil {
		t.Fatalf("unable to create service account: %s\n", err)
	}

	if err := g.AddPolicyBinding("projects", g.Caravan.Name, name, "roles/owner"); err != nil {
		t.Errorf("unable to add policy binding: %s\n", err)
	}

	if err := g.AddPolicyBinding("projects", g.Caravan.Name, name, "roles/owner"); err != nil {
		t.Errorf("idempotent: unable to add policy binding: %s\n", err)
	}

	if err := g.DeleteServiceAccount(name); err != nil {
		t.Fatalf("unable to delete service account: %s\n", err)
	}
}
