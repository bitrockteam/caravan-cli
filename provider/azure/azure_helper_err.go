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
