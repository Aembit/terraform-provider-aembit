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

resource "aembit_integration" "wiz" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://api.us17.app.wiz.io/graphql"
	oauth_client_credentials = {
		token_url = "https://auth.app.wiz.io/oauth/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

data "aembit_integrations" "test" {
	resource_set_id = aembit_resource_set.crs.id
    depends_on = [ aembit_integration.wiz ]
}
