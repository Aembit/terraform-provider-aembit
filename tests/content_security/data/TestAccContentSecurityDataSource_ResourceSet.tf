provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_resource_set" "crs" {
	name = "TF Acceptance Custom ResourceSet ${uuid()}"
	description = "TF Acceptance Custom ResourceSet"
	roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
}

resource "aembit_content_security" "crowdstrike" {
	resource_set_id = aembit_resource_set.crs.id
	name = "TF Acceptance CrowdStrike Content Security"
	is_active = true
	type = "CrowdStrikeAIDR"
	crowdstrike_falcon_aidr = {
		base_url = "https://api.crowdstrike.com/auth"
		encrypted_token = "test_token"
		fail_open = true
		max_retries = 10
		timeout_ms = 6000
	}
}

data "aembit_content_securities" "crowdstrike_datasource" {
	resource_set_id = aembit_resource_set.crs.id
    depends_on = [ aembit_content_security.crowdstrike ]
}
