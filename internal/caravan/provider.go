package caravan

type Provider interface {
	GenerateConfig() error
	CreateBucket(name string) error
	CreateLockTable(name string) error
}
