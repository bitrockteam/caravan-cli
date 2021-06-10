package aws_test

import (
	"caravan/internal/aws"
	"caravan/internal/caravan"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateConfig(t *testing.T) {

	uid := uuid.New().String()
	aws := aws.NewAWS(caravan.Config{
		Name:           "test-name",
		Provider:       "test-provider",
		Profile:        "test-profile",
		Region:         "test-region",
		Workdir:        ".testwd-" + uid,
		WorkdirProject: ".testwd-" + uid + "/test-name",
		TableName:      "test-table",
		BucketName:     "test-bucket",
	})

	defer os.RemoveAll(aws.CaravanConfig.Workdir)

	err := aws.GenerateConfig()
	if err != nil {
		t.Fatalf("error generating config: %s\n", err)
	}

	testCases := []struct {
		name  string
		ext   string
		fname string
	}{
		{"baking", "tfvars", filepath.Join(aws.CaravanConfig.WorkdirProject, aws.CaravanConfig.Provider+"-baking.tfvars")},
		{"infra", "tfvars", filepath.Join(aws.CaravanConfig.WorkdirProject, aws.CaravanConfig.Name+"-infra.tfvars")},
		{"backend", "tf", filepath.Join(aws.CaravanConfig.WorkdirProject, aws.CaravanConfig.Name+"-backend.tf")},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gold := filepath.Join("testdata", tc.name+".golden."+tc.ext)
			gen := tc.fname
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
