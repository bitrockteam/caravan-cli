package caravan

type Config struct {
	Name           string
	Region         string
	Workdir        string
	WorkdirProject string
	Profile        string
	Provider       string
	TableName      string
	BucketName     string
}

// Validate validate the configuration for the constraints
func (c Config) Validate() (err error) {
	//TODO validate configuration for different cloud providers
	return nil
}
