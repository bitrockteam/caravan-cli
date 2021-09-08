package azure_test

import (
	"caravan/internal/azure"
	"caravan/internal/caravan"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	config, err := caravan.NewConfigFromScratch("test-name", caravan.Azure, "europewest")
	if err != nil {
		t.FailNow()
	}

	config.SetWorkdir(dir, caravan.Azure)
	_ = config.SetDomain("test.me")
	config.SetAzureSubscriptionID("111-222-333")
	config.SetAzureResourceGroup("caravan-test-rg")
	config.SetAzureBakingResourceGroup("caravan-admin")
	config.SetAzureBakingClientID("exampleClientId")
	config.SetAzureBakingClientSecret("exampleClientSecret")
	config.SetAzureStorageAccount("sg-test-01")
	config.SetAzureStorageContainerName("tfstate")
	config.SetAzureTenantID("my-tenant-111")
	config.SetAzureClientID("client1")
	config.SetAzureClientSecret("pass1")
	az, _ := azure.New(config)

	testCases := []struct {
		name string
		gold string
	}{
		{"baking-vars", "baking.golden.tfvars"},
		{"infra-vars", "infra.golden.tfvars"},
		{"infra-backend", "infra.golden.tf"},
		{"platform-vars", "platform.golden.tfvars"},
		{"platform-backend", "platform.golden.tf"},
		{"application-backend", "application.golden.tf"},
		{"application-vars", "application.golden.tfvars"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gold := filepath.Join("testdata", tc.gold)
			templates, _ := az.GetTemplates()
			for _, tmp := range templates {
				if tmp.Name == tc.name {
					// fmt.Printf("%s\n", tc.name)
					// fmt.Printf("test: %s\n", tmp.Path)

					if err := tmp.Render(az.Caravan); err != nil {
						t.Errorf("error generating template %s: %s\n", tmp.Name, err)
					}

					want, err := ioutil.ReadFile(gold)
					if err != nil {
						t.Fatalf("error reading golden file: %s\n", err)
					}
					got, err := ioutil.ReadFile(tmp.Path)
					if err != nil {
						t.Fatalf("error reading current file: %s\n", err)
					}

					if strings.Trim(string(got), "\n") != strings.Trim(string(want), "\n") {
						t.Errorf("%s <-> %s: mismatch found with golden sample:\n%s\n%s\n", tmp.Path, gold, string(got), string(want))
					}
				}
			}
		})
	}
}
