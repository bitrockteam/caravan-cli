package caravan

const (
	AWS = "aws"
	GCP = "gcp"
)

type Provider interface {
	Init() error
	GenerateConfig() error
	CreateStateStore(name string) error
	DeleteStateStore(name string) error
	EmptyStateStore(name string) error
	CreateLock(name string) error
	DeleteLock(name string) error
	Clean() error
}
