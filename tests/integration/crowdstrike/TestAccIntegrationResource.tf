provider "aembit" {
}

resource "aembit_integration" "crowdstrike" {
	name = "TF Acceptance Crowdstrike"
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