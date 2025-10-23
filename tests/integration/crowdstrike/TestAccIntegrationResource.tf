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
	}
}