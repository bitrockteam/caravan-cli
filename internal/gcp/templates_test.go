package gcp_test

import (
	"caravan/internal/caravan"
	"caravan/internal/gcp"
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
	config, _ := caravan.NewConfigFromScratch("test-name", "gcp", "europe-west6")
	config.SetWorkdir(dir, "gcp")
	_ = config.SetDomain("test.me")
	config.ParentProject = "parent-project"
	gcp, _ := gcp.New(config)

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
			templates, _ := gcp.GetTemplates()
			for _, tmp := range templates {
				if tmp.Name == tc.name {
					// fmt.Printf("%s\n", tc.name)
					// fmt.Printf("test: %s\n", tmp.Path)

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
						t.Errorf("%s <-> %s: mismatch found with golden sample:\n%s\n%s\n", tmp.Path, gold, string(got), string(want))
					}
				}
			}
		})
	}
}