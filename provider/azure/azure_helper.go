package azure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Azure/go-autorest/autorest/azure"

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

type Helper struct {
	AzureArmConnection   *arm.Connection
	AzureGraphAuthorizer autorest.Authorizer
}

func NewHelper(useCLI bool) (*Helper, error) {
	a := &Helper{}
	var err error
	if a.AzureArmConnection, err = setupArmConnection(useCLI); err != nil {
		return nil, err
	}
	if a.AzureGraphAuthorizer, err = setupAuthorizationWithResource(useCLI, azure.PublicCloud.GraphEndpoint); err != nil {
		return nil, err
	}
	return a, nil
}

// CreateResourceGroup az group create --name "$RESOURCE_GROUP" --location "$LOCATION".
func (a Helper) CreateResourceGroup(ctx context.Context, resourceGroupName, subscriptionID, location string) error {
	if err := a.checkResourceGroup(ctx, resourceGroupName, subscriptionID); err == nil {
		log.Info().Msgf("resource group [%s] already exists", resourceGroupName)
		return nil
	}

	log.Info().Msgf("creating resource group [%s] in location [%s]", resourceGroupName, location)
	c := armresources.NewResourceGroupsClient(a.AzureArmConnection, subscriptionID)

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

func (a Helper) checkResourceGroup(ctx context.Context, resourceGroupName, subscriptionID string) error {
	log.Info().Msgf("checking existence of resource group [%s]", resourceGroupName)
	c := armresources.NewResourceGroupsClient(a.AzureArmConnection, subscriptionID)

	res, err := c.CheckExistence(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}

	if res.RawResponse.StatusCode != 204 {
		return NoAzureResourceGroupError{name: resourceGroupName}
	}
	return nil
}

// CreateStorageAccount az storage account create --name "$STORAGE_ACCOUNT" --resource-group "$RESOURCE_GROUP" --location "$LOCATION" --tags "owner=$OWNER".
func (a Helper) CreateStorageAccount(ctx context.Context, subscriptionID, storageAccountName, resourceGroupName, location string) error {
	if err := a.checkStorageAccount(ctx, subscriptionID, storageAccountName, resourceGroupName); err == nil {
		log.Info().Msgf("storage account [%s] already exists", storageAccountName)
		return nil
	}

	c := armstorage.NewStorageAccountsClient(a.AzureArmConnection, subscriptionID)

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

	log.Info().Msgf("creating storage account [%s]", storageAccountName)
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

func (a Helper) checkStorageAccount(ctx context.Context, subscriptionID, storageAccountName, resourceGroupName string) error {
	log.Info().Msgf("checking existence of storage account [%s] in resource group [%s]", storageAccountName, resourceGroupName)
	c := armstorage.NewStorageAccountsClient(a.AzureArmConnection, subscriptionID)

	if _, err := c.GetProperties(ctx, resourceGroupName, storageAccountName, nil); err != nil {
		return NoAzureStorageAccountError{name: storageAccountName}
	} else {
		return nil
	}
}

// CreateStorageContainer az storage container create --name "$CONTAINER_NAME" --resource-group "$RESOURCE_GROUP" --account-name "$STORAGE_ACCOUNT".
func (a Helper) CreateStorageContainer(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName string) error {
	if err := a.checkStorageContainer(ctx, subscriptionID, resourceGroupName, storageAccountName, containerName); err == nil {
		log.Info().Msgf("storage account container [%s] already exists", containerName)
		return nil
	}

	log.Info().Msgf("creating storage account container [%s] in [%s]", containerName, storageAccountName)
	c := armstorage.NewBlobContainersClient(a.AzureArmConnection, subscriptionID)

	_, err := c.Create(ctx, resourceGroupName, storageAccountName, containerName, armstorage.BlobContainer{}, nil)
	return err
}

func (a Helper) checkStorageContainer(ctx context.Context, subscriptionID, resourceGroupName, storageAccountName, containerName string) error {
	log.Info().Msgf("checking existence storage account container [%s] in [%s]", containerName, storageAccountName)
	c := armstorage.NewBlobContainersClient(a.AzureArmConnection, subscriptionID)

	if _, err := c.Get(ctx, resourceGroupName, storageAccountName, containerName, nil); err != nil {
		return NoAzureStorageContainerError{
			name:           containerName,
			storageAccount: storageAccountName,
		}
	} else {
		return nil
	}
}

// CreateRoleAssignment az role assignment create --scope "/subscriptions/${SUBSCRIPTION_ID}" --role "User Access Administrator" --assignee "$CLIENT_ID".
func (a Helper) CreateRoleAssignment(ctx context.Context, subscriptionID, scope, roleName, principalID string) error {
	log.Info().Msgf("creating role assignment for principal [%s] with role [%s] and scope [%s]", principalID, roleName, scope)

	c := armauthorization.NewRoleAssignmentsClient(a.AzureArmConnection, subscriptionID)
	c2 := armauthorization.NewRoleDefinitionsClient(a.AzureArmConnection)

	res := c2.List(scope, &armauthorization.RoleDefinitionsListOptions{Filter: to.StringPtr(fmt.Sprintf("roleName eq '%s'", roleName))})
	if res.Err() != nil {
		return res.Err()
	}
	res.NextPage(ctx)
	roleDefID := res.PageResponse().Value[0].ID

	log.Info().Msgf("using role definition [%s]", to.String(roleDefID))

	if err := a.checkRoleAssignment(ctx, subscriptionID, scope, *roleDefID, principalID); err == nil {
		log.Info().Msgf("role definition [%s] already assigned to principal [%s] with scope [%s]", roleName, principalID, scope)
		return nil
	}

	log.Info().Msgf("assigning role definition [%s] to principal [%s] with scope [%s]", roleName, principalID, scope)
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

func (a Helper) checkRoleAssignment(ctx context.Context, subscriptionID, scope, roleDefID, principalID string) error {
	c := armauthorization.NewRoleAssignmentsClient(a.AzureArmConnection, subscriptionID)

	res := c.ListForScope(scope, &armauthorization.RoleAssignmentsListForScopeOptions{
		Filter: to.StringPtr(fmt.Sprintf("roleDefinitionID eq '%s' && principalId eq '%s'", roleDefID, principalID)),
	})
	if res.Err() != nil {
		return res.Err()
	}
	res.NextPage(ctx)
	if len(res.PageResponse().Value) != 0 {
		return nil
	} else {
		return NoAzureRoleAssignmentError{
			roleDefinitionID: roleDefID,
			principalID:      principalID,
			scope:            scope,
		}
	}
}

func (a Helper) CreateServicePrincipal(ctx context.Context, tenantID, name string) (string, string, error) {
	appID := ""
	objectID, secret, err := a.checkServicePrincipal(ctx, tenantID, name)
	var errNoApp *NoAzureApplicationError
	var errNoSP *NoAzureServicePrincipalError

	if err == nil {
		log.Info().Msgf("service principal [%s] already existing", name)
		return *objectID, *secret, nil
	} else if ok := errors.As(err, errNoApp); ok {
		log.Error().Msg(errNoApp.Error())
	} else if ok := errors.As(err, errNoSP); ok {
		log.Error().Msg(errNoSP.Error())
		appID = errNoSP.appID
	} else {
		return "", "", err
	}

	c := graphrbac.NewServicePrincipalsClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer
	c2 := graphrbac.NewApplicationsClient(tenantID)
	c2.Authorizer = a.AzureGraphAuthorizer

	newSecret := uuid.NewV1().String()
	if appID == "" {
		log.Info().Msgf("creating ad application with name [%s]", name)
		app, err := c2.Create(ctx, graphrbac.ApplicationCreateParameters{
			DisplayName:             to.StringPtr(name),
			AvailableToOtherTenants: to.BoolPtr(false),
		})
		if err != nil {
			return "", "", err
		}
		appID = *app.AppID
	}

	log.Info().Msgf("creating ad service principal for application [%s]", name)
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

func (a Helper) checkServicePrincipal(ctx context.Context, tenantID, name string) (*string, *string, error) {
	c := graphrbac.NewServicePrincipalsClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer
	c2 := graphrbac.NewApplicationsClient(tenantID)
	c2.Authorizer = a.AzureGraphAuthorizer

	filterQuery := fmt.Sprintf("displayName eq '%s'", name)

	res, err := c2.List(ctx, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	err = res.NextWithContext(ctx)
	if err != nil {
		return nil, nil, NoAzureApplicationError{name: name}
	}
	apps := res.Values()
	if len(apps) == 0 {
		return nil, nil, NoAzureApplicationError{name: name}
	}

	res2, err := c.List(ctx, filterQuery)
	if err != nil {
		return nil, nil, err
	}
	err = res2.NextWithContext(ctx)
	if err != nil {
		return nil, nil, NoAzureServicePrincipalError{
			name:  name,
			appID: *(apps[0].AppID),
		}
	}
	sps := res2.Values()
	if len(sps) == 0 {
		return nil, nil, NoAzureServicePrincipalError{
			name:  name,
			appID: *(apps[0].AppID),
		}
	}

	sp := sps[0]
	objID := sp.ObjectID
	secret := (*sp.PasswordCredentials)[0].Value

	return objID, secret, nil
}

// CreateADPermission az ad app permission add --id "${CLIENT_ID}" --api 00000002-0000-0000-c000-000000000000 --api-permissions 1cda74f2-2616-4834-b122-5cb1b07f8a59=Role.
func (a Helper) CreateADPermission(ctx context.Context, tenantID, clientID, scope string) error {
	log.Info().Msgf("creating ad permission for client [%s] with scope [%s]", clientID, scope)
	c := graphrbac.NewOAuth2PermissionGrantClient(tenantID)
	c.Authorizer = a.AzureGraphAuthorizer

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

func setupAuthorizationWithResource(useCLI bool, resource string) (autorest.Authorizer, error) {
	if useCLI {
		return auth.NewAuthorizerFromCLIWithResource(resource)
	} else {
		return auth.NewAuthorizerFromEnvironmentWithResource(resource)
	}
}

func setupArmConnection(useCLI bool) (*arm.Connection, error) {
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
