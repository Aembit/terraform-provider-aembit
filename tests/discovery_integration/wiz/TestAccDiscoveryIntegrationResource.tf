provider "aembit" {
}

resource "aembit_discovery_integration" "wiz" {
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	endpoint = "https://endpoint.url"
	wiz_integration = {
		token_url = "https://teoken.url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}
