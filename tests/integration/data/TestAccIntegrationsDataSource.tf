provider "aembit" {
}

resource "aembit_integration" "wiz" {
	name = "TF Acceptance Wiz"
	is_active = false
	type = "WizIntegrationApi"
	sync_frequency = 3600
	endpoint = "https://api.us17.app.wiz.io/graphql"
	oauth_client_credentials = {
		token_url = "https://auth.app.wiz.io/oauth/token"
		client_id = "client_id"
		client_secret = "client_secret"
		audience = "audience"
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}

data "aembit_integrations" "test" {
    depends_on = [ aembit_integration.wiz ]
}
