---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "aembit_server_workload Resource - terraform-provider-aembit"
subcategory: ""
description: |-
  
---

# aembit_server_workload (Resource)

Resource to create and manage Server Workloads in a Aembit Cloud tenant.

## Example Usage
```terraform
resource "aembit_server_workload" "test" {
	name = "Name"
    description = "Description"
    is_active = true
	service_endpoint = {
		host = "test.host.com"
		port = 443
        tls = true
		app_protocol = "HTTP"
		transport_protocol = "TCP"
		requested_port = 443
        requested_tls = true
		tls_verification = "full"
		authentication_config = {
			"method" = "HTTP Authentication"
			"scheme" = "Bearer"
		}
	}
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name for the Server Workload.
- `service_endpoint` (Attributes) Service endpoint details. (see [below for nested schema](#nestedatt--service_endpoint))

### Optional

- `description` (String) Description for the Server Workload.
- `is_active` (Boolean) Active status of the Server Workload.
- `tags` (Map of String) Tags are key-value pairs.

### Read-Only

- `id` (String) Unique identifier of the Server Workload.

<a id="nestedatt--service_endpoint"></a>
### Nested Schema for `service_endpoint`

Required:

- `app_protocol` (String) Application Protocol of the Server Workload service endpoint. Possible values are: 
	* `Amazon Redshift`
	* `HTTP`
	* `MySQL`
	* `PostgreSQL`
	* `Redis`
	* `Snowflake`
- `host` (String) Hostname of the Server Workload service endpoint.
- `port` (Number) Port of the Server Workload service endpoint.
- `requested_port` (Number) Requested port of the Server Workload service endpoint.
- `tls_verification` (String) TLS verification configuration of the Server Workload service endpoint. Possible values are `full` (default) or `none`.
- `transport_protocol` (String) Transport protocol of the Server Workload service endpoint. This value must be set to the default `TCP`.

Optional:

- `authentication_config` (Attributes) Service authentication details. (see [below for nested schema](#nestedatt--service_endpoint--authentication_config))
- `http_headers` (Map of String) HTTP Headers are key-value pairs.
- `requested_tls` (Boolean) TLS requested on the Server Workload service endpoint.
- `tls` (Boolean) TLS indicated on the Server Workload service endpoint.

Read-Only:

- `external_id` (String) Unique identifier of the service endpoint.
- `id` (Number) Number identifier of the service endpoint.

<a id="nestedatt--service_endpoint--authentication_config"></a>
### Nested Schema for `service_endpoint.authentication_config`

Required:

- `method` (String) Server Workload Service authentication method. Possible values are: 
	* `API Key`
	* `HTTP Authentication`
	* `JWT Token Authentication`
	* `Password Authentication`
- `scheme` (String) Server Workload Service authentication scheme. Possible values are: 
	* For Authentation Method `API Key`:
		* `Header`
		* `Query Parameter`
	* For Authentation Method `HTTP Authentication`:
		* `Basic`
		* `Bearer`
		* `Header`
		* `AWS Signature v4`
	* For Authentation Method `JWT Token Authentication`:
		* `Snowflake JWT`
	* For Authentation Method `Password Authentication`:
		* `Password`

Optional:

- `config` (String) Server Workload Service authentication config. <br>This value is used to identify the HTTP Header or Query Parameter used for the associated authentication scheme. <br>**Note:** This value is required in cases where an HTTP Header or Query Parameter is required, for example with `HTTP Authentication` and scheme `Header`.



