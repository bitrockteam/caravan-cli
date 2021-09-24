package terraform_test

import (
	"io/ioutil"
	"os"
	"testing"

	"caravan-cli/terraform"
)

func TestTerraformInit(t *testing.T) {
	dir, err := ioutil.TempDir("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tf := terraform.Terraform{}
	if err := tf.Init(dir); err != nil {
		t.Errorf("error during terraform init: %s", err)
	}
	if tf.Workdir != dir {
		t.Errorf("error setting terraform workdir: got %s want %s\n", tf.Workdir, dir)
	}
}
