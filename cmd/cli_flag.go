package cmd

type CliFlag = string

const (
	FlagProject          CliFlag = "project"
	FlagProjectShort     CliFlag = "n"
	FlagProvider         CliFlag = "provider"
	FlagProviderShort    CliFlag = "p"
	FlagRegion           CliFlag = "region"
	FlagRegionShort      CliFlag = "r"
	FlagBranch           CliFlag = "branch"
	FlagBranchShort      CliFlag = "b"
	FlagDomain           CliFlag = "domain"
	FlagDomainShort      CliFlag = "d"
	FlagForce            CliFlag = "force"
	FlagForceShort       CliFlag = "f"
	FlagDeployNomad      CliFlag = "deploy-nomad"
	FlagLinuxDistro      CliFlag = "linux-distro"
	FlagLinuxDistroShort CliFlag = "l"
	FlagEdition          CliFlag = "edition"
	FlagEditionShort     CliFlag = "e"

	FlagGCPParentProject CliFlag = "gcp-parent-project"
	FlagGCPDnsZone       CliFlag = "gcp-dns-zone"
	FlagGCPOrgID         CliFlag = "gcp-org-id"
	FlagGCPBillingID     CliFlag = "gcp-billing-account-id"

	FlagAZResourceGroup  CliFlag = "az-resource-group"
	FlagAZSubscriptionID CliFlag = "az-subscription-id"
	FlagAZTenantID       CliFlag = "az-tenant-id"
	FlagAZLoginViaCLI    CliFlag = "az-use-cli"
)
