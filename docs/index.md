---
layout: ""
page_title: "Provider: Aembit Cloud"
description: |-
  This Aembit Cloud provider provides resources and data sources to manage the Aembit platform as infrastructure-as-code, through the Aembit management API.
---

# Aembit Provider

The Aembit provider interacts with the configuration of the Aembit platform via the management API. The provider requires credentials before it can be used.

## Getting Started

To get started using the Aembit Terraform provider, first you'll need an active Aembit cloud tenant.  Get instant access with a [Aembit trial account](https://useast2.aembit.io/signup), or read more about Aembit at [https://aembit.io](https://aembit.io)

## Authenticate using Aembit native authentication

Aembit supports authentication to the Aembit API using a native authentication capability which utilizes OIDC (Open ID Connect tokens) ID Tokens. This capability requires configuring your Aembit tenant with the appropriate components as follows:
* **Client Workload:** This Workload identifies the execution environment of the Terraform Provider, either in Terraform Cloud, GitHub Actions, or another Aembit-supported Serverless platform.
* **Trust Provider:** This Trust Provider ensures the authentication of the Client Workload using attestation of the platform ID Token and associated match rules.
  * Match Rules can be configured for platform-specific restrictions, for example a repository on GitHub or workspace ID on Terraform Cloud.
* **Credential Provider:** Created using the *Aembit Access Token* Credential Provider type, this Credential Provider can be configured with an Aembit Role that has permissions for the types of resources you want to configure.
  * **Prerequisite**: Configuring your Credential Provider with the *Aembit Access Token* type requires that you have the **read** permission for [**Roles**](https://docs.aembit.io/administration/roles/overview).
  * **Note**: The Aembit API hostname will be provided as an Audience value here and can be copied to the Server Workload hostname field.
* **Server Workload:** This Workload identifies the Aembit tenant-specific API endpoint.
  * The Host value can be copied from the Audience value of the Credential Provider.
  * The Port values should be set to 443 with TLS encryption for both the Port and Forward to Port options.
  * The Authentication section should be configured for *HTTP Authentication* with the *Bearer* authentication scheme. 
* **Access Policy:** This policy associates the previously configured elements and ensures that only the specific Terraform provider workload has access as defined.

After configuring these Aembit resources, the Client ID from the Trust Provider can be configured for the Aembit Terraform Provider, enabling automatic native authentication for the configured Workload.
The Client ID can be configured using the `client_id` field in the Aembit provider configuration block or with the `AEMBIT_CLIENT_ID` environment variable.

->**Resource Set Scoping**
When specifying the `resource_set_id` configuration to a custom Resource Set, the Aembit Cloud Terraform Provider authenticates using entities in the custom Resource Set. 
Please ensure that the entities you created above are in the same custom Resource Set.
</br>
The custom Resource Set can alternatively be specified by setting the `AEMBIT_RESOURCE_SET_ID` environment variable.

->**Terraform Cloud Configuration**
Setting the environment variable `TFC_WORKLOAD_IDENTITY_AUDIENCE` is required for Terraform Cloud Workspace ID Tokens. The value for this variable will be provided by your Aembit Cloud tenant Trust Provider and references your tenant-specific endpoint.

#### Sample Terraform Config

```terraform
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}

provider "aembit" {
  # This client_id configuration may be set here or in the AEMBIT_CLIENT_ID environment variable.
  # Note: This is a sample value and must be replaced with your Aembit Trust Provider generated value.
  client_id = "aembit:useast2:tenant:identity:github_idtoken:0bc4dbcd-e9c8-445b-ac90-28f47b8649cc"

  # Optional, defaults to the Default Resource Set
  # Note: This is a sample value and must be replaced with your generated Resource Set ID.
  resource_set_id = "d67afe77-6313-4b18-9b64-c0949b75bd1c"
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}
```

```shell
$ terraform plan
```

## Authenticate using an environment variable access token

When using the Aembit Terraform Provider on your desktop, you can retrieve an API access token from your tenant and use it by setting the `AEMBIT_TOKEN` environment variable.

**Note**: This token will only be valid for the lifetime of your Aembit UI session and should not be used in CI/CD pipelines.

#### Sample Terraform Config

```terraform
terraform {
  required_providers {
    aembit = {
      source = "aembit/aembit"
    }
  }
}


provider "aembit" {
}

resource "aembit_client_workload" "client" {
  # Resource configuration
}
```

```shell
$ export AEMBIT_TENANT_ID="tenant"
$ export AEMBIT_TOKEN="token-from-console"
$ terraform plan
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_id` (String) The Aembit Trust Provider Client ID to use for authentication to the Aembit Cloud Tenant instance (recommended).
- `resource_set_id` (String) The Aembit Resource Set to use for resources associated with this Terraform Provider.
- `tenant` (String) Tenant ID of the specific Aembit Cloud instance.
- `token` (String, Sensitive) Access Token to use for authentication to the Aembit Cloud Tenant instance.

