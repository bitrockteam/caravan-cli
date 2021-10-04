package azure

import "context"

type HelperInterface interface {
	CreateResourceGroup(ctx context.Context, resourceGroupName, subscriptionID, location string) error
	CreateStorageAccount(ctx context.Context, subscriptionID, storageAccountName, resourceGroupName, location string) error
	CreateStorageContainer(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName string) error
	CreateRoleAssignment(ctx context.Context, subscriptionID, scope, roleName, principalID string) error
	CreateServicePrincipal(ctx context.Context, tenantID, name string) (string, string, error)
	CreateADPermission(ctx context.Context, tenantID, clientID, scope string) error
}
