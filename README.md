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
|init| :heavy_check_mark: | :x: | :x: | :x: |
|bake| :heavy_check_mark: | :x: | :x: | :x: |
|up| :heavy_check_mark: | :x: | :x: | :x: | 
|status| :heavy_check_mark: | :x: | :x: | :x: | 
|update| :x: | :x: | :x: | :x: | 
|clean| :heavy_check_mark: | :x: | :x: | :x: | 

## Usage
A typical session is as follows:
```
./caravan init --project --provider aws
```
This will generate in the ```.caravan``` local folder the needed variables/templates for the correspondig provider selected. In the same folder the git repos with the relevant terraform code will be checked-out with the default branch (release branch) unless the ```--branch``` optional parameter is specified.

Once the init is performed the caravan environment can be started by issuing:
```
./caravan up
```

At each point in time the status of the ongoing deployment can be checked with:
```
./caravan status
```

To delete anenvironment the following command is available:
```
./caravan clean
```
A ```--force true``` option is needed to remove all the objects from the cloud bucket in order to avoid losing state.

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
