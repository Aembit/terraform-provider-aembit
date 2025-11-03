provider "aembit" {
}

resource "aembit_integration" "crowdstrike" {
	name = "TF Acceptance Crowdstrike"
	is_active = false
	type = "CrowdStrike"
	sync_frequency = 3600
	endpoint = "https://api.crowdstrike.com"
	oauth_client_credentials = {
		token_url = "https://api.crowdstrike.com/oauth2/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "crowdstrike" {
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
    depends_on = [ aembit_access_condition.crowdstrike ]
}
