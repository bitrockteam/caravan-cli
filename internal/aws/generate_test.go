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
	aws := aws.NewAWS(caravan.Config{
		Name:           "test-name",
		Provider:       "test-provider",
		Profile:        "test-profile",
		Region:         "test-region",
		Workdir:        dir,
		WorkdirProject: dir + "/test-name",
		TableName:      "test-table",
		BucketName:     "test-bucket",
	})

	testCases := []struct {
		name   string
		ext    string
		folder string
		fname  string
	}{
		{"baking", "tfvars", aws.CaravanConfig.Workdir, aws.CaravanConfig.Provider + "-baking.tfvars"},
		{"infra", "tfvars", aws.CaravanConfig.Workdir, aws.CaravanConfig.Name + "-infra.tfvars"},
		{"backend", "tf", aws.CaravanConfig.Workdir, aws.CaravanConfig.Name + "-backend.tf"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			gen := filepath.Join(tc.folder, tc.fname)
			gold := filepath.Join("testdata", tc.name+".golden."+tc.ext)

			switch tc.name {
			case "baking":
				err = aws.GenerateBaking(gen)
				if err != nil {
					t.Fatalf("error generating %s config: %s\n", tc.name, err)
				}
			case "infra":
				err = aws.GenerateInfra(gen)
				if err != nil {
					t.Fatalf("error generating %s config: %s\n", tc.name, err)
				}
			case "backend":
				err = aws.GenerateBackend(gen)
				if err != nil {
					t.Fatalf("error generating %s config: %s\n", tc.name, err)
				}
			}

			want, err := ioutil.ReadFile(gold)
			if err != nil {
				t.Fatalf("error reading golden file: %s\n", err)
			}
			got, err := ioutil.ReadFile(gen)
			if err != nil {
				t.Fatalf("error reading current file: %s\n", err)
			}
			if string(got) != string(want) {
				t.Errorf("%s <-> %s: mismatch found with golden sample:\n%s\n%s\n", gen, gold, string(got), string(want))
			}
		})

	}

}
