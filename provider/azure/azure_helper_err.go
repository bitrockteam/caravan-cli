package azure

import "fmt"

type NoAzureApplicationError struct {
	name string
}

func (e NoAzureApplicationError) Error() string {
	return fmt.Sprintf("azure_helper: no application with name [%s]", e.name)
}

type NoAzureServicePrincipalError struct {
	name  string
	appID string
}

func (e NoAzureServicePrincipalError) Error() string {
	return fmt.Sprintf("azure_helper: no service principal with name [%s] associated to app [%s]", e.name, e.appID)
}

type NoAzureResourceGroupError struct {
	name string
}

func (e NoAzureResourceGroupError) Error() string {
	return fmt.Sprintf("azure_helper: no resource group with name [%s]", e.name)
}

type NoAzureStorageAccountError struct {
	name string
}

func (e NoAzureStorageAccountError) Error() string {
	return fmt.Sprintf("azure_helper: no storage account with name [%s]", e.name)
}

type NoAzureStorageContainerError struct {
	name           string
	storageAccount string
}

func (e NoAzureStorageContainerError) Error() string {
	return fmt.Sprintf("azure_helper: no storage container with name [%s] in account [%s]", e.name, e.storageAccount)
}

type NoAzureRoleAssignmentError struct {
	roleDefinitionID string
	principalID      string
	scope            string
}

func (e NoAzureRoleAssignmentError) Error() string {
	return fmt.Sprintf("azure_helper: no role with id [%s] assigned for principal [%s] with scope [%s]", e.roleDefinitionID, e.principalID, e.scope)
}
