package cmd

type CliFlag = string

const (
	FlagProject       CliFlag = "project"
	FlagProjectShort  CliFlag = "n"
	FlagProvider      CliFlag = "provider"
	FlagProviderShort CliFlag = "p"
	FlagRegion        CliFlag = "region"
	FlagRegionShort   CliFlag = "r"
	FlagBranch        CliFlag = "branch"
	FlagBranchShort   CliFlag = "b"
	FlagDomain        CliFlag = "domain"
	FlagDomainShort   CliFlag = "d"

	FlagGCPParentProject CliFlag = "gcp-parent-project"
	FlagGCPDnsZone       CliFlag = "gcp-dns-zone"

	FlagAZResourceGroup  CliFlag = "az-resource-group"
	FlagAZSubscriptionID CliFlag = "az-subscription-id"
	FlagAZTenantID       CliFlag = "az-tenant-id"
	FlagAZLoginViaCLI    CliFlag = "az-use-cli"
)
