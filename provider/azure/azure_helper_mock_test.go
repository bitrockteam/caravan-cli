package azure_test

import (
	"context"
	"fmt"
)

type HelperMock struct {
	Subscriptions                 []string
	SubscriptionToResourceGroups  map[string][]string
	StorageAccountToResourceGroup map[string]string
	StorageAccountToContainer     map[string][]string
	SubscriptionToRoleAssignments map[string][]string
	TenantToServicePrincipals     map[string][]string
	TenantToADPermissions         map[string][]string
}

func NewHelperMock(subscriptionID string) *HelperMock {
	return &HelperMock{
		Subscriptions:                 []string{subscriptionID},
		SubscriptionToResourceGroups:  make(map[string][]string),
		StorageAccountToResourceGroup: make(map[string]string),
		StorageAccountToContainer:     make(map[string][]string),
		SubscriptionToRoleAssignments: make(map[string][]string),
		TenantToServicePrincipals:     make(map[string][]string),
		TenantToADPermissions:         make(map[string][]string),
	}
}

func (h HelperMock) subscriptionExists(id string) bool {
	for _, subscription := range h.Subscriptions {
		if id == subscription {
			return true
		}
	}
	return false
}

func (h HelperMock) CreateResourceGroup(ctx context.Context, resourceGroupName, subscriptionID, location string) error {
	if !h.subscriptionExists(subscriptionID) {
		return fmt.Errorf("subscription [%s] not exists", subscriptionID)
	}
	if val, ok := h.SubscriptionToResourceGroups[subscriptionID]; ok {
		h.SubscriptionToResourceGroups[subscriptionID] = append(val, resourceGroupName)
	} else {
		h.SubscriptionToResourceGroups[subscriptionID] = []string{resourceGroupName}
	}
	fmt.Printf("creating resource group [%s] in subscription [%s]\n", resourceGroupName, subscriptionID)
	return nil
}

func (h HelperMock) CreateStorageAccount(ctx context.Context, subscriptionID, storageAccountName, resourceGroupName, location string) error {
	if _, ok := h.StorageAccountToResourceGroup[storageAccountName]; !ok {
		h.StorageAccountToResourceGroup[storageAccountName] = resourceGroupName
	}
	fmt.Printf("creating storage account [%s] in resource group [%s]\n", storageAccountName, resourceGroupName)
	return nil
}

func (h HelperMock) CreateStorageContainer(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName string) error {
	if val, ok := h.StorageAccountToContainer[storageAccountName]; ok {
		h.StorageAccountToContainer[storageAccountName] = append(val, containerName)
	} else {
		h.StorageAccountToContainer[storageAccountName] = []string{containerName}
	}
	fmt.Printf("creating storage container [%s] in storage account [%s]\n", containerName, storageAccountName)
	return nil
}

func (h HelperMock) CreateRoleAssignment(ctx context.Context, subscriptionID, scope, roleName, principalID string) error {
	if !h.subscriptionExists(subscriptionID) {
		return fmt.Errorf("subscription [%s] not exists", subscriptionID)
	}
	if val, ok := h.SubscriptionToRoleAssignments[subscriptionID]; ok {
		h.SubscriptionToRoleAssignments[subscriptionID] = append(val, fmt.Sprintf("%s:::%s:::%s", scope, roleName, principalID))
	} else {
		h.SubscriptionToRoleAssignments[subscriptionID] = []string{fmt.Sprintf("%s:::%s:::%s", scope, roleName, principalID)}
	}
	fmt.Printf("creating role assignment %s:::%s:::%s in subscription [%s]\n", scope, roleName, principalID, subscriptionID)
	return nil
}

func (h HelperMock) CreateServicePrincipal(ctx context.Context, tenantID, name string) (string, string, error) {
	if val, ok := h.TenantToServicePrincipals[tenantID]; ok {
		h.TenantToServicePrincipals[tenantID] = append(val, name)
	} else {
		h.TenantToServicePrincipals[tenantID] = []string{name}
	}
	fmt.Printf("creating service principal [%s] in tenant [%s]\n", name, tenantID)
	return "id", "pass", nil
}

func (h HelperMock) CreateADPermission(ctx context.Context, tenantID, clientID, scope string) error {
	if val, ok := h.TenantToADPermissions[tenantID]; ok {
		h.TenantToADPermissions[tenantID] = append(val, fmt.Sprintf("%s:::%s", clientID, scope))
	} else {
		h.TenantToADPermissions[tenantID] = []string{fmt.Sprintf("%s:::%s", clientID, scope)}
	}
	fmt.Printf("creating ad permission %s:::%s in tenant [%s]\n", clientID, scope, tenantID)
	return nil
}
