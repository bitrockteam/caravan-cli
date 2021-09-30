package cli

type AzureConfig struct {
	AzureBakingSubscriptionID string `json:",omitempty"`
	AzureBakingResourceGroup  string `json:",omitempty"`
	AzureBakingClientID       string `json:",omitempty"`
	AzureBakingClientSecret   string `json:",omitempty"`
	AzureDNSResourceGroup     string `json:",omitempty"`
	AzureResourceGroup        string `json:",omitempty"`
	AzureStorageAccount       string `json:",omitempty"`
	AzureStorageContainerName string `json:",omitempty"`
	AzureClientID             string `json:",omitempty"`
	AzureClientSecret         string `json:",omitempty"`
	AzureTenantID             string `json:",omitempty"`
	AzureSubscriptionID       string `json:",omitempty"`
	AzureUseCLI               bool   `json:",omitempty"`
}

func (c *Config) SetAzureBakingSubscriptionID(s string) {
	c.AzureBakingResourceGroup = s
}

func (c *Config) SetAzureBakingResourceGroup(s string) {
	c.AzureBakingResourceGroup = s
}

func (c *Config) SetAzureBakingClientID(s string) {
	c.AzureBakingClientID = s
}

func (c *Config) SetAzureBakingClientSecret(s string) {
	c.AzureBakingClientSecret = s
}

func (c *Config) SetAzureDNSResourceGroup(s string) {
	c.AzureDNSResourceGroup = s
}

func (c *Config) SetAzureResourceGroup(s string) {
	c.AzureResourceGroup = s
}

func (c *Config) SetAzureStorageAccount(s string) {
	c.AzureStorageAccount = s
}
func (c *Config) SetAzureStorageContainerName(s string) {
	c.AzureStorageContainerName = s
}
func (c *Config) SetAzureClientID(s string) {
	c.AzureClientID = s
}
func (c *Config) SetAzureClientSecret(s string) {
	c.AzureClientSecret = s
}
func (c *Config) SetAzureTenantID(s string) {
	c.AzureTenantID = s
}
func (c *Config) SetAzureSubscriptionID(s string) {
	c.AzureSubscriptionID = s
}
