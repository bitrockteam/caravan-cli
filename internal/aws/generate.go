package aws

import (
	"fmt"
	"os"
	"text/template"
)

func (a *AWS) GenerateConfig() (err error) {

	fmt.Printf("generating config files on: %s\n", a.CaravanConfig.WorkdirProject)
	err = os.MkdirAll(a.CaravanConfig.WorkdirProject, 0777)
	if err != nil {
		return err
	}

	err = a.GenerateBaking(a.CaravanConfig.WorkdirBakingVars)
	if err != nil {
		return err
	}

	err = a.GenerateInfra(a.CaravanConfig.WorkdirInfraVars)
	if err != nil {
		return err
	}

	err = a.GenerateBackend(a.CaravanConfig.WorkdirInfraBackend)
	if err != nil {
		return err
	}
	return nil
}

func (a *AWS) GenerateBaking(path string) (err error) {

	t, err := template.New("baking").Parse(`build_on_aws      = true
build_image_name  = "caravan-centos-image"
aws_region        = "{{ .CaravanConfig.Region }}"
aws_instance_type = "t3.small"
`)

	if err != nil {
		return err
	}

	f, err := os.Create(path)
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

func (a *AWS) GenerateInfra(path string) (err error) {

	t, err := template.New("infra").Parse(`region                  = "{{ .CaravanConfig.Region }}"
awsprofile              = "{{ .CaravanConfig.Profile }}"
shared_credentials_file = "~/.aws/credentials"
prefix                  = "{{ .CaravanConfig.Name }}"
personal_ip_list        = ["0.0.0.0/0"]
use_le_staging          = true
external_domain         = "{{ .CaravanConfig.Domain }}"
tfstate_bucket_name     = "{{ .CaravanConfig.BucketName }}"
tfstate_table_name      = "{{ .CaravanConfig.TableName }}"
tfstate_region          = "{{ .CaravanConfig.Region }}"
`)

	if err != nil {
		return err
	}

	f, err := os.Create(path)
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

func (a *AWS) GenerateBackend(path string) (err error) {

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

	f, err := os.Create(path)
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
