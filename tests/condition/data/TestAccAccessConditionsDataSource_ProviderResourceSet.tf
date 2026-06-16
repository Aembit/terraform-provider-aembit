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

resource "aembit_integration" "crowdstrike" {
	provider = aembit.rs_manager
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
	provider = aembit.rs_manager
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
	provider = aembit.rs_manager
    depends_on = [ aembit_access_condition.crowdstrike ]
}
