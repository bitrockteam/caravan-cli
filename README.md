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
|up| WIP | :x: | :x: | :x: | 
|status| :x: | :x: | :x: | :x: | 
|update| :x: | :x: | :x: | :x: | 
|clean| :x: | :x: | :x: | :x: | 





## Develop

To build the cli execute:
```
cd cmd/caravan
go build .
```

Adding support for a new command is possible by adding a new command implmentation under `commands` folder  and making it accessible in the `cmd/main.go`.
As of now the configuration for the cli is managed in the  `internal/caravan` package `Config` struct.

For the execution of the command to be successful the following pre-requisite needs to be verified on the environment:
1- terraform installed and available in the `$PATH` variable
2- aws cli installed and with credentials provided in `.aws/credentials`
