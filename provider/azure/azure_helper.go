package azure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/authorization/armauthorization"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/date"
	uuid "github.com/satori/go.uuid"

	"github.com/Azure/go-autorest/autorest/to"
)

// CreateResourceGroup az group create --name "$RESOURCE_GROUP" --location "$LOCATION".
func (a Azure) CreateResourceGroup(resourceGroupName, subscriptionID, location string) error {
	if err := a.CheckResourceGroup(resourceGroupName, subscriptionID); err == nil {
		fmt.Printf("resource group [%s] already exists\n", resourceGroupName)
		return nil
	}

	fmt.Printf("creating resource group [%s] in location [%s]\n", resourceGroupName, location)
	c := armresources.NewResourceGroupsClient(a.AzureArmConnection, subscriptionID)
	ctx := context.TODO()
	_, err := c.CreateOrUpdate(
		ctx,
		resourceGroupName,
		armresources.ResourceGroup{
			Location: to.StringPtr(location),
			Tags:     map[string]*string{"owner": to.StringPtr("caravan-cli")},
		},
		nil)
	return err
}

func (a Azure) CheckResourceGroup(resourceGroupName, subscriptionID string) error {
	fmt.Printf("checking existence of resource group [%s]\n", resourceGroupName)
	c := armresources.NewResourceGroupsClient(a.AzureArmConnection, subscriptionID)
	ctx := context.TODO()
	res, err := c.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}

	if res.RawResponse.StatusCode != 204 {
		return fmt.Errorf("resource group [%s] not existing", resourceGroupName)
	}
	return nil
}

// CreateStorageAccount az storage account create --name "$STORAGE_ACCOUNT" --resource-group "$RESOURCE_GROUP" --location "$LOCATION" --tags "owner=$OWNER".
func (a Azure) CreateStorageAccount(subscriptionID, storageAccountName, resourceGroupName, location string) error {
	if err := a.CheckStorageAccount(subscriptionID, storageAccountName, resourceGroupName); err == nil {
		fmt.Printf("storage account [%s] already exists\n", storageAccountName)
		return nil
	}

	c := armstorage.NewStorageAccountsClient(a.AzureArmConnection, subscriptionID)
	ctx := context.TODO()
	res, err := c.CheckNameAvailability(ctx, armstorage.StorageAccountCheckNameAvailabilityParameters{
		Name: to.StringPtr(storageAccountName),
		Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
	},
		nil)
	if err != nil {
		return err
	}

	if !*res.NameAvailable {
		return fmt.Errorf(
			"storage account name [%s] not available: %w\nserver message: %v",
			storageAccountName, err, *res.Message)
	}

	fmt.Printf("creating storage account [%s]\n", storageAccountName)
	future, err := c.BeginCreate(
		ctx,
		resourceGroupName,
		storageAccountName,
		armstorage.StorageAccountCreateParameters{
			SKU:        &armstorage.SKU{Name: armstorage.SKUNameStandardLRS.ToPtr()},
			Kind:       armstorage.KindStorage.ToPtr(),
			Location:   to.StringPtr(location),
			Properties: nil,
		},
		nil)

	if err != nil {
		return fmt.Errorf("failed to start creating storage account: %w", err)
	}

	_, err = future.PollUntilDone(ctx, time.Second*5)
	if err != nil {
		return fmt.Errorf("failed to finish creating storage account: %w", err)
	}
	return nil
}

func (a Azure) CheckStorageAccount(subscriptionID, storageAccountName, resourceGroupName string) error {
	fmt.Printf("checking existence of storage account [%s] in resource group [%s]\n", storageAccountName, resourceGroupName)
	c := armstorage.NewStorageAccountsClient(a.AzureArmConnection, subscriptionID)
	ctx := context.TODO()
	_, err := c.GetProperties(ctx, resourceGroupName, storageAccountName, nil)
	return err
}

// CreateStorageContainer az storage container create --name "$CONTAINER_NAME" --resource-group "$RESOURCE_GROUP" --account-name "$STORAGE_ACCOUNT".
func (a Azure) CreateStorageContainer(subscriptionID, resourceGroupName, storageAccountName, containerName string) error {
	fmt.Printf("creating storage account container [%s] in [%s]\n", containerName, storageAccountName)
	c := armstorage.NewBlobContainersClient(a.AzureArmConnection, subscriptionID)
	ctx := context.TODO()
	_, err := c.Create(ctx, resourceGroupName, storageAccountName, containerName, armstorage.BlobContainer{}, nil)
	return err
}

// CreateRoleAssignment az role assignment create --scope "/subscriptions/${SUBSCRIPTION_ID}" --role "User Access Administrator" --assignee "$CLIENT_ID".
func (a Azure) CreateRoleAssignment(subscriptionID, scope, roleName, principalID string) error {
	fmt.Printf("creating role assignment for principal [%s] with role [%s] and scope [%s]\n", principalID, roleName, scope)
	ctx := context.TODO()

	c := armauthorization.NewRoleAssignmentsClient(a.AzureArmConnection, subscriptionID)
	c2 := armauthorization.NewRoleDefinitionsClient(a.AzureArmConnection)

	res := c2.List(scope, &armauthorization.RoleDefinitionsListOptions{Filter: to.StringPtr(fmt.Sprintf("roleName eq '%s'", roleName))})
	if res.Err() != nil {
		return res.Err()
	}
	res.NextPage(ctx)
	roleDefID := res.PageResponse().Value[0].ID

	fmt.Printf("using role definition [%s]\n", to.String(roleDefID))

	_, err := c.Create(
		ctx,
		scope,
		uuid.NewV1().String(),
		armauthorization.RoleAssignmentCreateParameters{
			Properties: &armauthorization.RoleAssignmentProperties{
				RoleDefinitionID: roleDefID,
				PrincipalID:      to.StringPtr(principalID),
				PrincipalType:    armauthorization.PrincipalTypeServicePrincipal.ToPtr(),
			}}, nil)

	return err
}

func (a Azure) CreateServicePrincipal(tenantID, name string) (string, string, error) {
	appID := ""
	objectID, secret, err := a.CheckServicePrincipal(tenantID, name)
	var errNoApp *NoAzureApplicationError
	var errNoSP *NoAzureServicePrincipalError

	if err != nil {
		fmt.Printf("service principal [%s] already existing", name)
		return *objectID, *secret, nil
	} else if ok := errors.As(err, errNoApp); ok {
		fmt.Println(errNoApp.Error())
	} else if ok := errors.As(err, errNoSP); ok {
		fmt.Println(errNoSP.Error())
		appID = errNoSP.appID
	} else {
		return "", "", err
	}

	c := graphrbac.NewServicePrincipalsClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer
	c2 := graphrbac.NewApplicationsClient(tenantID)
	c2.Authorizer = a.AzureGraphAuthorizer
	ctx := context.TODO()

	newSecret := uuid.NewV1().String()
	if appID == "" {
		fmt.Printf("creating ad application with name [%s]\n", name)
		app, err := c2.Create(ctx, graphrbac.ApplicationCreateParameters{
			DisplayName:             to.StringPtr(name),
			AvailableToOtherTenants: to.BoolPtr(false),
		})
		if err != nil {
			return "", "", err
		}
		appID = *app.AppID
	}

	fmt.Printf("creating ad service principal for application [%s]\n", name)
	time.Now().Add(time.Hour * 24 * 365)
	sp, err := c.Create(ctx, graphrbac.ServicePrincipalCreateParameters{
		AppID:          &appID,
		AccountEnabled: to.BoolPtr(true),
		PasswordCredentials: &[]graphrbac.PasswordCredential{{
			Value:     to.StringPtr(newSecret),
			StartDate: &date.Time{Time: time.Now()},
			EndDate:   &date.Time{Time: time.Now().Add(time.Hour * 24 * 365)},
		}},
	})
	if err != nil {
		return "", "", err
	}
	objectID = sp.ObjectID

	// FIXME: LOL, works on my machine with this. Basically there's a sync issue and the service principal is not properly
	//  propagated within Azure. Adding this sleep we increase the possibility of the SP being available for other API
	//  calls. Note: adding a c.Get(x,x) does not solve the problem, given it is immediately available with that API anyway.
	// time.Sleep(60 * time.Second)

	return *objectID, newSecret, nil
}

func (a Azure) CheckServicePrincipal(tenantID, name string) (*string, *string, error) {
	c := graphrbac.NewServicePrincipalsClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer
	c2 := graphrbac.NewApplicationsClient(tenantID)
	c2.Authorizer = a.AzureGraphAuthorizer
	ctx := context.TODO()

	filterQuery := fmt.Sprintf("displayName eq '%s'", name)

	res, err := c2.List(ctx, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	err = res.NextWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	apps := res.Values()
	if len(apps) == 0 {
		return nil, nil, fmt.Errorf("no application existing with name [%s]", name)
	}

	res2, err := c.List(ctx, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	err = res2.NextWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	sps := res2.Values()
	if len(sps) == 0 {
		return nil, nil, fmt.Errorf("no service principal found with name [%s]", name)
	}

	sp := sps[0]
	objID := sp.ObjectID
	secret := (*sp.PasswordCredentials)[0].Value

	return objID, secret, nil
}

// CreateADPermission az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 1cda74f2-2616-4834-b122-5cb1b07f8a59=Role.
func (a Azure) CreateADPermission(tenantID, clientID, scope string) error {
	fmt.Printf("creating ad permission for client [%s] with scope [%s]\n", clientID, scope)
	c := graphrbac.NewOAuth2PermissionGrantClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer
	ctx := context.TODO()
	_, err := c.Create(ctx, &graphrbac.OAuth2PermissionGrant{
		ClientID:   to.StringPtr(clientID),
		ResourceID: to.StringPtr("00000002-0000-0000-c000-000000000000"),
		Scope:      to.StringPtr(scope),
	})
	if err != nil {
		return err
	}
	return nil
}

func (a Azure) SetupAuthorization(useCLI bool) (autorest.Authorizer, error) {
	if useCLI {
		return auth.NewAuthorizerFromCLI()
	} else {
		return auth.NewAuthorizerFromEnvironment()
	}
}

func (a Azure) SetupAuthorizationWithResource(useCLI bool, resource string) (autorest.Authorizer, error) {
	if useCLI {
		return auth.NewAuthorizerFromCLIWithResource(resource)
	} else {
		return auth.NewAuthorizerFromEnvironmentWithResource(resource)
	}
}

func (a Azure) SetupConnection(useCLI bool) (*arm.Connection, error) {
	var credential azcore.TokenCredential
	var err error
	if useCLI {
		credential, err = azidentity.NewAzureCLICredential(nil)
	} else {
		credential, err = azidentity.NewEnvironmentCredential(nil)
	}
	if err != nil {
		return nil, err
	}
	return arm.NewDefaultConnection(credential, nil), nil
}
