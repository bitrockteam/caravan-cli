package caravan

type DeployLayer int

const (
	Infrastructure Status = iota
	Platform
	ApplicationSupport
)

type WithDeploy interface {
	// Deploy will execute the operations needed to deploy the different stack layers
	Deploy(DeployLayer) error
}

type WithBake interface {
	// Bake will execute the image baking procedures
	Bake() error
}

type WithStatus interface {
	// Status will output the current state of Caravan
	Status() error
}

type NewProvider interface {
	// GetTemplates returns the templates needed by the provider. The caller will handle persistence of the files.
	GetTemplates() ([]Template, error)

	// ValidateConfiguration performs a check on the configuration provided to the NewProvider implementation. For example it
	// might check that the provided instance size is valid
	ValidateConfiguration() error

	// InitProvider creates baseline resources like state stores, lock, projects, etc...
	InitProvider() error

	WithBake

	WithDeploy

	// Clean deletes everything, including baseline resources
	Clean() error

	WithStatus

	// Update upgrades versions etc...
	// Update() error
}
