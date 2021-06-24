package terraform_test

import (
	"caravan/internal/terraform"
	"io/ioutil"
	"os"
	"testing"
)

func TestTerraformInit(t *testing.T) {
	dir, err := ioutil.TempDir("", "caravan-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tf := terraform.NewTerraform(dir)
	err = tf.Init()
	if err != nil {
		t.Errorf("error during terraform init: %s", err)
	}
}
