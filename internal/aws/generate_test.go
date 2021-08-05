package aws_test

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	config, _ := caravan.NewConfigFromScratch("test-name", "aws", "eu-south-1")
	config.SetWorkdir(dir, "aws")
	_ = config.SetDomain("test.me")
	aws, _ := aws.New(*config)

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
			for _, tmp := range aws.Templates {
				if tmp.Name == tc.name {
					// fmt.Printf("%s\n", tc.name)
					// fmt.Printf("test: %s\n", tmp.Path)

					if err := aws.Generate(tmp); err != nil {
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

					if string(got) != string(want) {
						t.Errorf("%s <-> %s: mismatch found with golden sample:\n%s\n%s\n", tmp.Path, gold, string(got), string(want))
					}
				}
			}
		})
	}
}
