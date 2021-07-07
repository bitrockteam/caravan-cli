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
