# Aembit Terraform Provider

This is the repository for the Aembit Cloud Terraform Provider. Learn more about Aembit at https://aembit.io/

## Support, Bugs, Feature Requests

Any requests should be filed under the Issues section of this repository. All filed issues will be handled on a "best effort" basis.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.6
- [Go](https://golang.org/doc/install) >= 1.20

## Getting Started

The provider can be installed by running `terraform init`.

The provider block can be specified as follows:
```shell
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}


provider "aembit" {
}
```

### Testing

The Aembit Terraform Provider is regularly tested with every Aembit Cloud and Terraform Provider update through the use of Acceptance Testing.
These test can be run locally on your desktop and are additionally run automatically in the GitHub CI/CD pipeline using Aembit native authentication.

When running locally on your desktop, be sure to set the necessary environment variables:
```bash
export AEMBIT_TENANT_ID=<tenant-d>
export AEMBIT_TOKEN=<access-token-from-console>
```

### Documentation



Documentation can be verified using the [Terraform Registry Doc Preview](https://registry.terraform.io/tools/doc-preview).