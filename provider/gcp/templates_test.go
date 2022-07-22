package gcp_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/gcp"
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	ctx := context.Background()
	dir, err := ioutil.TempDir("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	config, _ := cli.NewConfigFromScratch("test-name", "gcp", "europe-west6")
	config.SetWorkdir(dir, "gcp")
	_ = config.SetDomain("test.me")
	_ = config.SetEdition("ent")
	config.GCPDNSZone = "dns-zone"
	config.GCPParentProject = "parent-project"
	config.GCPUserEmail = "test.name@test.me"
	gcp, err := gcp.New(ctx, config)
	if err != nil {
		t.Fatalf("unable to create gcp: %s", err)
	}

	testCases := []struct {
		name        string
		gold        string
		deployNomad bool
	}{
		{"baking-vars", "baking.golden.tfvars", true},
		{"infra-vars", "infra.golden.tfvars", true},
		{"infra-vars", "infra.golden.nonomad.tfvars", false},
		{"infra-backend", "infra.golden.tf", true},
		{"platform-vars", "platform.golden.tfvars", true},
		{"platform-vars", "platform.golden.nonomad.tfvars", false},
		{"platform-backend", "platform.golden.tf", true},
		{"application-backend", "application.golden.tf", true},
		{"application-vars", "application.golden.tfvars", true},
		{"application-vars", "application.golden.nonomad.tfvars", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.DeployNomad = tc.deployNomad
			gold := filepath.Join("testdata", tc.gold)
			templates, _ := gcp.GetTemplates(ctx)
			for _, tmp := range templates {
				if tmp.Name == tc.name {
					// log.Info().Msgf("%s", tc.name)
					// log.Info().Msgf("test: %s", tmp.Path)

					if err := tmp.Render(gcp.Caravan); err != nil {
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
						t.Errorf("%s <-> %s: mismatch found with golden sample:\n%s\n-----\n%s\n", tmp.Path, gold, string(got), string(want))
					}
				}
			}
		})
	}
}
