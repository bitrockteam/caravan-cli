package aws_test

import (
	"caravan-cli/cli"
	"caravan-cli/provider/aws"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	ctx := context.Background()
	dir, err := os.MkdirTemp("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	config, _ := cli.NewConfigFromScratch("test-name", "aws", "eu-south-1")
	config.SetWorkdir(dir, "aws")
	_ = config.SetDomain("test.me")
	_ = config.SetDistro("ubuntu-2204")
	_ = config.SetEdition("ent")
	aws, _ := aws.New(ctx, config)

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
			templates, _ := aws.GetTemplates(ctx)
			for _, tmp := range templates {
				if tmp.Name == tc.name {
					// log.Info().Msgf("%s", tc.name)
					// log.Info().Msgf("test: %s", tmp.Path)

					if err := tmp.Render(aws.Caravan); err != nil {
						t.Errorf("error generating template %s: %s\n", tmp.Name, err)
					}

					want, err := os.ReadFile(gold)
					if err != nil {
						t.Fatalf("error reading golden file: %s\n", err)
					}
					got, err := os.ReadFile(tmp.Path)
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
