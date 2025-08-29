provider "aembit" {
}

resource "aembit_discovery_integration" "wiz" {
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	endpoint = "https://endpoint"
	wiz_integration = {
		token_url = "https://url/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
}

data "aembit_discovery_integrations" "test" {
    depends_on = [ aembit_discovery_integration.wiz ]
}
