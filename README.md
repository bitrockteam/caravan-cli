# caravan-cli

This repo contains the source code for the command line tools and utilities used for caravan deployment.

## Status
The caravan-cli will provide the following capabilities (commands):

- init: creation of the needed configuration files and supporting state store and locking facilities for terraform
- bake: baking of the VM images for the given cloud provider
- up: startup of the caravan infrastructure deployment
- TBD

In the following table the support of the corresponding provider for each provider is reported:

|  | init | bake | up |
|--|--|--|--|
|aws| :heavy_check_mark: | :heavy_check_mark: | WIP | 
|gcp | :x: | :x: | :x: | 
|az | :x: | :x: | :x: | 
|oci | :x: | :x: | :x: | 

