// Provider interface to manage different cloud providers.
package provider

import (
	"caravan-cli/cli"
	"context"
)

const (
	AWS   = "aws"
	GCP   = "gcp"
	Azure = "azure"
)

type WithDeploy interface {
	// Deploy will execute the operations needed to deploy the different stack layers
	Deploy(context.Context, cli.DeployLayer) error
}

type WithBake interface {
	// Bake will execute the image baking procedures
	Bake(context.Context) error
}

type WithStatus interface {
	// Status will output the current state of Caravan
	Status(context.Context) error
}

type WithDestroy interface {
	// Destroy will execute the operations needed to destroy the different stack layers
	Destroy(context.Context, cli.DeployLayer) error
}

type Provider interface {
	// GetTemplates returns the templates needed by the provider. The caller will handle persistence of the files.
	GetTemplates(context.Context) ([]cli.Template, error)

	// ValidateConfiguration performs a check on the configuration provided to the Provider implementation. For example it
	// might check that the provided instance size is valid
	ValidateConfiguration(context.Context) error

	// InitProvider creates baseline resources like state stores, lock, projects, etc...
	InitProvider(context.Context) error

	WithBake

	WithDeploy

	WithDestroy

	// CleanProvider deletes cloud resources created during InitProvider
	CleanProvider(context.Context) error

	WithStatus

	// Update upgrades versions etc...
	// Update() error
}
