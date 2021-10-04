package azure_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider"
	"caravan-cli/provider/azure"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()
	type test struct {
		name  string
		error bool
		desc  string
	}

	tests := []test{
		{name: "test-me", error: false, desc: "ok"},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			c, err := cli.NewConfigFromScratch(tc.name, provider.Azure, "westeurope")
			if err != nil {
				t.Fatalf("unable to create config: %s\n", err)
			}
			withAzureEnvVariables(func() {
				_, err = azure.New(ctx, c)
				if err == nil && tc.error || err != nil && !tc.error {
					t.Errorf("something wen wrong: want %t but got %s", tc.error, err)
				}
			})
		})
	}
}

func TestAzure_Bake(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	assert.Panics(t, func() { _ = a.Bake(ctx) })
}

func TestAzure_Deploy(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	assert.Panics(t, func() { _ = a.Deploy(ctx, cli.Infrastructure) })
}

func TestAzure_Destroy(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	assert.Panics(t, func() { _ = a.Destroy(ctx, cli.Infrastructure) })
}

func TestAzure_CleanProvider(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	assert.Panics(t, func() { _ = a.CleanProvider(ctx) })
}

func TestAzure_InitProvider(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	if err := a.InitProvider(ctx); err.Error() == "mkdir: no such file or directory" {
		t.Errorf("unexpected error: %v", err)
	}
	// ok to fail when a failure in Config,Save()
}

func TestAzure_Status(t *testing.T) {
	ctx := context.Background()
	a := mockProvider()

	assert.Panics(t, func() { _ = a.Status(ctx) })
}

func mockProvider() *azure.Azure {
	cfg := mockConfig()
	return &azure.Azure{
		GenericProvider: provider.GenericProvider{
			Caravan: cfg,
		},
		AzureHelper: NewHelperMock(cfg.AzureSubscriptionID),
	}
}

func mockConfig() *cli.Config {
	return &cli.Config{
		Name: "caravan-az-test",
		AzureConfig: cli.AzureConfig{
			AzureResourceGroup:  "resourceGroup",
			AzureTenantID:       "444-555-666",
			AzureSubscriptionID: "111-222-333",
		},
	}
}

func withAzureEnvVariables(f func()) {
	env := map[string]string{
		"AZURE_CLIENT_ID": "dummy",
		"AZURE_TENANT_ID": "dummy",
		"AZURE_USERNAME":  "dummy",
		"AZURE_PASSWORD":  "dummy",
	}
	for k, v := range env {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}
	f()
}
