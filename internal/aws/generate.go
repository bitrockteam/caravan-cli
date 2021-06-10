package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func (a *AWS) GenerateConfig() (err error) {

	fmt.Printf("generating config files on: %s\n", a.CaravanConfig.WorkdirProject)
	if _, err := os.Stat(a.CaravanConfig.WorkdirProject); os.IsNotExist(err) {
		err := os.MkdirAll(a.CaravanConfig.WorkdirProject, 0777)
		if err != nil {
			return err
		}
	}

	err = a.generateBaking()
	if err != nil {
		return err
	}

	err = a.generateInfra()
	if err != nil {
		return err
	}

	err = a.generateBackend()
	if err != nil {
		return err
	}
	return nil
}

func (a *AWS) generateBaking() (err error) {

	t, err := template.New("baking").Parse(`build_on_aws      = true
build_image_name  = "caravan-centos-image"
aws_region        = "{{ .CaravanConfig.Region }}"
aws_instance_type = "t3.small"
`)

	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(a.CaravanConfig.Workdir, "caravan-baking", "terraform", a.CaravanConfig.Provider+"-baking.tfvars"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, a)
	if err != nil {
		return err
	}
	return nil
}

func (a *AWS) generateInfra() (err error) {

	t, err := template.New("infra").Parse(`region                  = "{{ .CaravanConfig.Region }}"
awsprofile              = "{{ .CaravanConfig.Profile }}"
shared_credentials_file = "~/.aws/credentials"
prefix                  = "{{ .CaravanConfig.Name }}"
personal_ip_list        = ["0.0.0.0/0"]
use_le_staging          = true
external_domain         = "my-real-domain.io"
tfstate_bucket_name     = "{{ .CaravanConfig.BucketName }}"
tfstate_table_name      = "{{ .CaravanConfig.TableName }}"
tfstate_region          = "{{ .CaravanConfig.Region }}"
`)

	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(a.CaravanConfig.Workdir, "caravan-infra-aws", a.CaravanConfig.Name+"-infra.tfvars"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, a)
	if err != nil {
		return err
	}
	return nil
}

func (a *AWS) generateBackend() (err error) {

	t, err := template.New("bakend").Parse(`terraform {
  backend "s3" {
    bucket         = "{{ .CaravanConfig.BucketName }}"
    key            = "infraboot/terraform/state/terraform.tfstate"
    region         = "{{ .CaravanConfig.Region }}"
    dynamodb_table = "{{ .CaravanConfig.TableName }}"
  }
}
`)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(a.CaravanConfig.Workdir, "caravan-infra-aws", a.CaravanConfig.Name+"-backend.tf"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, a)
	if err != nil {
		return err
	}
	return nil
}
