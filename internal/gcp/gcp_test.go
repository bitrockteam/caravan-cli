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

	c, err := caravan.NewConfigFromScratch("name", "gcp", "")
	if err != nil {
		t.Fatalf("unable to create config: %s\n", err)
	}
	name := c.Name + "-" + uid

	c.SetGCPOrgID("55685363496")

	g, err := gcp.New(*c)
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

	if err := g.DeleteProject(name); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}
	//idempotence
	if err := g.DeleteProject(name); err != nil {
		t.Fatalf("unable to delete project %s: %s\n", name, err)
	}
}
