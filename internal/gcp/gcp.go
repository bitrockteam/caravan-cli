package gcp

import (
	"caravan/internal/caravan"
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/grpc/status"
)

type GCP struct {
	Caravan   caravan.Config
	Templates []caravan.Template
}

func New(c caravan.Config) (g GCP, err error) {
	return GCP{Caravan: c}, nil
}

func (g GCP) Init() error {
	if err := g.CreateProject(g.Caravan.Name, g.Caravan.GCPOrgID); err != nil {
		return err
	}
	return nil
}

func (g GCP) Clean() error {
	if err := g.DeleteProject(g.Caravan.Name); err != nil {
		return err
	}
	return nil
}

func (g GCP) CreateBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) DeleteBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}
func (g GCP) EmptyBucket(name string) error {
	fmt.Printf("NOP\n")
	return nil
}
func (g GCP) CreateLockTable(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) DeleteLockTable(name string) error {
	fmt.Printf("NOP\n")
	return nil
}

func (g GCP) GenerateConfig() error {
	fmt.Printf("NOP\n")
	return nil
}

// CreateProject creates a project in GCP and waits for the completion.
func (g GCP) CreateProject(id, orgID string) error {
	ctx := context.Background()
	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get resourcemanager: %w", err)
	}

	p := &cloudresourcemanager.Project{
		ProjectId:   id,
		Parent:      orgID,
		DisplayName: id,
	}

	op, err := cloudresourcemanagerService.Projects.Create(p).Context(ctx).Do()
	if err != nil {
		fmt.Printf("error creating project %s, %s: %s\n", orgID, id, err)
		s, _ := status.FromError(err)
		if strings.Contains(s.Message(), "alreadyExists") {
			return nil
		}
		return err
	}

	for i := 0; i < 10; i++ {
		resp, err := cloudresourcemanagerService.Operations.Get(op.Name).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("error getting operation %s/%s: %s", id, orgID, err)
		}
		if resp.Done {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timed out creating project %s", id)
}

// DeleteProject deletes a project from its project id.
func (g GCP) DeleteProject(name string) error {
	ctx := context.Background()

	cloudresourcemanagerService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return fmt.Errorf("unable to get resourcemanager: %w", err)
	}

	// check project name
	resp, err := cloudresourcemanagerService.Projects.Search().Query("id:" + name).Context(ctx).Do()
	if err != nil {
		return err
	}

	// project doesn't exists
	if len(resp.Projects) == 0 {
		fmt.Printf("project doesn't exist: %s\n", name)
		return nil
	}

	if len(resp.Projects) != 1 {
		return fmt.Errorf("unable to uniquely identify project: %s (%d)", name, len(resp.Projects))
	}

	if resp.Projects[0].State == "DELETE_REQUESTED" {
		return nil
	}

	_, err = cloudresourcemanagerService.Projects.Delete(resp.Projects[0].Name).Context(ctx).Do()
	if err != nil {

		// check project name
		resp, err := cloudresourcemanagerService.Projects.Search().Query("id:" + name).Context(ctx).Do()
		if err != nil {
			return err
		}

		if len(resp.Projects) != 1 {
			return fmt.Errorf("unable to uniquely identify project: %s (%d)", name, len(resp.Projects))
		}

		// project doesn't exists
		if len(resp.Projects) == 0 {
			fmt.Printf("project doesn't exist: %s\n", name)
			return nil
		}

		if resp.Projects[0].State == "DELETE_REQUESTED" {
			return nil
		}

		return err
	}
	return nil
}
