provider "aembit" {
}

resource "aembit_integration" "crowdstrike" {
	name = "TF Acceptance Crowdstrike"
	is_active = true
	type = "CrowdStrike"
	sync_frequency = 3600
	endpoint = "https://endpoint"
	oauth_client_credentials = {
		token_url = "https://url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

resource "aembit_access_condition" "crowdstrike" {
	name = "TF Acceptance Crowdstrike"
	is_active = true
	integration_id = aembit_integration.crowdstrike.id
	crowdstrike_conditions = {
		max_last_seen = 3600
		match_hostname = true
		match_serial_number = true
		prevent_rfm = true
	}
}