package caravan

const (
	AWS = "aws"
	GCP = "gcp"
)

type Provider interface {
	Init() error
	GenerateConfig() error
	CreateBucket(name string) error
	DeleteBucket(name string) error
	EmptyBucket(name string) error
	CreateLockTable(name string) error
	DeleteLockTable(name string) error
	Clean() error
}
