provider "aembit" {
    alias = "rs_loader"
}

data "aembit_resource_sets" "all" {
    provider = aembit.rs_loader
}

locals {
    tf_testing_rs_id = [for rs in data.aembit_resource_sets.all.resource_sets : rs.id if rs.name == "TF Testing"][0]
}

// Create a Provider and Resource in the TF Testing Resource Set
provider "aembit" {
    alias = "rs_manager"
    resource_set_id = local.tf_testing_rs_id
}

resource "aembit_content_security" "crowdstrike" {
	provider = aembit.rs_manager
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
	provider = aembit.rs_manager
    depends_on = [ aembit_content_security.crowdstrike ]
}
