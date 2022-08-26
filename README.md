# caravan-cli

This repo contains the source code for the command line tools and utilities used for caravan deployment.

## Status
The caravan-cli will provide the following capabilities (commands):

- init: creation of the needed configuration files and supporting state store and locking facilities for terraform
- bake: baking of the VM images for the given cloud provider
- up: startup of the caravan infrastructure deployment
- status: provides a status of the deployment and running components
- update: update a running instance with new versions
- clean: destroy the project
- TBD

In the following table the support of the corresponding provider for each provider is reported:

|  | aws | gcp | az | oci |
|--|--|--|--|--|
|init| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: |
|bake| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: |
|up| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: | 
|status| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: | 
|update| :x: | :x: | :x: | :x: | 
|clean| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: | 
|multiOS| :heavy_check_mark: | :heavy_check_mark: | :x: | :x: |

## Usage

A typical sessioni, after the baking process was successfully completed,  is as follows:

### Init AWS

#### Prerequisites

* Authentication: currently the CLI leverages the provider's cloud SDK libraries for authentication. In an AWS context the user must be able to login with the aws cli.
* AMI images: The caravan's tools image must be available in the region where the init command is applied otherwise the ```up``` command later will fail.

#### Command line examples

```
./caravan-cli init --provider aws --project <project_name> --domain <domain_name>
```
Another example overriding some default values:

```
./caravan-cli init --provider aws --project <project_name> --linux-distro ubuntu-2204 --branch main --domain <domain_name> --region <aws_region>

```

### Init GCP

#### Prerequisites
In GCP context the following conditions must be met for the ```init``` and ```up``` command to be successful:

* Authentication: the gcloud cli must be installed and the application default authentication must be provided with the following command: ``` gcloud auth application-default login ```
* Parent project: a parent project (gcp-parent-project)  where the VM images are stored must be available. A cloud-dns zone should also be available (gcp-dns-zone)
* Billing account and organization: as part of the init step a new project is created to isolate the caravan's resources. For the creation an existing organization ID (gcp-org-id) and Billing account (gcp-billing-account-id) must be provided.
* User access rights: the authenticated user should be allowed to create in the parent project the terraform service account needed to create/access the VM images

#### Command line examples
```
./caravan-cli init --provider gcp --project <project_name> --linux-distro centos-7 --branch main --domain <doman> --region <gcp_region> --gcp-dns-zone <gcp_dns_zone> --gcp-parent-project <gcp_parent_project> --gcp-org-id <gcp_org_id> --gcp-billing-account-id <gcp_billing_account_id>
```

This will generate in the ```.caravan``` local folder the needed variables/templates for the correspondig provider selected. In the same folder the git repos with the relevant terraform code will be checked-out with the default branch (release branch) unless the ```--branch``` optional parameter is specified.

### Up

Once the init is performed the caravan environment can be started by issuing:
```
./caravan up
```

### Status

At each point in time the status of the ongoing deployment can be checked with:
```
./caravan status
```

### Delete

To delete anenvironment the following command is available:
```
./caravan clean
```
A ```--force true``` option is provided in case of an hard clean is needed. This option will execute the ```terraform destroy``` and remove all the state, regardless. After a forced delete is applied it's suggested to manually check that no resources are left over. 

## Develop

To build the cli execute:
```
go build .
```

The CLI is based on the Cobra cli library (https://github.com/spf13/cobra).
Adding support for a new command is possible by adding a new implmentation under `cmd` folder. The ```cobra``` command once installed can also be leveraged. Adding a new command is possible by issuing:
```
cobra add <cmd>
```

The configuration and state for the CLI is managed in the  `internal/caravan` package `Config` struct.

For the execution of the command to be successful the following pre-requisite needs to be verified on the environment:

- terraform installed and available in the `$PATH` environment variable
- aws cli installed and with credentials provided in `.aws/credentials`

### Pre Commit Checks

We use https://pre-commit.com/ for executing some check for each commit locally before push

Follow https://pre-commit.com/#install for the installation process.

We require the following binary installed in the machine:

- https://golangci-lint.run/usage/install/#local-installation
