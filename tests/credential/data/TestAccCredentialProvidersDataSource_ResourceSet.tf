provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_resource_set" "crs" {
	name = "TF Acceptance Custom ResourceSet"
	description = "TF Acceptance Custom ResourceSet"
	roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
}

resource "aembit_credential_provider" "oauth" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance OAuth"
	is_active = true
	oauth_client_credentials = {
		token_url = "https://aembit.io/token"
		client_id = "test_client_id"
		client_secret = "test_client_secret"
		scopes = "test_scopes"
		credential_style = "authHeader"
		custom_parameters = [
			{
				key = "key"
				value = "value"
				value_type = "literal"
			}
		]
	}
}

data "aembit_credential_providers" "test" {
    depends_on = [ aembit_credential_provider.oauth ]
}
