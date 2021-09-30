// Microsoft Azure provider.
package azure

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
	"fmt"
)

type Azure struct {
	provider.GenericProvider
	AzureHelper *Helper
}

func New(c *cli.Config) (Azure, error) {
	a := Azure{}
	var err error
	a.Caravan = c
	if err = a.ValidateConfiguration(); err != nil {
		return a, err
	}
	if a.AzureHelper, err = NewHelper(a.Caravan.AzureUseCLI); err != nil {
		return a, err
	}
	return a, nil
}

func (a Azure) GetTemplates() ([]cli.Template, error) {
	baking := cli.Template{
		Name: "baking-vars",
		Text: bakingTfVarsTmpl,
		Path: a.Caravan.WorkdirBakingVars,
	}
	infra := cli.Template{
		Name: "infra-vars",
		Text: infraTfVarsTmpl,
		Path: a.Caravan.WorkdirInfraVars,
	}
	infraBackend := cli.Template{
		Name: "infra-backend",
		Text: infraBackendTmpl,
		Path: a.Caravan.WorkdirInfraBackend,
	}
	platform := cli.Template{
		Name: "platform-vars",
		Text: platformTfVarsTmpl,
		Path: a.Caravan.WorkdirPlatformVars,
	}
	platformBackend := cli.Template{
		Name: "platform-backend",
		Text: platformBackendTmpl,
		Path: a.Caravan.WorkdirPlatformBackend,
	}
	applicationSupport := cli.Template{
		Name: "application-vars",
		Text: applicationTfVarsTmpl,
		Path: a.Caravan.WorkdirApplicationVars,
	}
	applicationSupportBackend := cli.Template{
		Name: "application-backend",
		Text: applicationSupportBackendTmpl,
		Path: a.Caravan.WorkdirApplicationBackend,
	}

	return []cli.Template{
		baking,
		infra,
		infraBackend,
		platform,
		platformBackend,
		applicationSupport,
		applicationSupportBackend,
	}, nil
}

func (a Azure) ValidateConfiguration() error {
	//TODO: implement me
	return nil
}

func (a Azure) InitProvider() error {
	var err error
	err = a.AzureHelper.CreateResourceGroup(a.Caravan.AzureResourceGroup, a.Caravan.AzureSubscriptionID, a.Caravan.Region)
	if err != nil {
		return err
	}

	//TODO: create storage account (prefix)sa
	saName := fmt.Sprintf("crv%ssa", a.Caravan.Name)
	err = a.AzureHelper.CreateStorageAccount(a.Caravan.AzureSubscriptionID, saName, a.Caravan.AzureResourceGroup, a.Caravan.Region)
	if err != nil {
		return err
	}
	a.Caravan.SetAzureStorageAccount(saName)

	//TODO: create storage container tfstate
	containerName := "tfstate"
	err = a.AzureHelper.CreateStorageContainer(a.Caravan.AzureSubscriptionID, a.Caravan.AzureResourceGroup, a.Caravan.AzureStorageAccount, containerName)
	if err != nil {
		return err
	}
	a.Caravan.SetAzureStorageContainerName(containerName)

	//TODO: create service principal (prefix)-tf-sp Contributor on the RG + ParentRG
	clientID, clientSecret, err := a.AzureHelper.CreateServicePrincipal(a.Caravan.AzureTenantID, fmt.Sprintf("%s-tf-sp", a.Caravan.Name))
	if err != nil {
		return err
	}
	a.Caravan.SetAzureClientID(clientID)
	a.Caravan.SetAzureClientSecret(clientSecret)
	err = a.AzureHelper.CreateRoleAssignment(
		a.Caravan.AzureSubscriptionID,
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", a.Caravan.AzureSubscriptionID, a.Caravan.AzureResourceGroup),
		"Contributor",
		a.Caravan.AzureClientID,
	)
	if err != nil {
		return err
	}
	if a.Caravan.AzureDNSResourceGroup != a.Caravan.AzureResourceGroup {
		err = a.AzureHelper.CreateRoleAssignment(
			a.Caravan.AzureSubscriptionID,
			fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", a.Caravan.AzureSubscriptionID, a.Caravan.AzureDNSResourceGroup),
			"Contributor",
			a.Caravan.AzureClientID,
		)
		if err != nil {
			return err
		}
	}
	if a.Caravan.AzureBakingResourceGroup != a.Caravan.AzureResourceGroup {
		err = a.AzureHelper.CreateRoleAssignment(
			a.Caravan.AzureSubscriptionID,
			fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", a.Caravan.AzureSubscriptionID, a.Caravan.AzureBakingResourceGroup),
			"Contributor",
			a.Caravan.AzureClientID,
		)
		if err != nil {
			return err
		}
	}
	//TODO: bunch of permissions in AD

	// # Grant Application.ReadWrite.All
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 1cda74f2-2616-4834-b122-5cb1b07f8a59=Role
	err = a.AzureHelper.CreateADPermission(a.Caravan.AzureTenantID, a.Caravan.AzureClientID, "ReadWrite.All")
	if err != nil {
		return err
	}
	// # Grant User.Read
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 311a71cc-e848-46a1-bdf8-97ff7156d8e6=Scope
	err = a.AzureHelper.CreateADPermission(a.Caravan.AzureTenantID, a.Caravan.AzureClientID, "User.Read")
	if err != nil {
		return err
	}
	// # Grant Directory.ReadWrite.All
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 78c8a3c8-a07e-4b9e-af1b-b5ccab50a175=Role
	err = a.AzureHelper.CreateADPermission(a.Caravan.AzureTenantID, a.Caravan.AzureClientID, "Directory.ReadWrite.All")
	if err != nil {
		return err
	}

	// # Apply changes
	// az ad app permission grant --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000

	//TODO: allow access to backend for TF
	err = a.AzureHelper.CreateRoleAssignment(
		a.Caravan.AzureSubscriptionID,
		fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", a.Caravan.AzureSubscriptionID, a.Caravan.AzureResourceGroup),
		"Storage Blob Data Contributor",
		a.Caravan.AzureClientID,
	)
	if err != nil {
		return err
	}

	//TODO: allow assigning roles to other entites for TF
	err = a.AzureHelper.CreateRoleAssignment(
		a.Caravan.AzureSubscriptionID,
		fmt.Sprintf("/subscriptions/%s", a.Caravan.AzureSubscriptionID),
		"User Access Administrator",
		a.Caravan.AzureClientID,
	)
	if err != nil {
		return err
	}

	err = a.Caravan.Save()
	if err != nil {
		return err
	}

	return nil
}

func (a Azure) Bake() error {
	panic("implement me")
}

func (a Azure) Deploy(layer cli.DeployLayer) error {
	panic("implement me")
}

func (a Azure) Destroy(layer cli.DeployLayer) error {
	panic("implement me")
}

func (a Azure) CleanProvider() error {
	//TODO: implement me
	return nil
}

func (a Azure) Status() error {
	panic("implement me")
}
