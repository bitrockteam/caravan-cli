package aws

import (
	"caravan/internal/caravan"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TODO refactor with common generate code.
func (a AWS) GenerateConfig() (err error) {
	fmt.Printf("generating config files on: %s\n", a.Caravan.WorkdirProject)
	if err := os.MkdirAll(a.Caravan.WorkdirProject, 0777); err != nil {
		return err
	}

	for _, t := range a.Templates {
		fmt.Printf("generating %v:%s \n", t.Name, t.Path)
		if err := a.Generate(t); err != nil {
			return err
		}
	}

	return nil
}

func loadTemplates(a AWS) []caravan.Template {
	return []caravan.Template{
		{
			Name: "baking-vars",
			Text: `build_on_aws      = true
build_image_name  = "caravan-centos-image"
aws_region        = "{{ .Caravan.Region }}"
aws_instance_type = "t3.small"
`,
			Path: a.Caravan.WorkdirBakingVars,
		},
		{
			Name: "infra-vars",
			Text: `region                  = "{{ .Caravan.Region }}"
awsprofile              = "{{ .Caravan.Profile }}"
shared_credentials_file = "~/.aws/credentials"
prefix                  = "{{ .Caravan.Name }}"
personal_ip_list        = ["0.0.0.0/0"]
use_le_staging          = true
external_domain         = "{{ .Caravan.Domain }}"
tfstate_bucket_name     = "{{ .Caravan.BucketName }}"
tfstate_table_name      = "{{ .Caravan.TableName }}"
tfstate_region          = "{{ .Caravan.Region }}"
`,
			Path: a.Caravan.WorkdirInfraVars,
		},
		{
			Name: "infra-backend",
			Text: `terraform {
  backend "s3" {
    bucket         = "{{ .Caravan.BucketName }}"
    key            = "infraboot/terraform/state/terraform.tfstate"
    region         = "{{ .Caravan.Region }}"
    dynamodb_table = "{{ .Caravan.TableName }}"
  }
}
`,
			Path: a.Caravan.WorkdirInfraBackend,
		},
		{
			Name: "platform-backend",
			Text: `terraform {
  backend "s3" {
    bucket         = "{{ .Caravan.BucketName }}"
    key            = "platform/terraform/state/terraform.tfstate"
    region         = "{{ .Caravan.Region }}"
    dynamodb_table = "{{ .Caravan.TableName }}"
  }
}
`,
			Path: a.Caravan.WorkdirPlatformBackend,
		},
		{
			Name: "platform-vars",
			Text: `
vault_endpoint  = "https://vault.{{.Caravan.Name}}.{{.Caravan.Domain}}"
consul_endpoint = "https://consul.{{.Caravan.Name}}.{{.Caravan.Domain}}"
nomad_endpoint  = "https://nomad.{{.Caravan.Name}}.{{.Caravan.Domain}}"

vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{.Caravan.Provider}}/ca_certs.pem"

auth_providers = ["{{.Caravan.Provider}}"]

aws_region                  = "{{ .Caravan.Region }}"
aws_shared_credentials_file = "~/.{{.Caravan.Provider}}/credentials"
aws_profile                 = "default"

bootstrap_state_backend_provider   = "{{ .Caravan.Provider }}"
bootstrap_state_bucket_name_prefix = "{{ .Caravan.BucketName }}"
bootstrap_state_object_name_prefix = "infraboot/terraform/state"
s3_bootstrap_region                = "{{ .Caravan.Region }}"
`,
			Path: a.Caravan.WorkdirPlatformVars,
		},
		{
			Name: "application-vars",
			Text: `
vault_endpoint  = "https://vault.{{.Caravan.Name}}.{{.Caravan.Domain}}"
consul_endpoint = "https://consul.{{.Caravan.Name}}.{{.Caravan.Domain}}"
nomad_endpoint  = "https://nomad.{{.Caravan.Name}}.{{.Caravan.Domain}}"
domain = "{{.Caravan.Name}}.{{.Caravan.Domain}}"

artifacts_source_prefix    = ""
container_registry         = ""
services_domain            = "service.consul"
dc_names                   = ["{{.Caravan.Provider}}-dc"]
cloud                      = "{{.Caravan.Provider}}"
jenkins_volume_external_id = ""


vault_skip_tls_verify = true
consul_insecure_https = true
ca_cert_file          = "../caravan-infra-{{.Caravan.Provider}}/ca_certs.pem"
`,
			Path: a.Caravan.WorkdirApplicationVars,
		},
		{
			Name: "application-backend",
			Text: `terraform {
  backend "s3" {
    bucket         = "{{ .Caravan.BucketName }}"
    key            = "appsupport/terraform/state/terraform.tfstate"
    region         = "{{ .Caravan.Region }}"
    dynamodb_table = "{{ .Caravan.TableName }}"
  }
}
`,
			Path: a.Caravan.WorkdirApplicationBackend,
		},
	}
}

func (a AWS) Generate(t caravan.Template) (err error) {
	temp, err := template.New(t.Name).Parse(t.Text)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(t.Path), 0777); err != nil {
		return err
	}
	f, err := os.Create(t.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := temp.Execute(f, a); err != nil {
		return err
	}
	return nil
}
