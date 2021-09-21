// Microsoft Azure provider.
package azure

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
)

type Azure struct {
	provider.GenericProvider
}

func New(c *caravan.Config) (Azure, error) {
	a := Azure{}
	a.Caravan = c
	if err := a.ValidateConfiguration(); err != nil {
		return a, err
	}
	return a, nil
}

func (a Azure) GetTemplates() ([]cli.Template, error) {
	baking := caravan.Template{
		Name: "baking-vars",
		Text: bakingTfVarsTmpl,
		Path: a.Caravan.WorkdirBakingVars,
	}
	infra := caravan.Template{
		Name: "infra-vars",
		Text: infraTfVarsTmpl,
		Path: a.Caravan.WorkdirInfraVars,
	}
	infraBackend := caravan.Template{
		Name: "infra-backend",
		Text: infraBackendTmpl,
		Path: a.Caravan.WorkdirInfraBackend,
	}
	platform := caravan.Template{
		Name: "platform-vars",
		Text: platformTfVarsTmpl,
		Path: a.Caravan.WorkdirPlatformVars,
	}
	platformBackend := caravan.Template{
		Name: "platform-backend",
		Text: platformBackendTmpl,
		Path: a.Caravan.WorkdirPlatformBackend,
	}
	applicationSupport := caravan.Template{
		Name: "application-vars",
		Text: applicationTfVarsTmpl,
		Path: a.Caravan.WorkdirApplicationVars,
	}
	applicationSupportBackend := caravan.Template{
		Name: "application-backend",
		Text: applicationSupportBackendTmpl,
		Path: a.Caravan.WorkdirApplicationBackend,
	}

	return []caravan.Template{
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
	//TODO: create resource group (prefix)-rg
	// az group create --name "$RESOURCE_GROUP" --location "$LOCATION" --tags "owner=$OWNER"

	//TODO: create storage account (prefix)sa
	// az storage account create --name "$STORAGE_ACCOUNT" --resource-group "$RESOURCE_GROUP" --location "$LOCATION" --tags "owner=$OWNER"

	//TODO: create storage container tfstate
	// az storage container create --name "$CONTAINER_NAME" --resource-group "$RESOURCE_GROUP" --account-name "$STORAGE_ACCOUNT"

	//TODO: create service principal (prefix)-tf-sp Contributor on the RG + ParetntRG
	// SERVICE_PRINCIPAL=$(az ad sp create-for-rbac \
	//  --name="${PREFIX}-tf-sp" \
	//  --role="Contributor" \
	//  --scopes "/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP}" "/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${PARENT_RESOURCE_GROUP}")

	//TODO: bunch of permissions in AD
	// # Grant Application.ReadWrite.All
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 1cda74f2-2616-4834-b122-5cb1b07f8a59=Role
	// # Grant User.Read
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 311a71cc-e848-46a1-bdf8-97ff7156d8e6=Scope
	// # Grant Directory.ReadWrite.All
	// az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 78c8a3c8-a07e-4b9e-af1b-b5ccab50a175=Role
	// # Apply changes
	// az ad app permission grant --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000

	//TODO: allow access to backend for TF
	// az role assignment create \
	//  --scope "/subscriptions/${SUBSCRIPTION_ID}/resourceGroups/${RESOURCE_GROUP}" \
	//  --role "Storage Blob Data Contributor" \
	//  --assignee "$CLIENT_ID"

	//TODO: allow assigning roles to other entites for TF
	// az role assignment create \
	//  --scope "/subscriptions/${SUBSCRIPTION_ID}" \
	//  --role "User Access Administrator" \
	//  --assignee "$CLIENT_ID"
	return nil
}

func (a Azure) Bake() error {
	panic("implement me")
}

func (a Azure) Deploy(layer caravan.DeployLayer) error {
	panic("implement me")
}

func (a Azure) Destroy(layer caravan.DeployLayer) error {
	panic("implement me")
}

func (a Azure) CleanProvider() error {
	//TODO: implement me
	return nil
}

func (a Azure) Status() error {
	panic("implement me")
}
