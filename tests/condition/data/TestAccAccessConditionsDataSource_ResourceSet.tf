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

resource "aembit_integration" "crowdstrike" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance Crowdstrike"
	is_active = false
	type = "CrowdStrike"
	sync_frequency = 3600
	endpoint = "https://api.crowdstrike.com"
	oauth_client_credentials = {
		token_url = "https://api.crowdstrike.com/oauth2/token"
		client_id = "client_id"
		client_secret = "client_secret"
	}
}

resource "aembit_access_condition" "crowdstrike" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance Crowdstrike"
	is_active = false
	integration_id = aembit_integration.crowdstrike.id
	crowdstrike_conditions = {
		max_last_seen = 3600
		match_hostname = true
		match_serial_number = true
		prevent_rfm = true
		match_mac_address = true
		match_local_ip = true
	}
}

data "aembit_access_conditions" "test" {
	resource_set_id = aembit_resource_set.crs.id
    depends_on = [ aembit_access_condition.crowdstrike ]
}
