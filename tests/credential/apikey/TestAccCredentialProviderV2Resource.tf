provider "aembit" {
}

resource "aembit_credential_provider_v2" "api_key" {
	name = "TF Acceptance API Key"
	is_active = true
	api_key = {
		api_key = "test_api_key"
	}
}
